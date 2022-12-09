package s3action

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/set"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/xerrors"
	"sort"
)

//Action s3 action
type Action string

// ActionSet - set of actions.
type ActionSet map[Action]struct{}

// Add - add action to the set.
func (as ActionSet) Add(action Action) {
	as[action] = struct{}{}
}

// Match - matches object name with anyone of action pattern in action set.
func (as ActionSet) Match(action Action) bool {
	for r := range as {
		if r.Match(action) {
			return true
		}

		// This is a special case where GetObjectVersion
		// means GetObject is enabled implicitly.
		switch r {
		case GetObjectVersionAction:
			if action == GetObjectAction {
				return true
			}
		}
	}

	return false
}

// NewActionSet - creates new action set.
func NewActionSet(actions ...Action) ActionSet {
	actionSet := make(ActionSet)
	for _, action := range actions {
		actionSet.Add(action)
	}

	return actionSet
}

// Equals - checks whether given action set is equal to current action set or not.
func (as ActionSet) Equals(sactionSet ActionSet) bool {
	// If length of set is not equal to length of given set, the
	// set is not equal to given set.
	if len(as) != len(sactionSet) {
		return false
	}

	// As both sets are equal in length, check each elements are equal.
	for k := range as {
		if _, ok := sactionSet[k]; !ok {
			return false
		}
	}

	return true
}

// Clone clones ActionSet structure
func (as ActionSet) Clone() ActionSet {
	return NewActionSet(as.ToSlice()...)
}

// ToSlice - returns slice of actions from the action set.
func (as ActionSet) ToSlice() []Action {
	var actions []Action
	for action := range as {
		actions = append(actions, action)
	}
	return actions
}

// Validate checks if all actions are valid
func (as ActionSet) Validate() error {
	for _, action := range as.ToSlice() {
		if !action.IsValid() {
			return errors.New(fmt.Sprintf("unsupported action '%v'", action))
		}
	}
	return nil
}

// MarshalJSON - encodes ActionSet to JSON data.
func (as ActionSet) MarshalJSON() ([]byte, error) {
	if len(as) == 0 {
		return nil, errors.New("empty actions not allowed")
	}

	return json.Marshal(as.ToSlice())
}

func (as ActionSet) MarshalMsgpack() ([]byte, error) {
	if len(as) == 0 {
		return nil, errors.New("empty actions not allowed")
	}

	return msgpack.Marshal(as.ToSlice())
}

func (as ActionSet) String() string {
	var actions []string
	for action := range as {
		actions = append(actions, string(action))
	}
	sort.Strings(actions)

	return fmt.Sprintf("%v", actions)
}

// UnmarshalJSON - decodes JSON data to ActionSet.
func (as *ActionSet) UnmarshalJSON(data []byte) error {
	var sset set.StringSet
	if err := json.Unmarshal(data, &sset); err != nil {
		return err
	}

	if len(sset) == 0 {
		return errors.New("empty actions not allowed")
	}

	*as = make(ActionSet)
	for _, s := range sset.ToSlice() {
		action := Action(s)
		if action.IsValid() {
			as.Add(action)
		} else {
			return xerrors.Errorf("unsupported action '%v'", s)
		}
	}

	return nil
}

func (as *ActionSet) UnmarshalMsgpack(data []byte) error {
	var sset set.StringSet
	if err := msgpack.Unmarshal(data, &sset); err != nil {
		return err
	}

	if len(sset) == 0 {
		return errors.New("empty actions not allowed")
	}

	*as = make(ActionSet)
	for _, s := range sset.ToSlice() {
		action := Action(s)
		if action.IsValid() {
			as.Add(action)
		} else {
			return xerrors.Errorf("unsupported action '%v'", s)
		}
	}

	return nil
}
