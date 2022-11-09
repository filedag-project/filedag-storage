package condition

import (
	"encoding/json"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
)

// ValueSet - unique list of values.
type ValueSet map[Value]struct{}

// Add - adds given value to value set.
func (set ValueSet) Add(value Value) {
	set[value] = struct{}{}
}

// ToSlice converts ValueSet to a slice of Value
func (set ValueSet) ToSlice() []Value {
	var values []Value
	for k := range set {
		values = append(values, k)
	}
	return values
}

func valueArrayToValueSet(values []Value) (*ValueSet, error) {
	vset := make(ValueSet)
	for _, v := range values {
		if _, found := vset[v]; found {
			return nil, fmt.Errorf("duplicate value found '%v'", v)
		}

		vset.Add(v)
	}
	return &vset, nil
}

// MarshalJSON - encodes ValueSet to JSON data.
func (set ValueSet) MarshalJSON() ([]byte, error) {
	values := set.ToSlice()
	if len(values) == 0 {
		return nil, fmt.Errorf("invalid value set %v", set)
	}

	return json.Marshal(values)
}

// UnmarshalJSON - decodes JSON data.
func (set *ValueSet) UnmarshalJSON(data []byte) error {
	var v Value
	if err := json.Unmarshal(data, &v); err == nil {
		*set = make(ValueSet)
		set.Add(v)
		return nil
	}

	var values []Value
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	if len(values) < 1 {
		return fmt.Errorf("invalid value")
	}

	vset, err := valueArrayToValueSet(values)
	if err != nil {
		return nil
	}
	*set = *vset
	return nil
}

func (set ValueSet) MarshalMsgpack() ([]byte, error) {
	values := set.ToSlice()
	if len(values) == 0 {
		return nil, fmt.Errorf("invalid value set %v", set)
	}

	return msgpack.Marshal(values)
}

func (set *ValueSet) UnmarshalMsgpack(data []byte) error {
	var v Value
	if err := msgpack.Unmarshal(data, &v); err == nil {
		*set = make(ValueSet)
		set.Add(v)
		return nil
	}

	var values []Value
	if err := msgpack.Unmarshal(data, &values); err != nil {
		return err
	}

	if len(values) < 1 {
		return fmt.Errorf("invalid value")
	}

	vset, err := valueArrayToValueSet(values)
	if err != nil {
		return nil
	}
	*set = *vset
	return nil
}

// Clone clones ValueSet structure
func (set ValueSet) Clone() ValueSet {
	return NewValueSet(set.ToSlice()...)
}

// NewValueSet - returns new value set containing given values.
func NewValueSet(values ...Value) ValueSet {
	set := make(ValueSet)

	for _, value := range values {
		set.Add(value)
	}

	return set
}
