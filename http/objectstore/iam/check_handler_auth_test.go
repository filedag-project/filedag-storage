package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"testing"
)

func TestV2CheckRequestAuthType(t *testing.T) {
	var aSys AuthSys
	aSys.Init()
	req := testsign.MustNewSignedV2Request("GET", "http://127.0.0.1:9000", 0, nil, t)
	err := aSys.checkRequestAuthType(context.Background(), req, s3action.ListAllMyBucketsAction, "test", "testobject")
	fmt.Println(api_errors.GetAPIError(err))
}
func TestV4CheckRequestAuthType(t *testing.T) {
	var aSys AuthSys
	aSys.Init()
	req := testsign.MustNewSignedV4Request("GET", "http://127.0.0.1:9000", 0, nil, "test", "test", "s3", t)
	err := aSys.checkRequestAuthType(context.Background(), req, s3action.ListAllMyBucketsAction, "test", "testobject")
	fmt.Println(api_errors.GetAPIError(err))
}
