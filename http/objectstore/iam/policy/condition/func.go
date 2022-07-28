package condition

// CondFunction - condition function interface.
type CondFunction interface {
	// evaluate() - evaluates this condition function with given values.
	evaluate(values map[string][]string) bool

	// String() - returns string representation of function.
	String() string
}

// Conditions - list of functions.
type Conditions []CondFunction
