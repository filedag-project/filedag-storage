package blockpinner

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/dsindex"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	logging "github.com/ipfs/go-log/v2"
	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/obj/atlas"
	"path"
	"sync"
)

var log = logging.Logger("Pinner")

// ErrNotPinned is returned when trying to unpin items that are not pinned.
var ErrNotPinned = fmt.Errorf("not pinned or pinned indirectly")

// Pinner implements the Pinner interface
type Pinner struct {
	autoSync bool
	Lock     sync.RWMutex

	dstore ds.Datastore

	CidDIndex dsindex.Indexer
	CidRIndex dsindex.Indexer
	NameIndex dsindex.Indexer

	clean int64
	Dirty int64
}

const (
	linkRecursive = "recursive"
	linkDirect    = "direct"
	linkIndirect  = "indirect"
	linkInternal  = "internal"
	linkNotPinned = "not pinned"
	linkAny       = "any"
	linkAll       = "all"
)

// Mode allows to specify different types of pin (recursive, direct etc.).
// See the Pin Modes constants for a full list.
type Mode int

// Pin Modes
const (
	// Recursive pins pin the target cids along with any reachable children.
	Recursive Mode = iota

	// Direct pins pin just the target cid.
	Direct

	// Indirect pins are cids who have some ancestor pinned recursively.
	Indirect

	// Internal pins are cids used to keep the internal state of the pinner.
	Internal

	// NotPinned
	NotPinned

	// Any refers to any pinned cid
	Any
)

// ModeToString returns a human-readable name for the Mode.
func ModeToString(mode Mode) (string, bool) {
	m := map[Mode]string{
		Recursive: linkRecursive,
		Direct:    linkDirect,
		Indirect:  linkIndirect,
		Internal:  linkInternal,
		NotPinned: linkNotPinned,
		Any:       linkAny,
	}
	s, ok := m[mode]
	return s, ok
}

// sync datastore after every 50 cid repairs
const syncRepairFrequency = 50
const (
	basePath     = "/pins"
	pinKeyPath   = "/pins/pin"
	indexKeyPath = "/pins/index"
	dirtyKeyPath = "/pins/state/dirty"
)

var (
	linkDirects, linkRecursives string

	pinCidDIndexPath string
	pinCidRIndexPath string
	pinNameIndexPath string

	dirtyKey = ds.NewKey(dirtyKeyPath)

	pinAtl atlas.Atlas
)

type pin struct {
	Id       string
	Cid      cid.Cid
	Metadata map[string]interface{}
	Mode     Mode
	Name     string
}

func init() {
	directStr, ok := ModeToString(Direct)
	if !ok {
		panic("could not find Direct pin enum")
	}
	linkDirects = directStr

	recursiveStr, ok := ModeToString(Recursive)
	if !ok {
		panic("could not find Recursive pin enum")
	}
	linkRecursives = recursiveStr

	pinCidRIndexPath = path.Join(indexKeyPath, "cidRindex")
	pinCidDIndexPath = path.Join(indexKeyPath, "cidDindex")
	pinNameIndexPath = path.Join(indexKeyPath, "NameIndex")

	pinAtl = atlas.MustBuild(
		atlas.BuildEntry(pin{}).StructMap().
			AddField("Cid", atlas.StructMapEntry{SerialName: "cid"}).
			AddField("Metadata", atlas.StructMapEntry{SerialName: "metadata", OmitEmpty: true}).
			AddField("Mode", atlas.StructMapEntry{SerialName: "mode"}).
			AddField("Name", atlas.StructMapEntry{SerialName: "name", OmitEmpty: true}).
			Complete(),
		atlas.BuildEntry(cid.Cid{}).Transform().
			TransformMarshal(atlas.MakeMarshalTransformFunc(func(live cid.Cid) ([]byte, error) { return live.MarshalBinary() })).
			TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(func(serializable []byte) (cid.Cid, error) {
				c := cid.Cid{}
				err := c.UnmarshalBinary(serializable)
				if err != nil {
					return cid.Cid{}, err
				}
				return c, nil
			})).Complete(),
	)
	pinAtl = pinAtl.WithMapMorphism(atlas.MapMorphism{KeySortMode: atlas.KeySortMode_Strings})
}
func (p *pin) dsKey() ds.Key {
	return ds.NewKey(path.Join(pinKeyPath, p.Id))
}

