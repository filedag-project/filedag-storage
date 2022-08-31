package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/set"
	"github.com/vmihailenco/msgpack/v5"
)

// Principal - policy principal.
//"Principal": {
//  "AWS": [
//    "arn:aws:iam::123456789012:root",
//    "999999999999",
//    "CanonicalUser": "79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be"
//  ]
//}
//The user, service, or account that receives permissions that are defined in a policy.
//The principal is A in the statement "A has permission to do B to C."
type Principal struct {
	AWS set.StringSet
}

// IsValid - checks whether Principal is valid or not.
func (p Principal) IsValid() bool {
	return len(p.AWS) != 0
}

// Equals - returns true if principals are equal.
func (p Principal) Equals(pp Principal) bool {
	return p.AWS.Equals(pp.AWS)
}

// Intersection - returns principals available in both Principal.
func (p Principal) Intersection(principal Principal) set.StringSet {
	return p.AWS.Intersection(principal.AWS)
}

// MarshalJSON - encodes Principal to JSON data.
func (p Principal) MarshalJSON() ([]byte, error) {
	if !p.IsValid() {
		return nil, errors.New(fmt.Sprintf("invalid principal %v", p))
	}

	// subtype to avoid recursive call to MarshalJSON()
	type subPrincipal Principal
	sp := subPrincipal(p)
	return json.Marshal(sp)
}

func (p Principal) MarshalMsgpack() ([]byte, error) {
	if !p.IsValid() {
		return nil, errors.New(fmt.Sprintf("invalid principal %v", p))
	}

	// subtype to avoid recursive call to MarshalJSON()
	type subPrincipal Principal
	sp := subPrincipal(p)
	return msgpack.Marshal(sp)
}

// Match - matches given principal is wildcard matching with Principal.
func (p Principal) Match(principal string) bool {
	for _, pattern := range p.AWS.ToSlice() {
		if set.MatchSimple(pattern, principal) {
			return true
		}
	}

	return false
}

// UnmarshalJSON - decodes JSON data to Principal.
func (p *Principal) UnmarshalJSON(data []byte) error {
	// subtype to avoid recursive call to UnmarshalJSON()
	type subPrincipal Principal
	var sp subPrincipal

	if err := json.Unmarshal(data, &sp); err != nil {
		var s string
		if err = json.Unmarshal(data, &s); err != nil {
			return err
		}

		if s != "*" {
			return errors.New(fmt.Sprintf("invalid principal '%v'", s))
		}

		sp.AWS = set.CreateStringSet("*")
	}

	*p = Principal(sp)

	return nil
}

func (p *Principal) UnmarshalMsgpack(data []byte) error {
	// subtype to avoid recursive call to UnmarshalJSON()
	type subPrincipal Principal
	var sp subPrincipal

	if err := msgpack.Unmarshal(data, &sp); err != nil {
		var s string
		if err = msgpack.Unmarshal(data, &s); err != nil {
			return err
		}

		if s != "*" {
			return errors.New(fmt.Sprintf("invalid principal '%v'", s))
		}

		sp.AWS = set.CreateStringSet("*")
	}

	*p = Principal(sp)

	return nil
}

// Clone clones Principal structure
func (p Principal) Clone() Principal {
	return NewPrincipal(p.AWS.ToSlice()...)
}

// NewPrincipal - creates new Principal.
func NewPrincipal(principals ...string) Principal {
	return Principal{AWS: set.CreateStringSet(principals...)}
}
