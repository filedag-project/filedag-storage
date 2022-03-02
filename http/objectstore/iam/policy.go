package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/set"
	"strings"
	"time"
)

// MappedPolicy represents a policy name mapped to a user or group
type MappedPolicy struct {
	Policies string `json:"policy"`
}

// converts a mapped policy into a slice of distinct policies
func (mp MappedPolicy) toSlice() []string {
	var policies []string
	for _, policy := range strings.Split(mp.Policies, ",") {
		policy = strings.TrimSpace(policy)
		if policy == "" {
			continue
		}
		policies = append(policies, policy)
	}
	return policies
}

func (mp MappedPolicy) policySet() set.StringSet {
	return set.CreateStringSet(mp.toSlice()...)
}

func newMappedPolicy(policy string) MappedPolicy {
	return MappedPolicy{Policies: policy}
}

// PolicyDoc represents an IAM policy with some metadata.
type PolicyDoc struct {
	Version    int `json:",omitempty"`
	Policy     policy.Policy
	CreateDate time.Time `json:",omitempty"`
	UpdateDate time.Time `json:",omitempty"`
}

func newPolicyDoc(p policy.Policy) PolicyDoc {
	now := time.Now().UTC().Round(time.Millisecond)
	return PolicyDoc{
		Version:    1,
		Policy:     p,
		CreateDate: now,
		UpdateDate: now,
	}
}

// defaultPolicyDoc - used to wrap a default policy as PolicyDoc.
func defaultPolicyDoc(p policy.Policy) PolicyDoc {
	return PolicyDoc{
		Version: 1,
		Policy:  p,
	}
}

func (d *PolicyDoc) update(p policy.Policy) {
	now := time.Now().UTC().Round(time.Millisecond)
	d.UpdateDate = now
	if d.CreateDate.IsZero() {
		d.CreateDate = now
	}
	d.Policy = p
}