// Pinned represents CID which has been pinned with a pinning strategy.
// The Via field allows to identify the pinning parent of this CID, in the
// case that the item is not pinned directly (but rather pinned recursively
// by some ascendant).
type Pinned struct {
	Key  cid.Cid
	Mode Mode
	Via  cid.Cid
}

//AddPin add pin
func (p *Pinner) AddPin(ctx context.Context, c cid.Cid, mode Mode, name string) (string, error) {
	// Create new pin and Dstore in datastore
	pp := newPin(c, mode, name)

	// Serialize pin
	pinData, err := encodePin(pp)
	if err != nil {
		return "", fmt.Errorf("could not encode pin: %v", err)
	}

	p.setDirty(ctx)

	// Store the pin
	err = p.dstore.Put(ctx, pp.dsKey(), pinData)
	if err != nil {
		return "", err
	}

	// Store CID index
	switch mode {
	case Recursive:
		err = p.CidRIndex.Add(ctx, c.KeyString(), pp.Id)
	case Direct:
		err = p.CidDIndex.Add(ctx, c.KeyString(), pp.Id)
	default:
		panic("pin mode must be recursive or direct")
	}
	if err != nil {
		return "", fmt.Errorf("could not add pin cid index: %v", err)
	}

	if name != "" {
		// Store name index
		err = p.NameIndex.Add(ctx, name, pp.Id)
		if err != nil {
			if mode == Recursive {
				e := p.CidRIndex.Delete(ctx, c.KeyString(), pp.Id)
				if e != nil {
					log.Errorf("error deleting index: %s", e)
				}
			} else {
				e := p.CidDIndex.Delete(ctx, c.KeyString(), pp.Id)
				if e != nil {
					log.Errorf("error deleting index: %s", e)
				}
			}
			return "", fmt.Errorf("could not add pin name index: %v", err)
		}
	}

	return pp.Id, nil
}

// RemovePinsForCid removes all pins for a cid that has the specified mode.
// Returns true if any pins, and all corresponding CID index entries, were
// removed.  Otherwise, returns false.
func (p *Pinner) RemovePinsForCid(ctx context.Context, c cid.Cid, mode Mode) (bool, error) {
	// Search for pins by CID
	var ids []string
	var err error
	cidKey := c.KeyString()
	switch mode {
	case Recursive:
		ids, err = p.CidRIndex.Search(ctx, cidKey)
	case Direct:
		ids, err = p.CidDIndex.Search(ctx, cidKey)
	case Any:
		ids, err = p.CidRIndex.Search(ctx, cidKey)
		if err != nil {
			return false, err
		}
		dIds, err := p.CidDIndex.Search(ctx, cidKey)
		if err != nil {
			return false, err
		}
		if len(dIds) != 0 {
			ids = append(ids, dIds...)
		}
	}
	if err != nil {
		return false, err
	}

	var removed bool

	// Remove the pin with the requested mode
	for _, pid := range ids {
		var pp *pin
		pp, err = p.loadPin(ctx, pid)
		if err != nil {
			if err == ds.ErrNotFound {
				p.setDirty(ctx)
				// Fix index; remove index for pin that does not exist
				switch mode {
				case Recursive:
					_, err = p.CidRIndex.DeleteKey(ctx, cidKey)
					if err != nil {
						return false, fmt.Errorf("error deleting index: %s", err)
					}
				case Direct:
					_, err = p.CidDIndex.DeleteKey(ctx, cidKey)
					if err != nil {
						return false, fmt.Errorf("error deleting index: %s", err)
					}
				case Any:
					_, err = p.CidRIndex.DeleteKey(ctx, cidKey)
					if err != nil {
						return false, fmt.Errorf("error deleting index: %s", err)
					}
					_, err = p.CidDIndex.DeleteKey(ctx, cidKey)
					if err != nil {
						return false, fmt.Errorf("error deleting index: %s", err)
					}
				}
				if err = p.FlushPins(ctx, true); err != nil {
					return false, err
				}
				// Mark this as removed since it removed an index, which is
				// what prevents determines if an item is pinned.
				removed = true
				log.Error("found CID index with missing pin")
				continue
			}
			return false, err
		}
		if mode == Any || pp.Mode == mode {
			err = p.removePin(ctx, pp)
			if err != nil {
				return false, err
			}
			removed = true
		}
	}
	return removed, nil
}

