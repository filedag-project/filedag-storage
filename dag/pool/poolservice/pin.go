package poolservice

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/blockpinner"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	"sync"
)

func (d *Pool) Pin(ctx context.Context, c cid.Cid, recurse bool) error {
	cidKey := c.KeyString()

	d.pinner.Lock.Lock()
	defer d.pinner.Lock.Unlock()

	if recurse {
		found, err := d.pinner.CidRIndex.HasAny(ctx, cidKey)
		if err != nil {
			return err
		}
		if found {
			return nil
		}

		dirtyBefore := d.pinner.Dirty

		// Only look again if something has changed.
		if d.pinner.Dirty != dirtyBefore {
			found, err = d.pinner.CidRIndex.HasAny(ctx, cidKey)
			if err != nil {
				return err
			}
			if found {
				return nil
			}
		}

		// TODO: remove this to support multiple pins per CID
		found, err = d.pinner.CidDIndex.HasAny(ctx, cidKey)
		if err != nil {
			return err
		}
		if found {
			_, err = d.pinner.RemovePinsForCid(ctx, c, blockpinner.Direct)
			if err != nil {
				return err
			}
		}

		_, err = d.pinner.AddPin(ctx, c, blockpinner.Recursive, "")
		if err != nil {
			return err
		}
	} else {
		found, err := d.pinner.CidRIndex.HasAny(ctx, cidKey)
		if err != nil {
			return err
		}
		if found {
			return fmt.Errorf("%s already pinned recursively", c.String())
		}

		_, err = d.pinner.AddPin(ctx, c, blockpinner.Direct, "")
		if err != nil {
			return err
		}
	}
	return d.pinner.FlushPins(ctx, false)
}

func (d *Pool) UnPin(ctx context.Context, c cid.Cid, recursive bool) error {
	cidKey := c.KeyString()

	d.pinner.Lock.Lock()
	defer d.pinner.Lock.Unlock()

	// TODO: use Ls() to lookup pins when new pinning API available
	/*
		matchSpec := map[string][]string {
			"cid": []string{c.String}
		}
		matches := p.Ls(matchSpec)
	*/
	has, err := d.pinner.CidRIndex.HasAny(ctx, cidKey)
	if err != nil {
		return err
	}

	if has {
		if !recursive {
			return fmt.Errorf("%s is pinned recursively", c.String())
		}
	} else {
		has, err = d.pinner.CidDIndex.HasAny(ctx, cidKey)
		if err != nil {
			return err
		}
		if !has {
			return blockpinner.ErrNotPinned
		}
	}

	removed, err := d.pinner.RemovePinsForCid(ctx, c, blockpinner.Any)
	if err != nil {
		return err
	}
	if !removed {
		return nil
	}

	return d.pinner.FlushPins(ctx, false)
}

func (d *Pool) IsPinned(ctx context.Context, cid cid.Cid) bool {
	pinned, err := d.CheckIfPinned(ctx, cid)
	if err != nil {
		return false
	}
	if len(pinned) == 0 {
		return false
	}
	return pinned[0].Pinned()
}

// CheckIfPinned checks if a set of keys are pinned, more efficient than
// calling IsPinned for each key, returns the pinned status of cid(s)
//
// TODO: If a CID is pinned by multiple pins, should they all be reported?
func (d *Pool) CheckIfPinned(ctx context.Context, cids ...cid.Cid) ([]blockpinner.Pinned, error) {
	pinned := make([]blockpinner.Pinned, 0, len(cids))
	toCheck := cid.NewSet()

	d.pinner.Lock.RLock()
	defer d.pinner.Lock.RUnlock()

	// First check for non-Indirect pins directly
	for _, c := range cids {
		cidKey := c.KeyString()
		has, err := d.pinner.CidRIndex.HasAny(ctx, cidKey)
		if err != nil {
			return nil, err
		}
		if has {
			pinned = append(pinned, blockpinner.Pinned{Key: c, Mode: blockpinner.Recursive})
		} else {
			has, err = d.pinner.CidDIndex.HasAny(ctx, cidKey)
			if err != nil {
				return nil, err
			}
			if has {
				pinned = append(pinned, blockpinner.Pinned{Key: c, Mode: blockpinner.Direct})
			} else {
				toCheck.Add(c)
			}
		}
	}

	var e error
	visited := cid.NewSet()
	err := d.pinner.CidRIndex.ForEach(ctx, "", func(key, value string) bool {
		var rk cid.Cid
		rk, e = cid.Cast([]byte(key))
		if e != nil {
			return false
		}
		e = d.Walk(ctx, rk, func(c cid.Cid) bool {
			if toCheck.Len() == 0 || !visited.Visit(c) {
				return false
			}

			if toCheck.Has(c) {
				pinned = append(pinned, blockpinner.Pinned{Key: c, Mode: blockpinner.Indirect, Via: rk})
				toCheck.Remove(c)
			}

			return true
		}, concurrent())
		if e != nil {
			return false
		}
		return toCheck.Len() > 0
	})
	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, e
	}

	// Anything left in toCheck is not pinned
	for _, k := range toCheck.Keys() {
		pinned = append(pinned, blockpinner.Pinned{Key: k, Mode: blockpinner.NotPinned})
	}

	return pinned, nil
}

