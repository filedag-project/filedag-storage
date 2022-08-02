package condition

import "testing"

func TestKey_IsValid(t *testing.T) {
	testcases := []struct {
		key            Key
		expectedResult bool
	}{
		{
			key:            NewKey(S3Prefix, ""),
			expectedResult: true,
		},
		{
			key:            NewKey(S3Prefix, "aa"),
			expectedResult: true,
		},
		{
			key:            NewKey(S3XAmzCopySource, ""),
			expectedResult: true,
		},
		{
			key:            NewKey("unknown", ""),
			expectedResult: false,
		},
		{
			key:            NewKey("unknown", "aa"),
			expectedResult: false,
		},
	}
	for i, testcase := range testcases {
		if testcase.key.IsValid() != testcase.expectedResult {
			t.Errorf("testcase %v should be %v", i, testcase.expectedResult)
		}
	}
}
func TestKey_parseKey(t *testing.T) {
	testcases := []struct {
		keyString      string
		key            Key
		expectedResult bool
	}{
		{
			keyString:      string(S3Prefix),
			key:            NewKey(S3Prefix, ""),
			expectedResult: true,
		},
		{
			keyString:      string(S3XAmzCopySource),
			key:            NewKey(S3Prefix, ""),
			expectedResult: false,
		},
	}
	for i, testcase := range testcases {
		key, err := parseKey(testcase.keyString)
		if err != nil {
			t.Errorf("testcases %v error:%v", i, err)
		}
		if (key.Name() == testcase.key.Name()) != testcase.expectedResult {
			t.Errorf("testcases %v key name should be %v", i, testcase.key.Name())
		}
	}
}
