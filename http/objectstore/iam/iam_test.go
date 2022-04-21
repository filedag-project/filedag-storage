package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"testing"
)

func TestIdentityAMSys_UserApi(t *testing.T) {
	uleveldb.DBClient, _ = uleveldb.OpenDb(tmppath)
	var iamSys IdentityAMSys
	iamSys.Init()
	//var accessKey = "test1"
	//var secretKey = "test12345"
	ctx := context.Background()
	testCases := []struct {
		isRemove  bool
		accessKey string
		secretKey string
	}{
		// Test case - 1.
		// Fetching the entire User and validating its contents.
		{
			isRemove:  true,
			accessKey: "adminTest1",
			secretKey: "adminTest1",
		},
		// Test case - 2.
		// wrong The same user name already exists ..
		{
			isRemove:  true,
			accessKey: "adminTest2",
			secretKey: "adminTest2",
		},
		// Test case - 3.
		// error  access key length should be between 3 and 20.
		{
			isRemove:  false,
			accessKey: "1",
			secretKey: "test1234",
		},
		// Test case - 4.
		// error  secret key length should be between 3 and 20.
		{
			isRemove:  false,
			accessKey: "test2",
			secretKey: "1",
		},
	}
	//add user
	fmt.Println("********add user********")
	for _, testCase := range testCases {
		if testCase.isRemove {
			//remove user
			err := iamSys.RemoveUser(ctx, testCase.accessKey)
			if err != nil {
				fmt.Println(err)
			}
		}
		//add user
		err := iamSys.AddUser(ctx, testCase.accessKey, testCase.secretKey)
		if err != nil {
			fmt.Println(err)
		}
	}

	//get user
	fmt.Println("********list user********")
	for _, testCase := range testCases {
		if testCase.isRemove {
			//get user
			cred, bool := iamSys.GetUser(ctx, testCase.accessKey)
			if bool {
				fmt.Println(bool, cred)
			}
		}
	}

	//list user
	fmt.Println("********list user********")
	users, err := iamSys.GetUserList(ctx, "")
	if err != nil {
		fmt.Println(users)
	}
	bytes, err := json.Marshal(users)
	fmt.Println(string(bytes))

	//remove user
	fmt.Println("********remove user********")
	for _, testCase := range testCases {
		//remove user
		err := iamSys.RemoveUser(ctx, testCase.accessKey)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//GetPolicyDocument Get PolicyDocument
func GetPolicyDocument(policyD *string) (policyDocument policy.PolicyDocument, err error) {
	if err = json.Unmarshal([]byte(*policyD), &policyDocument); err != nil {
		return policy.PolicyDocument{}, err
	}
	return policyDocument, err
}

func TestIdentityAMSys_UserPolicyApi(t *testing.T) {
	uleveldb.DBClient, _ = uleveldb.OpenDb(tmppath)
	var iamSys IdentityAMSys
	iamSys.Init()
	var userName = "test1"
	var policyName = "read2"

	policyDocumentString := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test1/*"}]}`
	//put user policy
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		log.Errorf("GetPolicyDocument:%v", err)
	}
	err = iamSys.PutUserPolicy(context.Background(), userName, policyName, policyDocument)
	if err != nil {
		log.Errorf("PutUserPolicy:%v", err)
	}

	//list user policy
	users, err := iamSys.GetUserPolices(context.Background(), userName)
	if err != nil {
		log.Errorf("GetUserPolices:%v", err)
	}
	log.Info(users)
}
