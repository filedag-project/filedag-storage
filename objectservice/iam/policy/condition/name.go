package condition

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	// names
	stringEquals              = "StringEquals"
	stringNotEquals           = "StringNotEquals"
	stringEqualsIgnoreCase    = "StringEqualsIgnoreCase"
	stringNotEqualsIgnoreCase = "StringNotEqualsIgnoreCase"
	stringLike                = "StringLike"
	stringNotLike             = "StringNotLike"
	binaryEquals              = "BinaryEquals"
	null                      = "Null"
)

var names = map[string]struct{}{
	stringEquals:              {},
	stringNotEquals:           {},
	stringEqualsIgnoreCase:    {},
	stringNotEqualsIgnoreCase: {},
	binaryEquals:              {},
	stringLike:                {},
	stringNotLike:             {},
	null:                      {},
}

type name struct {
	name string
}

func (n name) String() string {
	return n.name
}

// IsValid - checks if name is valid or not.
func (n name) IsValid() bool {
	_, found := names[n.name]
	return found
}

// MarshalJSON - encodes name to JSON data.
func (n name) MarshalJSON() ([]byte, error) {
	if !n.IsValid() {
		return nil, fmt.Errorf("invalid name %v", n)
	}

	return json.Marshal(n.String())
}

// UnmarshalJSON - decodes JSON data to condition name.
func (n *name) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsedName, err := parseName(s)
	if err != nil {
		return err
	}

	*n = parsedName
	return nil
}

func parseName(s string) (name, error) {
	tokens := strings.Split(s, ":")
	var n name
	switch len(tokens) {
	case 0, 1:
		n = name{name: s}
	case 2:
		n = name{name: tokens[1]}
	default:
		return n, fmt.Errorf("invalid condition name '%v'", s)
	}

	if n.IsValid() {
		return n, nil
	}

	return n, fmt.Errorf("invalid condition name '%v'", s)
}
