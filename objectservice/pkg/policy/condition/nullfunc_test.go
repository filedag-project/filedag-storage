package condition

import "testing"

func TestNullFunc_evaluate(t *testing.T) {
	testCases := []struct {
		name           string
		key            KeyName
		valuesToCheck  bool // values to check if condition is satisfied
		values         bool
		expectedResult bool
	}{
		{
			name:           "test1",
			key:            S3Prefix,
			valuesToCheck:  true,
			values:         true,
			expectedResult: false,
		},
		{
			name:           "test1",
			key:            S3Prefix,
			valuesToCheck:  true,
			values:         false,
			expectedResult: true,
		},
		{
			name:           "test2",
			key:            S3Prefix,
			valuesToCheck:  false,
			values:         true,
			expectedResult: true,
		},
		{
			name:           "test3",
			key:            S3Prefix,
			valuesToCheck:  false,
			values:         false,
			expectedResult: false,
		},
	}
	for i, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			c, err := newNullFunc(testcase.key.ToKey(), NewValueSet(NewBoolValue(testcase.valuesToCheck)), "")
			if err != nil {
				t.Errorf("error creating null func: %v", err)
			}
			m := make(map[string][]string)
			if testcase.values {
				m[testcase.key.Name()] = []string{""} // values to check if condition is satisfied
			}
			if c.evaluate(m) != testcase.expectedResult {
				t.Errorf("testcase %v should be %v", i, testcase.expectedResult)
			}
		})
	}
}
