package upolicy

import "golang.org/x/xerrors"

//DagPoolPolicy is the policy of dagpool
type DagPoolPolicy string

var (
	//OnlyRead only read
	OnlyRead DagPoolPolicy = "only-read"
	//OnlyWrite only write
	OnlyWrite DagPoolPolicy = "only-write"
	//ReadWrite read and write
	ReadWrite DagPoolPolicy = "read-write"
)

//AccessDenied is the error when the user is not allowed to access the dag
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

//CheckValid check the policy is valid
func CheckValid(pol string) bool {
	if DagPoolPolicy(pol) != OnlyRead && DagPoolPolicy(pol) != OnlyWrite && DagPoolPolicy(pol) != ReadWrite {
		return false
	}
	return true
}
