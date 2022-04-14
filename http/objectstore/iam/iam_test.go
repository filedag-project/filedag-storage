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
	uleveldb.DBClient, _ = uleveldb.OpenDb(path)
	var iamSys IdentityAMSys
	iamSys.Init()
	var accessKey = "test1"
	var secretKey = "test12345"

	//add user
	err := iamSys.AddUser(context.Background(), accessKey, secretKey)
	if err != nil {
		fmt.Println(err)
	}

	//list user
	users, err := iamSys.GetUserList(context.Background(), "")
	if err != nil {
		fmt.Println(users)
	}

	//get user
	cred, bool := iamSys.GetUser(context.Background(), accessKey)
	fmt.Println(bool, cred)

	//remove user
	err = iamSys.RemoveUser(context.Background(), accessKey)
	if err != nil {
		fmt.Println(err)
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
	uleveldb.DBClient, _ = uleveldb.OpenDb(path)
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
