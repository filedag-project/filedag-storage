package action

import (
	"errors"
	"fmt"
)

type Action string

// ActionSet - set of actions.
type ActionSet map[Action]struct{}

// Add - add action to the set.
func (actionSet ActionSet) Add(action Action) {
	actionSet[action] = struct{}{}
}

// Contains - checks given action exists in the action set.
func (actionSet ActionSet) Contains(action Action) bool {
	_, found := actionSet[action]
	return found
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
func (actionSet ActionSet) Equals(sactionSet ActionSet) bool {
	// If length of set is not equal to length of given set, the
	// set is not equal to given set.
	if len(actionSet) != len(sactionSet) {
		return false
	}

	// As both sets are equal in length, check each elements are equal.
	for k := range actionSet {
		if _, ok := sactionSet[k]; !ok {
			return false
		}
	}

	return true
}

// Clone clones ActionSet structure
func (actionSet ActionSet) Clone() ActionSet {
	return NewActionSet(actionSet.ToSlice()...)
}

// ToSlice - returns slice of actions from the action set.
func (actionSet ActionSet) ToSlice() []Action {
	var actions []Action
	for action := range actionSet {
		actions = append(actions, action)
	}
	return actions
}

// Validate checks if all actions are valid
func (actionSet ActionSet) Validate() error {
	for _, action := range actionSet.ToSlice() {
		if !action.IsValid() {
			return errors.New(fmt.Sprintf("unsupported action '%v'", action))
		}
	}
	return nil
}
