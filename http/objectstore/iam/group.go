package iam

// GroupDesc is a type that holds group info along with the policy
// attached to it.
type GroupDesc struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Members []string `json:"members"`
	Policy  string   `json:"policy"`
}
