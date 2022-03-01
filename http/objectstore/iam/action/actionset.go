package action

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
