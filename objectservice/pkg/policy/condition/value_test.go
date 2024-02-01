package condition

import "testing"

func TestValue_getValuesByKey(t *testing.T) {
	testCases := []struct {
		name           string
		key            KeyName
		values         map[string][]string
		expectedResult []string
	}{
		{
			name:           "test1",
			key:            S3Prefix,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt"}},
			expectedResult: []string{"object.txt"},
		},
		{
			name:           "test1",
			key:            S3Prefix,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt", "object2.txt"}},
			expectedResult: []string{"object.txt", "object2.txt"},
		},
		{
			name:           "test1",
			key:            S3Prefix,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object2.txt", "object1.txt"}},
			expectedResult: []string{"object2.txt", "object1.txt"},
		},
		{
			name:           "test2",
			key:            S3XAmzCopySource,
			values:         map[string][]string{S3Prefix.ToKey().Name(): {"object.txt"}},
			expectedResult: []string{},
		},
		{
			name:           "test2",
			key:            S3XAmzCopySource,
			values:         map[string][]string{S3XAmzCopySource.ToKey().Name(): {"object.txt"}},
			expectedResult: []string{"object.txt"},
		},
	}
	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			values := getValuesByKey(testcase.values, testcase.key.ToKey())
			if len(values) != len(testcase.expectedResult) {
				t.Errorf("testcase %v should be %v", testcase.name, testcase.expectedResult)
			}
			for j, value := range values {
				if value != testcase.expectedResult[j] {
					t.Errorf("testcase %v should be %v", testcase.name, testcase.expectedResult)
				}
			}
		})
	}
}
