package upolicy

import "golang.org/x/xerrors"

//DagPoolPolicy is the policy of dagpool
type DagPoolPolicy string

var (
	//ReadOnly defines read only permission
	ReadOnly DagPoolPolicy = "read-only"
	//WriteOnly defines write only permission
	WriteOnly DagPoolPolicy = "write-only"
	//ReadWrite defines read and write permission
	ReadWrite DagPoolPolicy = "read-write"
)

//AccessDenied is the error when the user is not allowed to access the dag
var AccessDenied = xerrors.Errorf("access denied")

func (dpp DagPoolPolicy) Allow(policy DagPoolPolicy) bool {
	switch dpp {
	case ReadWrite, policy:
		return true
	}
	return false
}

//CheckValid check the policy is valid
func CheckValid(policy string) bool {
	switch DagPoolPolicy(policy) {
	case ReadOnly, WriteOnly, ReadWrite:
		return true
	}
	return false
}
