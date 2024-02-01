package store

import (
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"testing"
)

func TestGetPolicyName(t *testing.T) {
	testCases := []struct {
		name       string
		statements []policy.Statement
		bucketName string
		prefix     string
		expect     BucketPolicy
	}{
		{
			name: bucketPolicyReadWrite,
			statements: []policy.Statement{
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("filedagadmin"),
					Actions:    s3action.NewActionSet(commonBucketActions),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("*"),
					Actions:    s3action.NewActionSet(append(writeOnlyObjectActions.ToSlice(), commonBucketActions, writeOnlyBucketActions, readOnlyBucketActions, readOnlyObjectActions)...),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "*"), policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
			},
			bucketName: "testbucket",
			expect:     bucketPolicyReadWrite,
		},
		{
			name: bucketPolicyReadOnly,
			statements: []policy.Statement{
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("filedagadmin"),
					Actions:    s3action.NewActionSet(commonBucketActions),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("*"),
					Actions:    s3action.NewActionSet(commonBucketActions, writeOnlyBucketActions, readOnlyBucketActions, readOnlyObjectActions),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "*"), policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
			},
			bucketName: "testbucket",
			expect:     bucketPolicyReadOnly,
		},
		{
			name: bucketPolicyWriteOnly,
			statements: []policy.Statement{
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("filedagadmin"),
					Actions:    s3action.NewActionSet(commonBucketActions),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("*"),
					Actions:    s3action.NewActionSet(append(writeOnlyObjectActions.ToSlice(), commonBucketActions, writeOnlyBucketActions)...),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "*"), policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
			},
			bucketName: "testbucket",
			expect:     bucketPolicyWriteOnly,
		},
		{
			name: string(bucketPolicyNone),
			statements: []policy.Statement{
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("filedagadmin"),
					Actions:    s3action.NewActionSet(commonBucketActions),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
				{
					SID:        "",
					Effect:     policy.Allow,
					Principal:  policy.NewPrincipal("*"),
					Actions:    s3action.NewActionSet(commonBucketActions),
					Resources:  policy.NewResourceSet(policy.NewResource("testbucket", "*"), policy.NewResource("testbucket", "")),
					Conditions: nil,
				},
			},
			bucketName: "testbucket",
			expect:     bucketPolicyNone,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			re := GetPolicyName(testCase.statements, testCase.bucketName, testCase.prefix)
			if re != testCase.expect {
				t.Fatalf("%v,expect %v,get %v ", testCase.name, testCase.expect, re)
			}
		})
	}

}
