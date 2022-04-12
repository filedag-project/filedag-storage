package iam

// GroupDesc is a type that holds group info along with the policy
// attached to it.
type GroupDesc struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Members []string `json:"members"`
	Policy  string   `json:"policy"`
}

// GroupInfo contains info about a group
type GroupInfo struct {
	Name    string   `json:"name"`
	Version int      `json:"version"`
	Status  string   `json:"status"`
	Members []string `json:"members"`
}

func newGroupInfo(members []string) GroupInfo {
	return GroupInfo{Version: 1, Status: statusEnabled, Members: members}
}
