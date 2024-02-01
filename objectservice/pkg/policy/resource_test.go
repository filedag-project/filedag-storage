package policy

import (
	"testing"
)

func TestResourceSet_Match(t *testing.T) {
	testCases := []struct {
		name           string
		bucketName     string
		keyName        string
		resource       string
		expectedResult bool
	}{
		{
			name:           "test1",
			bucketName:     "*",
			keyName:        "",
			resource:       "mybucket",
			expectedResult: true,
		},
		{
			name:           "test1",
			bucketName:     "*",
			keyName:        "",
			resource:       "mybucket/myobject",
			expectedResult: true,
		},
		{
			name:           "test2",
			bucketName:     "mybucket*",
			keyName:        "",
			resource:       "mybucket",
			expectedResult: true,
		},
		{
			name:           "test2",
			bucketName:     "mybucket*",
			keyName:        "",
			resource:       "mybucket/myobject",
			expectedResult: true,
		},
		{
			name:           "test3",
			bucketName:     "",
			keyName:        "*",
			resource:       "/mybucket",
			expectedResult: true,
		},
		{
			name:           "test3",
			bucketName:     "*",
			keyName:        "*",
			resource:       "mybucket/myobject",
			expectedResult: true,
		},
		{
			name:           "test3",
			bucketName:     "mybucket",
			keyName:        "*",
			resource:       "mybucket/myobject",
			expectedResult: true,
		},
		{
			name:           "test4",
			bucketName:     "mybucket*",
			keyName:        "/myobject",
			resource:       "mybucket/myobject",
			expectedResult: true,
		},
		{
			name:           "test4",
			bucketName:     "mybucket*",
			keyName:        "/myobject",
			resource:       "mybucket11/myobject",
			expectedResult: true,
		},
		{
			name:           "test5",
			bucketName:     "mybucket?0",
			keyName:        "/11/22/*",
			resource:       "mybucket20/11/22/1.jpg",
			expectedResult: true,
		},
		{
			name:           "test5",
			bucketName:     "mybucket?0",
			keyName:        "",
			resource:       "mybucket50",
			expectedResult: true,
		},
		{
			name:           "test6",
			bucketName:     "mybucket",
			keyName:        "",
			resource:       "mybucket",
			expectedResult: true,
		},
		{
			name:           "test6",
			bucketName:     "",
			keyName:        "*",
			resource:       "mybucket/myobject",
			expectedResult: false,
		},
		{
			name:           "test6",
			bucketName:     "*",
			keyName:        "*",
			resource:       "mybucket", //  */*
			expectedResult: false,
		},
		{
			name:           "test6",
			bucketName:     "mybucket",
			keyName:        "*",
			resource:       "mybucket11/myobject",
			expectedResult: false,
		}, {

			name:           "test6",
			bucketName:     "mybucket",
			keyName:        "",
			resource:       "mybucket/myobject",
			expectedResult: false,
		},
	}
	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			r := NewResourceSet(NewResource(testcase.bucketName, testcase.keyName))
			result := r.Match(testcase.resource, nil)
			if result != testcase.expectedResult {
				t.Errorf("expected %v, got %v", testcase.expectedResult, result)
			}
		})
	}
}
func TestResourceSet_Equals(t *testing.T) {
	testCases := []struct {
		name           string
		bucketName     string
		keyName        string
		bucketName1    string
		keyName1       string
		expectedResult bool
	}{
		{
			name:           "test1",
			bucketName:     "*",
			keyName:        "",
			bucketName1:    "*",
			keyName1:       "",
			expectedResult: true,
		},
		{
			name:           "test1",
			bucketName:     "*",
			keyName:        "*",
			bucketName1:    "*",
			keyName1:       "",
			expectedResult: false,
		},
		{
			name:           "test2",
			bucketName:     "",
			keyName:        "*",
			bucketName1:    "",
			keyName1:       "",
			expectedResult: false,
		},
		{
			name:           "test2",
			bucketName:     "",
			keyName:        "*",
			bucketName1:    "*",
			keyName1:       "*",
			expectedResult: false,
		},
		{
			name:           "test3",
			bucketName:     "*",
			keyName:        "",
			bucketName1:    "*",
			keyName1:       "*",
			expectedResult: false,
		},
		{
			name:           "test3",
			bucketName:     "*",
			keyName:        "",
			bucketName1:    "",
			keyName1:       "",
			expectedResult: false,
		},
		{
			name:           "test3",
			bucketName:     "*",
			keyName:        "*",
			bucketName1:    "",
			keyName1:       "",
			expectedResult: false,
		},
		{
			name:           "test3",
			bucketName:     "mybucket",
			keyName:        "abiu",
			bucketName1:    "mybucket",
			keyName1:       "abiu",
			expectedResult: true,
		},
	}
	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			r := NewResourceSet(NewResource(testcase.bucketName, testcase.keyName))
			result := r.Equals(NewResourceSet(NewResource(testcase.bucketName1, testcase.keyName1)))
			if result != testcase.expectedResult {
				t.Errorf("expected %v, got %v", testcase.expectedResult, result)
			}
		})
	}
}

func TestResourceSet_Validate(t *testing.T) {}
