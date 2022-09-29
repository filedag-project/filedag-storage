package policy

import (
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/iam/set"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/xerrors"
	"sort"
)

// ResourceSet - set of resources in policy statement.
type ResourceSet map[Resource]struct{}

// Add - adds resource to resource set.
func (resourceSet ResourceSet) Add(resource Resource) {
	resourceSet[resource] = struct{}{}
}

// Match - matches object name with anyone of resource pattern in resource set.
func (resourceSet ResourceSet) Match(resource string, conditionValues map[string][]string) bool {
	for r := range resourceSet {
		if r.Match(resource, conditionValues) {
			return true
		}
	}

	return false
}

// Equals - checks whether given resource set is equal to current resource set or not.
func (resourceSet ResourceSet) Equals(sresourceSet ResourceSet) bool {
	// If length of set is not equal to length of given set, the
	// set is not equal to given set.
	if len(resourceSet) != len(sresourceSet) {
		return false
	}

	// As both sets are equal in length, check each elements are equal.
	for k := range resourceSet {
		if _, ok := sresourceSet[k]; !ok {
			return false
		}
	}

	return true
}

// MarshalJSON - encodes ResourceSet to JSON data.
func (resourceSet ResourceSet) MarshalJSON() ([]byte, error) {
	if len(resourceSet) == 0 {
		return nil, xerrors.Errorf("empty resources not allowed")
	}

	return json.Marshal(resourceSet.ToSlice())
}

func (resourceSet ResourceSet) MarshalMsgpack() ([]byte, error) {
	if len(resourceSet) == 0 {
		return nil, xerrors.Errorf("empty resources not allowed")
	}

	return msgpack.Marshal(resourceSet.ToSlice())
}

func (resourceSet ResourceSet) String() string {
	resources := []string{}
	for resource := range resourceSet {
		resources = append(resources, resource.String())
	}
	sort.Strings(resources)

	return fmt.Sprintf("%v", resources)
}

// UnmarshalJSON - decodes JSON data to ResourceSet.
func (resourceSet *ResourceSet) UnmarshalJSON(data []byte) error {
	var sset set.StringSet
	if err := json.Unmarshal(data, &sset); err != nil {
		return err
	}

	*resourceSet = make(ResourceSet)
	for _, s := range sset.ToSlice() {
		resource, err := parseResource(s)
		if err != nil {
			return err
		}

		if _, found := (*resourceSet)[resource]; found {
			return xerrors.Errorf("duplicate resource '%v' found", s)
		}

		resourceSet.Add(resource)
	}

	return nil
}

func (resourceSet *ResourceSet) UnmarshalMsgpack(data []byte) error {
	var sset set.StringSet
	if err := msgpack.Unmarshal(data, &sset); err != nil {
		return err
	}

	*resourceSet = make(ResourceSet)
	for _, s := range sset.ToSlice() {
		resource, err := parseResource(s)
		if err != nil {
			return err
		}

		if _, found := (*resourceSet)[resource]; found {
			return xerrors.Errorf("duplicate resource '%v' found", s)
		}

		resourceSet.Add(resource)
	}

	return nil
}

// Validate - validates ResourceSet is for given bucket or not.
func (resourceSet ResourceSet) Validate(bucketName string) error {
	for resource := range resourceSet {
		if err := resource.Validate(bucketName); err != nil {
			return err
		}
	}

	return nil
}

// ToSlice - returns slice of resources from the resource set.
func (resourceSet ResourceSet) ToSlice() []Resource {
	resources := []Resource{}
	for resource := range resourceSet {
		resources = append(resources, resource)
	}

	return resources
}

// Clone clones ResourceSet structure
func (resourceSet ResourceSet) Clone() ResourceSet {
	return NewResourceSet(resourceSet.ToSlice()...)
}

// BucketResourceExists - checks if at least one bucket resource exists in the set.
func (resourceSet ResourceSet) BucketResourceExists() bool {
	for resource := range resourceSet {
		if resource.isBucketPattern() {
			return true
		}
	}

	return false
}

// ObjectResourceExists - checks if at least one object resource exists in the set.
func (resourceSet ResourceSet) ObjectResourceExists() bool {
	for resource := range resourceSet {
		if resource.isObjectPattern() {
			return true
		}
	}

	return false
}

// NewResourceSet - creates new resource set.
func NewResourceSet(resources ...Resource) ResourceSet {
	resourceSet := make(ResourceSet)
	for _, resource := range resources {
		resourceSet.Add(resource)
	}

	return resourceSet
}