// loadPin loads a single pin from the datastore.
func (p *Pinner) loadPin(ctx context.Context, pid string) (*pin, error) {
	pinData, err := p.dstore.Get(ctx, ds.NewKey(path.Join(pinKeyPath, pid)))
	if err != nil {
		return nil, err
	}
	return decodePin(pid, pinData)
}
func (p *Pinner) removePin(ctx context.Context, pp *pin) error {
	p.setDirty(ctx)
	var err error

	// Remove cid index from datastore
	if pp.Mode == Recursive {
		err = p.CidRIndex.Delete(ctx, pp.Cid.KeyString(), pp.Id)
	} else {
		err = p.CidDIndex.Delete(ctx, pp.Cid.KeyString(), pp.Id)
	}
	if err != nil {
		return err
	}

	if pp.Name != "" {
		// Remove name index from datastore
		err = p.NameIndex.Delete(ctx, pp.Name, pp.Id)
		if err != nil {
			return err
		}
	}

	// The pin is removed last so that an incomplete remove is detected by a
	// pin that has a missing index.
	err = p.dstore.Delete(ctx, pp.dsKey())
	if err != nil {
		return err
	}

	return nil
}

// setDirty updates the dirty counter and saves a dirty state in the datastore
// if the state was previously clean
func (p *Pinner) setDirty(ctx context.Context) {
	wasClean := p.Dirty == p.clean
	p.Dirty++

	if !wasClean {
		return // do not save; was already dirty
	}

	data := []byte{1}
	err := p.dstore.Put(ctx, dirtyKey, data)
	if err != nil {
		log.Errorf("failed to set pin dirty flag: %s", err)
		return
	}
	err = p.dstore.Sync(ctx, dirtyKey)
	if err != nil {
		log.Errorf("failed to sync pin dirty flag: %s", err)
	}
}

