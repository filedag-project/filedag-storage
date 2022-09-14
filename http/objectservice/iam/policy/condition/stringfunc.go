package condition

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/filedag-project/filedag-storage/http/objectservice/iam/set"
)

func substitute(values map[string][]string) func(string) string {
	return func(v string) string {
		for _, key := range CommonKeys {
			// Empty values are not supported for policy variables.
			if rvalues, ok := values[key.Name()]; ok && rvalues[0] != "" {
				v = strings.Replace(v, key.VarName(), rvalues[0], -1)
			}
		}
		return v
	}
}

//https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_String
// StringEquals Exact matching, case sensitive
// StringNotEquals Negated matching
// StringEqualsIgnoreCase Exact matching, ignoring case
// StringNotEqualsIgnoreCase Negated matching, ignoring case
// StringLike Case-sensitive matching. The values can include multi-character match wildcards (*) and single-character match wildcards (?) anywhere in the string.
// StringNotLike Negated case-sensitive matching. The values can include multi-character match wildcards (*) or single-character match wildcards (?) anywhere in the string.
type stringFunc struct {
	n          name
	k          Key
	values     set.StringSet
	ignoreCase bool
	base64     bool
	negate     bool
}

func (f stringFunc) eval(values map[string][]string) bool {
	rvalues := set.CreateStringSet(getValuesByKey(values, f.k)...)
	fvalues := f.values.ApplyFunc(substitute(values))
	if f.ignoreCase {
		rvalues = rvalues.ApplyFunc(strings.ToLower)
		fvalues = fvalues.ApplyFunc(strings.ToLower)
	}
	ivalues := rvalues.Intersection(fvalues)
	return !ivalues.IsEmpty()
}

func (f stringFunc) evaluate(values map[string][]string) bool {
	result := f.eval(values)
	if f.negate {
		return !result
	}
	return result
}

func (f stringFunc) key() Key {
	return f.k
}

func (f stringFunc) name() name {
	return f.n
}

func (f stringFunc) String() string {
	valueStrings := f.values.ToSlice()
	sort.Strings(valueStrings)
	return fmt.Sprintf("%v:%v:%v", f.n, f.k, valueStrings)
}

func (f stringFunc) toMap() map[Key]ValueSet {
	if !f.k.IsValid() {
		return nil
	}

	values := NewValueSet()
	for _, value := range f.values.ToSlice() {
		if f.base64 {
			values.Add(NewStringValue(base64.StdEncoding.EncodeToString([]byte(value))))
		} else {
			values.Add(NewStringValue(value))
		}
	}

	return map[Key]ValueSet{
		f.k: values,
	}
}

func (f stringFunc) copy() stringFunc {
	return stringFunc{
		n:          f.n,
		k:          f.k,
		values:     f.values.Union(set.NewStringSet()),
		ignoreCase: f.ignoreCase,
		base64:     f.base64,
		negate:     f.negate,
	}
}

func (f stringFunc) clone() CondFunction {
	c := f.copy()
	return &c
}

// stringLikeFunc - String like function. It checks whether value by Key in given
// values map is widcard matching in condition values.
// For example,
//   - if values = ["mybucket/foo*"], at evaluate() it returns whether string
//     in value map for Key is wildcard matching in values.
type stringLikeFunc struct {
	stringFunc
}

func (f stringLikeFunc) eval(values map[string][]string) bool {
	rvalues := getValuesByKey(values, f.k)
	fvalues := f.values.ApplyFunc(substitute(values))
	for _, v := range rvalues {
		matched := !fvalues.FuncMatch(set.Match, v).IsEmpty()
		if matched {
			return true
		}
	}
	return false
}

// evaluate() - evaluates to check whether value by Key in given values is wildcard
// matching in condition values.
func (f stringLikeFunc) evaluate(values map[string][]string) bool {
	result := f.eval(values)
	if f.negate {
		return !result
	}
	return result
}

func (f stringLikeFunc) clone() CondFunction {
	return &stringLikeFunc{stringFunc: f.copy()}
}

