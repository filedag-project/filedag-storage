package s3action

import "testing"

func TestActionSet_Equals(t *testing.T) {
	testCases := []struct {
		name           string
		actionSet      ActionSet
		actionSet1     ActionSet
		expectedResult bool
	}{
		{
			name:           "test1",
			actionSet:      NewActionSet(PutObjectAction),
			actionSet1:     NewActionSet(PutObjectAction),
			expectedResult: true,
		},
		{
			name:           "test1",
			actionSet:      NewActionSet("*"),
			actionSet1:     NewActionSet(PutObjectAction),
			expectedResult: false,
		},
		{
			name:           "test1",
			actionSet:      SupportedActions,
			actionSet1:     NewActionSet(PutObjectAction),
			expectedResult: false,
		},
	}
	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {
			if testCase.actionSet.Equals(testCase.actionSet1) != testCase.expectedResult {
				t.Errorf("Test case failed: %s", testCase.name)
			}
		})
	}
}
func TestActionSet_Validate(t *testing.T) {
	testCases := []struct {
		name           string
		actionSet      ActionSet
		expectedResult bool
	}{
		{
			name:           "test1",
			actionSet:      NewActionSet(PutObjectAction),
			expectedResult: true,
		},
		{
			name:           "test1",
			actionSet:      NewActionSet("*"),
			expectedResult: true,
		},
		{
			name:           "test1",
			actionSet:      SupportedActions,
			expectedResult: true,
		},
		{
			name:           "test2",
			actionSet:      NewActionSet("abcd"),
			expectedResult: false,
		},
		{
			name:           "test2",
			actionSet:      NewActionSet(PutObjectAction + "*"),
			expectedResult: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if (testCase.actionSet.Validate() == nil) != testCase.expectedResult {
				t.Errorf("Test case failed: %s", testCase.name)
			}
		})
	}
}