func newPin(c cid.Cid, mode Mode, name string) *pin {
	return &pin{
		Id:   path.Base(ds.RandomKey().String()),
		Cid:  c,
		Name: name,
		Mode: mode,
	}
}
func encodePin(p *pin) ([]byte, error) {
	b, err := cbor.MarshalAtlased(p, pinAtl)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// New creates a new Pinner and loads its keysets from the given datastore. If
// there is no data present in the datastore, then an empty Pinner is returned.
//
// By default, changes are automatically flushed to the datastore.  This can be
// disabled by calling SetAutosync(false), which will require that Flush be
// called explicitly.
func New(ctx context.Context, dstore ds.Datastore) (*Pinner, error) {
	p := &Pinner{
		autoSync:  true,
		CidDIndex: dsindex.New(dstore, ds.NewKey(pinCidDIndexPath)),
		CidRIndex: dsindex.New(dstore, ds.NewKey(pinCidRIndexPath)),
		NameIndex: dsindex.New(dstore, ds.NewKey(pinNameIndexPath)),
		dstore:    dstore,
	}

	data, err := dstore.Get(ctx, dirtyKey)
	if err != nil {
		if err == ds.ErrNotFound {
			return p, nil
		}
		return nil, fmt.Errorf("cannot load dirty flag: %v", err)
	}
	if data[0] == 1 {
		p.Dirty = 1

		err = p.rebuildIndexes(ctx)
		if err != nil {
			return nil, fmt.Errorf("cannot rebuild indexes: %v", err)
		}
	}

	return p, nil
}

// rebuildIndexes uses the stored pins to rebuild secondary indexes.  This
// resolves any discrepancy between secondary indexes and pins that could
// result from a program termination between saving the two.
func (p *Pinner) rebuildIndexes(ctx context.Context) error {
	// Load all pins from the datastore.
	q := query.Query{
		Prefix: pinKeyPath,
	}
	results, err := p.dstore.Query(ctx, q)
	if err != nil {
		return err
	}
	defer results.Close()

	var checkedCount, repairedCount int

	// Iterate all pins and check if the corresponding recursive or direct
	// index is missing.  If the index is missing then create the index.
	for r := range results.Next() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if r.Error != nil {
			return fmt.Errorf("cannot read index: %v", r.Error)
		}
		ent := r.Entry
		pp, err := decodePin(path.Base(ent.Key), ent.Value)
		if err != nil {
			return err
		}

		indexKey := pp.Cid.KeyString()

		var indexer, staleIndexer dsindex.Indexer
		var idxrName, staleIdxrName string
		if pp.Mode == Recursive {
			indexer = p.CidRIndex
			staleIndexer = p.CidDIndex
			idxrName = linkRecursive
			staleIdxrName = linkDirect
		} else if pp.Mode == Direct {
			indexer = p.CidDIndex
			staleIndexer = p.CidRIndex
			idxrName = linkDirect
			staleIdxrName = linkRecursive
		} else {
			log.Error("unrecognized pin mode:", pp.Mode)
			continue
		}

		// Remove any stale index from unused indexer
		ok, err := staleIndexer.HasValue(ctx, indexKey, pp.Id)
		if err != nil {
			return err
		}
		if ok {
			// Delete any stale index
			log.Errorf("deleting stale %s pin index for cid %v", staleIdxrName, pp.Cid.String())
			if err = staleIndexer.Delete(ctx, indexKey, pp.Id); err != nil {
				return err
			}
		}

		// Check that the indexer indexes this pin
		ok, err = indexer.HasValue(ctx, indexKey, pp.Id)
		if err != nil {
			return err
		}

		var repaired bool
		if !ok {
			// Do not rebuild if index has an old value with leading slash
			ok, err = indexer.HasValue(ctx, indexKey, "/"+pp.Id)
			if err != nil {
				return err
			}
			if !ok {
				log.Errorf("repairing %s pin index for cid: %s", idxrName, pp.Cid.String())
				// There was no index found for this pin.  This was either an
				// incomplete add or and incomplete delete of a pin.  Either
				// way, restore the index to complete the add or to undo the
				// incomplete delete.
				if err = indexer.Add(ctx, indexKey, pp.Id); err != nil {
					return err
				}
				repaired = true
			}
		}
		// Check for missing name index
		if pp.Name != "" {
			ok, err = p.NameIndex.HasValue(ctx, pp.Name, pp.Id)
			if err != nil {
				return err
			}
			if !ok {
				log.Errorf("repairing name pin index for cid: %s", pp.Cid.String())
				if err = p.NameIndex.Add(ctx, pp.Name, pp.Id); err != nil {
					return err
				}
			}
			repaired = true
		}

		if repaired {
			repairedCount++
		}
		checkedCount++
		if checkedCount%syncRepairFrequency == 0 {
			p.FlushPins(ctx, true)
		}
	}

	log.Errorf("checked %d pins for invalid indexes, repaired %d pins", checkedCount, repairedCount)
	return p.FlushPins(ctx, true)
}
func (p *Pinner) FlushPins(ctx context.Context, force bool) error {
	if !p.autoSync && !force {
		return nil
	}
	if err := p.dstore.Sync(ctx, ds.NewKey(basePath)); err != nil {
		return fmt.Errorf("cannot sync pin state: %v", err)
	}
	p.setClean(ctx)
	return nil
}

// setClean saves a clean state value in the datastore if the state was
// previously dirty
func (p *Pinner) setClean(ctx context.Context) {
	if p.Dirty == p.clean {
		return // already clean
	}

	data := []byte{0}
	err := p.dstore.Put(ctx, dirtyKey, data)
	if err != nil {
		log.Errorf("failed to set clear dirty flag: %s", err)
		return
	}
	if err = p.dstore.Sync(ctx, dirtyKey); err != nil {
		log.Errorf("failed to sync cleared pin dirty flag: %s", err)
		return
	}
	p.clean = p.Dirty // set clean
}
func decodePin(pid string, data []byte) (*pin, error) {
	p := &pin{Id: pid}
	err := cbor.UnmarshalAtlased(cbor.DecodeOptions{}, data, p, pinAtl)
	if err != nil {
		return nil, err
	}
	return p, nil
}
