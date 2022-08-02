package condition

import "testing"

func TestKey_IsValid(t *testing.T) {
	testcases := []struct {
		name           string
		key            Key
		expectedResult bool
	}{
		{
			name:           "test1",
			key:            NewKey(S3Prefix, ""),
			expectedResult: true,
		},
		{
			name:           "test1",
			key:            NewKey(S3Prefix, "aa"),
			expectedResult: true,
		},
		{
			name:           "test2",
			key:            NewKey(S3XAmzCopySource, ""),
			expectedResult: true,
		},
		{
			name:           "test3",
			key:            NewKey("unknown", ""),
			expectedResult: false,
		},
		{
			name:           "test3",
			key:            NewKey("unknown", "aa"),
			expectedResult: false,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.key.IsValid() != testcase.expectedResult {
				t.Errorf("testcase %v should be %v", testcase.name, testcase.expectedResult)
			}
		})
	}
}
func TestKey_parseKey(t *testing.T) {
	testcases := []struct {
		name           string
		keyString      string
		key            Key
		expectedResult bool
	}{
		{
			name:           "test1",
			keyString:      string(S3Prefix),
			key:            NewKey(S3Prefix, ""),
			expectedResult: true,
		},
		{
			name:           "test2",
			keyString:      string(S3XAmzCopySource),
			key:            NewKey(S3Prefix, ""),
			expectedResult: false,
		},
	}
	for i, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			key, err := parseKey(testcase.keyString)
			if err != nil {
				t.Errorf("testcases %v error:%v", i, err)
			}
			if (key.Name() == testcase.key.Name()) != testcase.expectedResult {
				t.Errorf("testcases %v key name should be %v", i, testcase.key.Name())
			}
		})
	}
}
