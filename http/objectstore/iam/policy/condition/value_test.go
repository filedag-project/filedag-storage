package condition

import "testing"

func TestValue_getValuesByKey(t *testing.T) {
	testCases := []struct {
		key            KeyName
		values         map[string][]string
		expectedResult []string
	}{
		{
			key:            S3Prefix,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt"}},
			expectedResult: []string{"object.txt"},
		},
		{
			key:            S3Prefix,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt", "object2.txt"}},
			expectedResult: []string{"object.txt", "object2.txt"},
		},
		{
			key:            S3Prefix,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object2.txt", "object1.txt"}},
			expectedResult: []string{"object2.txt", "object1.txt"},
		},
		{
			key:            S3XAmzCopySource,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt"}},
			expectedResult: []string{},
		},
		{
			key:            S3XAmzCopySource,
			values:         map[string][]string{S3XAmzCopySource.ToKey().Name(): {"object.txt"}},
			expectedResult: []string{"object.txt"},
		},
	}
	for i, testcase := range testCases {
		values := getValuesByKey(testcase.values, testcase.key.ToKey())
		if len(values) != len(testcase.expectedResult) {
			t.Errorf("testcase %v should be %v", i, testcase.expectedResult)
		}
		for j, value := range values {
			if value != testcase.expectedResult[j] {
				t.Errorf("testcase %v should be %v", i, testcase.expectedResult)
			}
		}
	}
}