// walkOptions represent the parameters of a graph walking algorithm
type walkOptions struct {
	SkipRoot     bool
	Concurrency  int
	ErrorHandler func(c cid.Cid, err error) error
}

// concurrent is a WalkOption indicating that node fetching should be done in
// parallel, with the default concurrency factor.
// NOTE: When using that option, the walk order is *not* guarantee.
// NOTE: It *does not* make multiple concurrent calls to the passed `visit` function.
func concurrent() WalkOption {
	return func(walkOptions *walkOptions) {
		walkOptions.Concurrency = defaultConcurrentFetch
	}
}

// defaultConcurrentFetch is the default maximum number of concurrent fetches
// that 'fetchNodes' will start at a time
const defaultConcurrentFetch = 32

// WalkOption is a setter for walkOptions
type WalkOption func(*walkOptions)

// Walk WalkGraph will walk the dag in order (depth first) starting at the given root.
func (d *Pool) Walk(ctx context.Context, c cid.Cid, visit func(cid.Cid) bool, options ...WalkOption) error {
	visitDepth := func(c cid.Cid, depth int) bool {
		return visit(c)
	}

	return d.WalkDepth(ctx, c, visitDepth, options...)
}

// WalkDepth walks the dag starting at the given root and passes the current
// depth to a given visit function. The visit function can be used to limit DAG
// exploration.
func (d *Pool) WalkDepth(ctx context.Context, c cid.Cid, visit func(cid.Cid, int) bool, options ...WalkOption) error {
	opts := &walkOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.Concurrency > 1 {
		return d.parallelWalkDepth(ctx, c, visit, opts)
	} else {
		return d.sequentialWalkDepth(ctx, c, 0, visit, opts)
	}
}
func (d *Pool) parallelWalkDepth(ctx context.Context, root cid.Cid, visit func(cid.Cid, int) bool, options *walkOptions) error {
	type cidDepth struct {
		cid   cid.Cid
		depth int
	}

	type linksDepth struct {
		links []*format.Link
		depth int
	}

	feed := make(chan cidDepth)
	out := make(chan linksDepth)
	done := make(chan struct{})

	var visitlk sync.Mutex
	var wg sync.WaitGroup

	errChan := make(chan error)
	fetchersCtx, cancel := context.WithCancel(ctx)
	defer wg.Wait()
	defer cancel()
	for i := 0; i < options.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cdepth := range feed {
				ci := cdepth.cid
				depth := cdepth.depth

				var shouldVisit bool

				// bypass the root if needed
				if !(options.SkipRoot && depth == 0) {
					visitlk.Lock()
					shouldVisit = visit(ci, depth)
					visitlk.Unlock()
				} else {
					shouldVisit = true
				}

				if shouldVisit {
					links, err := d.GetLinks(ctx, ci)
					if err != nil && options.ErrorHandler != nil {
						err = options.ErrorHandler(root, err)
					}
					if err != nil {
						select {
						case errChan <- err:
						case <-fetchersCtx.Done():
						}
						return
					}

					outLinks := linksDepth{
						links: links,
						depth: depth + 1,
					}

					select {
					case out <- outLinks:
					case <-fetchersCtx.Done():
						return
					}
				}
				select {
				case done <- struct{}{}:
				case <-fetchersCtx.Done():
				}
			}
		}()
	}
	defer close(feed)

	send := feed
	var todoQueue []cidDepth
	var inProgress int

	next := cidDepth{
		cid:   root,
		depth: 0,
	}

	for {
		select {
		case send <- next:
			inProgress++
			if len(todoQueue) > 0 {
				next = todoQueue[0]
				todoQueue = todoQueue[1:]
			} else {
				next = cidDepth{}
				send = nil
			}
		case <-done:
			inProgress--
			if inProgress == 0 && !next.cid.Defined() {
				return nil
			}
		case linksDepth := <-out:
			for _, lnk := range linksDepth.links {
				cd := cidDepth{
					cid:   lnk.Cid,
					depth: linksDepth.depth,
				}

				if !next.cid.Defined() {
					next = cd
					send = feed
				} else {
					todoQueue = append(todoQueue, cd)
				}
			}
		case err := <-errChan:
			return err

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
func (d *Pool) sequentialWalkDepth(ctx context.Context, root cid.Cid, depth int, visit func(cid.Cid, int) bool, options *walkOptions) error {
	if !(options.SkipRoot && depth == 0) {
		if !visit(root, depth) {
			return nil
		}
	}
	links, err := d.GetLinks(ctx, root)
	if err != nil && options.ErrorHandler != nil {
		err = options.ErrorHandler(root, err)
	}
	if err != nil {
		return err
	}

	for _, lnk := range links {
		if err := d.sequentialWalkDepth(ctx, lnk.Cid, depth+1, visit, options); err != nil {
			return err
		}
	}
	return nil
}