func valuesToStringSlice(n string, values ValueSet) ([]string, error) {
	valueStrings := []string{}

	for value := range values {
		s, err := value.GetString()
		if err != nil {
			return nil, fmt.Errorf("value must be a string for %v condition", n)
		}

		valueStrings = append(valueStrings, s)
	}

	return valueStrings, nil
}

func validateStringValues(n string, key Key, values set.StringSet) error {
	// todo: validate key and values

	return nil
}

func newStringFunc(n string, key Key, values ValueSet, qualifier string, ignoreCase, base64, negate bool) (*stringFunc, error) {
	valueStrings, err := valuesToStringSlice(n, values)
	if err != nil {
		return nil, err
	}

	sset := set.CreateStringSet(valueStrings...)
	if err := validateStringValues(n, key, sset); err != nil {
		return nil, err
	}

	return &stringFunc{
		n:          name{name: n},
		k:          key,
		values:     sset,
		ignoreCase: ignoreCase,
		base64:     base64,
		negate:     negate,
	}, nil
}

// newStringEqualsFunc - returns new StringEquals function.
func newStringEqualsFunc(key Key, values ValueSet, qualifier string) (CondFunction, error) {
	return newStringFunc(stringEquals, key, values, qualifier, false, false, false)
}

// NewStringEqualsFunc - returns new StringEquals function.
func NewStringEqualsFunc(qualifier string, key Key, values ...string) (CondFunction, error) {
	vset := NewValueSet()
	for _, value := range values {
		vset.Add(NewStringValue(value))
	}
	return newStringFunc(stringEquals, key, vset, qualifier, false, false, false)
}

// newStringNotEqualsFunc - returns new StringNotEquals function.
func newStringNotEqualsFunc(key Key, values ValueSet, qualifier string) (CondFunction, error) {
	return newStringFunc(stringNotEquals, key, values, qualifier, false, false, true)
}

// newStringEqualsIgnoreCaseFunc - returns new StringEqualsIgnoreCase function.
func newStringEqualsIgnoreCaseFunc(key Key, values ValueSet, qualifier string) (CondFunction, error) {
	return newStringFunc(stringEqualsIgnoreCase, key, values, qualifier, true, false, false)
}

// newStringNotEqualsIgnoreCaseFunc - returns new StringNotEqualsIgnoreCase function.
func newStringNotEqualsIgnoreCaseFunc(key Key, values ValueSet, qualifier string) (CondFunction, error) {
	return newStringFunc(stringNotEqualsIgnoreCase, key, values, qualifier, true, false, true)
}

// newBinaryEqualsFunc - returns new BinaryEquals function.
func newBinaryEqualsFunc(key Key, values ValueSet, qualifier string) (CondFunction, error) {
	valueStrings, err := valuesToStringSlice(binaryEquals, values)
	if err != nil {
		return nil, err
	}

	return NewBinaryEqualsFunc(qualifier, key, valueStrings...)
}

// NewBinaryEqualsFunc - returns new BinaryEquals function.
func NewBinaryEqualsFunc(qualifier string, key Key, values ...string) (CondFunction, error) {
	vset := NewValueSet()
	for _, value := range values {
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, err
		}
		vset.Add(NewStringValue(string(data)))
	}
	return newStringFunc(binaryEquals, key, vset, qualifier, false, true, false)
}

// newStringLikeFunc - returns new StringLike function.
func newStringLikeFunc(key Key, values ValueSet, qualifier string) (CondFunction, error) {
	sf, err := newStringFunc(stringLike, key, values, qualifier, false, false, false)
	if err != nil {
		return nil, err
	}

	return &stringLikeFunc{*sf}, nil
}

// newStringNotLikeFunc - returns new StringNotLike function.
func newStringNotLikeFunc(key Key, values ValueSet, qualifier string) (CondFunction, error) {
	sf, err := newStringFunc(stringNotLike, key, values, qualifier, false, false, true)
	if err != nil {
		return nil, err
	}

	return &stringLikeFunc{*sf}, nil
}
