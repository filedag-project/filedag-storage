package store

//import (
//	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
//	"testing"
//)
//
//func TestGetPolicyName(t *testing.T) {
//	testCases := []struct {
//		name       string
//		statements []policy.Statement
//		bucketName string
//		prefix     string
//		expect     BucketPolicy
//	}{{
//		name: "1",
//		statements: []policy.Statement{{
//			SID:        "",
//			Effect:     policy.Allow,
//			Principal:  policy.NewPrincipal("filedagadmin"),
//			Actions:    commonBucketActions,
//			Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "")),
//			Conditions: nil,
//		}},
//		bucketName: "testbucket",
//		expect:     BucketPolicyNone,
//	}}
//	for _, testCase := range testCases {
//		re := GetPolicyName(testCase.statements, testCase.bucketName, testCase.prefix)
//		if re != testCase.expect {
//			t.Fatalf("%v,expect %v,get %v ", testCase.name, testCase.expect, re)
//		}
//	}
//
//}
