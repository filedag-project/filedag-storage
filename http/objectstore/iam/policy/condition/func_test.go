package condition

import "testing"

func TestConditions_Evaluate(t *testing.T) {
	testCases := []struct {
		key            KeyName
		valuesToCheck  string // values to check if condition is satisfied
		values         map[string][]string
		expectedResult bool
	}{
		{
			key:            S3Prefix,
			valuesToCheck:  "object.txt",
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt"}},
			expectedResult: true,
		},
		{
			key:            S3Prefix,
			valuesToCheck:  "object.txt",
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt", "object2.txt"}},
			expectedResult: true,
		},
		{
			key:            S3Prefix,
			valuesToCheck:  "object2.txt",
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt", "object1.txt"}},
			expectedResult: false,
		},
		{
			key:            S3XAmzCopySource,
			valuesToCheck:  "object.txt",
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt"}},
			expectedResult: false,
		},
		{
			key:            S3XAmzCopySource,
			valuesToCheck:  "object.txt",
			values:         map[string][]string{S3XAmzCopySource.ToKey().Name(): {"object.txt"}},
			expectedResult: true,
		},
	}
	for i, testcase := range testCases {
		c, _ := NewStringEqualsFunc("", testcase.key.ToKey(), testcase.valuesToCheck)
		cf := NewConFunctions(c)
		eva := cf.Evaluate(testcase.values)
		if eva != testcase.expectedResult {
			t.Errorf("testcase %v should be %v", i, testcase.expectedResult)
		}
	}

}
func TestConditions_Equals(t *testing.T) {
	testCases := []struct {
		key            KeyName
		keyToCheck     KeyName // key to check if condition is satisfied
		values         string
		valuesToCheck  string // values to check if condition is satisfied
		expectedResult bool
	}{
		{
			key:            S3Prefix,
			keyToCheck:     S3Prefix,
			valuesToCheck:  "object.txt",
			values:         "object.txt",
			expectedResult: true,
		},
		{
			key:            S3Prefix,
			keyToCheck:     S3Prefix,
			valuesToCheck:  "object1.txt",
			values:         "object.txt",
			expectedResult: false,
		},
		{
			key:            S3Prefix,
			keyToCheck:     S3XAmzCopySource,
			valuesToCheck:  "object.txt",
			values:         "object.txt",
			expectedResult: false,
		},
		{
			key:            S3Prefix,
			keyToCheck:     S3XAmzCopySource,
			valuesToCheck:  "object.txt",
			values:         "object2.txt",
			expectedResult: false,
		},
	}
	for i, testcase := range testCases {
		c1, _ := NewStringEqualsFunc("", testcase.key.ToKey(), testcase.values)
		cf1 := NewConFunctions(c1)
		c2, _ := NewStringEqualsFunc("", testcase.keyToCheck.ToKey(), testcase.valuesToCheck)
		cf2 := NewConFunctions(c2)

		if cf1.Equals(cf2) != testcase.expectedResult {
			t.Errorf("testcase %v should be %v", i, testcase.expectedResult)
		}

	}
}
