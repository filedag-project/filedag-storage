package condition

import (
	"fmt"
	"reflect"
	"strconv"
)

// nullFunc - Null condition function. It checks whether Key is not present in given
// values or not.
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_Null
type nullFunc struct {
	k     Key
	value bool
}

// evaluate() - evaluates to check whether Key is present in given values or not.
// Depending on condition boolean value, this function returns true or false.
func (f nullFunc) evaluate(values map[string][]string) bool {
	rvalues := getValuesByKey(values, f.k)
	if f.value {
		return len(rvalues) == 0
	}
	return len(rvalues) != 0
}

// key() - returns condition key which is used by this condition function.
func (f nullFunc) key() Key {
	return f.k
}

// name() - returns "Null" condition name.
func (f nullFunc) name() name {
	return name{name: null}
}

func (f nullFunc) String() string {
	return fmt.Sprintf("%v:%v:%v", null, f.k, f.value)
}

// toMap - returns map representation of this function.
func (f nullFunc) toMap() map[Key]ValueSet {
	if !f.k.IsValid() {
		return nil
	}

	return map[Key]ValueSet{
		f.k: NewValueSet(NewBoolValue(f.value)),
	}
}

func (f nullFunc) clone() condFunction {
	return &nullFunc{
		k:     f.k,
		value: f.value,
	}
}

func newNullFunc(key Key, values ValueSet, qualifier string) (condFunction, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("only one value is allowed for Null condition")
	}

	var value bool
	for v := range values {
		switch v.GetType() {
		case reflect.Bool:
			value, _ = v.GetBool()
		case reflect.String:
			var err error
			s, _ := v.GetString()
			if value, err = strconv.ParseBool(s); err != nil {
				return nil, fmt.Errorf("value must be a boolean string for Null condition")
			}
		default:
			return nil, fmt.Errorf("value must be a boolean for Null condition")
		}
	}

	return &nullFunc{key, value}, nil
}
