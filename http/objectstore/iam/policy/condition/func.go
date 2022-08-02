package condition

import (
	"encoding/json"
	"fmt"
	"sort"
)

// CondFunction - condition function interface.
//https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html
//String condition operators
//Numeric condition operators
//Date condition operators
//Boolean condition operators
//Binary condition operators
//IP address condition operators
//Amazon Resource Name (ARN) condition operators
//...IfExists condition operators
//Condition operator to check existence of condition keys
type CondFunction interface {
	// evaluate() - evaluates this condition function with given values.
	evaluate(values map[string][]string) bool

	// key() - returns condition key used in this function.
	key() Key

	// name() - returns condition name of this function.
	name() name

	//String () - returns string representation of function.
	String() string

	// toMap - returns map representation of this function.
	toMap() map[Key]ValueSet

	// clone - returns copy of this function.
	clone() CondFunction
}

// Conditions - list of functions.
type Conditions []CondFunction

// Evaluate - evaluates all functions with given values map. Each function is evaluated
// sequencely and next function is called only if current function succeeds.
func (cs Conditions) Evaluate(values map[string][]string) bool {
	for _, f := range cs {
		if !f.evaluate(values) {
			return false
		}
	}

	return true
}

// Keys - returns list of keys used in all functions.
func (cs Conditions) Keys() KeySet {
	keySet := NewKeySet()

	for _, f := range cs {
		keySet.Add(f.key())
	}

	return keySet
}

// Clone clones Conditions structure
func (cs Conditions) Clone() Conditions {
	funcs := []CondFunction{}
	for _, f := range cs {
		funcs = append(funcs, f.clone())
	}
	return funcs
}

// Equals returns true if two Conditions structures are equal
func (cs Conditions) Equals(funcs Conditions) bool {
	if len(cs) != len(funcs) {
		return false
	}
	for _, fi := range cs {
		fistr := fi.String()
		found := false
		for _, fj := range funcs {
			if fistr == fj.String() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// MarshalJSON - encodes Conditions to JSON data.
func (cs Conditions) MarshalJSON() ([]byte, error) {
	nm := make(map[string]map[string]ValueSet)

	for _, f := range cs {
		fname := f.name().String()
		if _, ok := nm[fname]; !ok {
			nm[fname] = map[string]ValueSet{}
		}
		for k, v := range f.toMap() {
			nm[fname][k.String()] = v
		}
	}

	return json.Marshal(nm)
}

func (cs Conditions) String() string {
	funcStrings := []string{}
	for _, f := range cs {
		s := fmt.Sprintf("%v", f)
		funcStrings = append(funcStrings, s)
	}
	sort.Strings(funcStrings)

	return fmt.Sprintf("%v", funcStrings)
}

var conditionFuncMap = map[string]func(Key, ValueSet, string) (CondFunction, error){
	stringEquals:              newStringEqualsFunc,
	stringNotEquals:           newStringNotEqualsFunc,
	stringEqualsIgnoreCase:    newStringEqualsIgnoreCaseFunc,
	stringNotEqualsIgnoreCase: newStringNotEqualsIgnoreCaseFunc,
	binaryEquals:              newBinaryEqualsFunc,
	stringLike:                newStringLikeFunc,
	stringNotLike:             newStringNotLikeFunc,

	null: newNullFunc,
	// todo Add  conditions
}

// UnmarshalJSON - decodes JSON data to Conditions.
func (cs *Conditions) UnmarshalJSON(data []byte) error {
	nm := make(map[string]map[string]ValueSet)
	if err := json.Unmarshal(data, &nm); err != nil {
		return err
	}

	if len(nm) == 0 {
		return fmt.Errorf("condition must not be empty")
	}

	var funcs []CondFunction
	for nameString, args := range nm {
		n, err := parseName(nameString)
		if err != nil {
			return err
		}

		for keyString, values := range args {
			key, err := parseKey(keyString)
			if err != nil {
				return err
			}

			fn, ok := conditionFuncMap[n.name]
			if !ok {
				return fmt.Errorf("condition %v is not handled", n)
			}

			f, err := fn(key, values, "")
			if err != nil {
				return err
			}

			funcs = append(funcs, f)
		}
	}

	*cs = funcs

	return nil
}

// GobEncode - encodes Conditions to gob data.
func (cs Conditions) GobEncode() ([]byte, error) {
	return cs.MarshalJSON()
}

// GobDecode - decodes gob data to Conditions.
func (cs *Conditions) GobDecode(data []byte) error {
	return cs.UnmarshalJSON(data)
}

// NewConFunctions - returns new Conditions with given function list.
func NewConFunctions(conditions ...CondFunction) Conditions {
	return conditions
}
