package s3api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var w *httptest.ResponseRecorder
var router = mux.NewRouter()

func TestMain(m *testing.M) {
	var err error
	uleveldb.DBClient, err = uleveldb.OpenDb("./test")
	if err != nil {
		return
	}
	defer uleveldb.DBClient.Close()
	iamapi.NewIamApiServer(router)
	NewS3Server(router)
	os.Exit(m.Run())
}
func reqTest(r *http.Request) *bytes.Buffer {
	// mock a response logger
	w = httptest.NewRecorder()
	// Let the server process the mock request and record the returned response content
	router.ServeHTTP(w, r)
	return w.Body
}
func TestS3ApiServer_BucketHandler(t *testing.T) {
	// mock an HTTP request
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, "/testbucket", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("put:", reqTest(reqPutBucket).String())

	reqHeadBucket := testsign.MustNewSignedV4Request(http.MethodHead, "/testbucket", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("head:", reqTest(reqHeadBucket).String())

	reqListBucket := testsign.MustNewSignedV4Request(http.MethodGet, "/", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	var resp1 response.ListAllMyBucketsResult
	utils.XmlDecoder(reqTest(reqListBucket), &resp1, reqListBucket.ContentLength)
	fmt.Println("list:", resp1)

	reqDeleteBucket := testsign.MustNewSignedV4Request(http.MethodDelete, "/testbucket", 0,
		nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("delete:", reqTest(reqDeleteBucket).String())

	resp2 := response.ListAllMyBucketsResult{}
	utils.XmlDecoder(reqTest(reqListBucket), &resp2, reqListBucket.ContentLength)
	fmt.Println("list:", resp2)
}

func TestS3ApiServer_BucketPolicyHandler(t *testing.T) {
	u := "/testbucket"
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).String())

	p := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucket/*"}]}`
	reqPut := testsign.MustNewSignedV4Request(http.MethodPut, u+"?policy", int64(len(p)), strings.NewReader(p),
		"s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("put:", reqTest(reqPut).String())

	reqGet := testsign.MustNewSignedV4Request(http.MethodGet, u+"?policy", 0, nil, "s3",
		DefaultTestAccessKey, DefaultTestSecretKey, t)
	resp1 := policy.Policy{}
	json.Unmarshal([]byte(reqTest(reqGet).String()), &resp1)
	fmt.Println("get:", resp1)

	reqDel := testsign.MustNewSignedV4Request(http.MethodDelete, u+"?policy", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("del:", reqTest(reqDel).String())
}

//func TestS3ApiServer_GetBucketLocationHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test22"
//	//req.Header.Set("Content-Type", "text/plain")
//	req := testsign.MustNewSignedV4Request(http.MethodHead, u+"?location", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//
//func TestS3ApiServer_GetBucketAclHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?acl=", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//func TestS3ApiServer_PutBucketAclHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test"
//	a := `<?xml version="1.0" encoding="UTF-8"?>
//<AccessControlPolicy xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
//  <Owner>
//    <ID>*** Owner-Canonical-User-ID ***</ID>
//    <DisplayName>owner-display-name</DisplayName>
//  </Owner>
//  <AccessControlList>
//    <Grant>
//      <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
//               xsi:type="Canonical User">
//        <ID>*** Owner-Canonical-User-ID ***</ID>
//        <DisplayName>display-name</DisplayName>
//      </Grantee>
//      <Permission>FULL_CONTROL</Permission>
//    </Grant>
//  </AccessControlList>
//</AccessControlPolicy>`
//	req := testsign.MustNewSignedV4Request(http.MethodPut, u+"?acl=", int64(len(a)), strings.NewReader(a), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
