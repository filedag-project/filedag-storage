package userpolicy

import "golang.org/x/xerrors"

type DagPoolPolicy string

var (
	OnlyRead  DagPoolPolicy = "only-read"
	OnlyWrite DagPoolPolicy = "only-write"
	ReadWrite DagPoolPolicy = "read-write"
)
var AccessDenied = xerrors.Errorf("access denied")

func (d *DagPoolPolicy) Allow(policy DagPoolPolicy) bool {
	if *d == ReadWrite {
		return true
	} else if *d == policy {
		return true
	} else {
		return false
	}
}
func CheckValid(pol string) bool {
	if DagPoolPolicy(pol) != OnlyRead && DagPoolPolicy(pol) != OnlyWrite && DagPoolPolicy(pol) != ReadWrite {
		return false
	}
	return true
}
