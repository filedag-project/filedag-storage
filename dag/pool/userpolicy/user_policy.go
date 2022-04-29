package userpolicy

type DagPoolPolicy string

var (
	OnlyRead  DagPoolPolicy = "only-read"
	OnlyWrite DagPoolPolicy = "only-write"
	ReadWrite DagPoolPolicy = "read-write"
)

func (d *DagPoolPolicy) Allow(policy DagPoolPolicy) bool {
	if *d == ReadWrite {
		return true
	} else if *d == policy {
		return false
	} else {
		return false
	}
}
