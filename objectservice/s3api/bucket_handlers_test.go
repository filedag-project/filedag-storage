package s3api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/objectservice/iamapi"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var w *httptest.ResponseRecorder
var router = mux.NewRouter()

func TestMain(m *testing.M) {
	db, err := uleveldb.OpenDb((&testing.T{}).TempDir())
	if err != nil {
		println(err)
		return
	}
	defer db.Close()
	cred, err := auth.CreateCredentials(auth.DefaultAccessKey, auth.DefaultSecretKey)
	if err != nil {
		println(err)
		return
	}
	authSys := iam.NewAuthSys(db, cred)
	poolCli, done := client.NewMockPoolClient(&testing.T{})
	defer done()
	dagServ := merkledag.NewDAGService(blockservice.New(poolCli, offline.Exchange(poolCli)))
	storageSys := store.NewStorageSys(context.TODO(), dagServ, db)
	bmSys := store.NewBucketMetadataSys(db)
	storageSys.SetNewBucketNSLock(bmSys.NewNSLock)
	storageSys.SetHasBucket(bmSys.HasBucket)
	bmSys.SetEmptyBucket(storageSys.EmptyBucket)
	cleanData := func(accessKey string) {
		ctx := context.Background()
		bkts, err := bmSys.GetAllBucketsOfUser(ctx, accessKey)
		if err != nil {
			log.Errorf("GetAllBucketsOfUser error: %v", err)
		}
		for _, bkt := range bkts {
			if err = storageSys.CleanObjectsInBucket(ctx, bkt.Name); err != nil {
				log.Errorf("CleanObjectsInBucket error: %v", err)
				continue
			}
			if err = bmSys.DeleteBucket(ctx, bkt.Name); err != nil {
				log.Errorf("DeleteBucket error: %v", err)
			}
		}
	}
	iamapi.NewIamApiServer(router, authSys, cleanData)
	NewS3Server(router, authSys, bmSys, storageSys)
	os.Exit(m.Run())
}
func reqTest(r *http.Request) *httptest.ResponseRecorder {
	// mock a response logger
	w = httptest.NewRecorder()
	// Let the server process the mock request and record the returned response content
	router.ServeHTTP(w, r)
	return w
}

func TestS3ApiServer_PutBucketHandler(t *testing.T) {
	bucketName := "/testbucketput"
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName:         bucketName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		// mock an HTTP request
		reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		result := reqTest(reqPutBucket)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
	}

}
func TestS3ApiServer_HeadBucketHandler(t *testing.T) {
	bucketName := "/testbuckethead"
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName:         bucketName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		// Test case - 3.
		// wrong accessKey.
		{
			bucketName:         "/1",
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		// mock an HTTP request
		reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		result1 := reqTest(reqPutBucket)
		if result1.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
		}

		reqHeadBucket := utils.MustNewSignedV4Request(http.MethodHead, testCase.bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		result2 := reqTest(reqHeadBucket)
		if result2.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result2.Code)
		}
	}

}
func TestS3ApiServer_ListBucketHandler(t *testing.T) {
	bucketName := "/testbucketlist"
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		isPut     bool
		accessKey string
		secretKey string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{

		// Test case - 1.
		// wrong accessKey.
		{
			isPut:              false,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			isPut:              true,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 3.
		// wrong accessKey.
		{
			isPut:              true,
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		// mock an HTTP request
		if testCase.isPut {
			reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result1 := reqTest(reqPutBucket)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
			}
		}
		reqListBucket := utils.MustNewSignedV4Request(http.MethodGet, "/", 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		result := reqTest(reqListBucket)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		var resp response.ListAllMyBucketsResult
		utils.XmlDecoder(result.Body, &resp, reqListBucket.ContentLength)
		fmt.Printf("case:%v  list:%v\n", i+1, resp)

		//reqDeleteBucket := testsign.MustNewSignedV4Request(http.MethodDelete, "/testbucket", 0,
		//	nil, "s3", testCase.accessKey, testCase.secretKey, t)
		//result4 := reqTest(reqDeleteBucket)
		//if result4.Code != testCase.expectedRespStatus {
		//	t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result4.Code)
		//}
		//
		//resp2 := response.ListAllMyBucketsResult{}
		//utils.XmlDecoder(reqTest(reqListBucket).Body, &resp2, reqListBucket.ContentLength)
		//fmt.Printf("case:%v  list:%v\n", i+1,resp)
	}

}
func TestS3ApiServer_DeleteBucketHandler(t *testing.T) {
	bucketName := "/testbucketdel"
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		isPut     bool
		accessKey string
		secretKey string
		// expected output.
		expectedPutRespStatus int // expected response status body.
		expectedDelRespStatus int
	}{

		// Test case - 1.
		// wrong accessKey.
		{
			isPut:                 false,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusNotFound,
			expectedDelRespStatus: http.StatusNotFound,
		},

		// Test case - 2.
		// wrong accessKey.
		{
			isPut:                 true,
			accessKey:             "1",
			secretKey:             "1",
			expectedPutRespStatus: http.StatusForbidden,
			expectedDelRespStatus: http.StatusForbidden,
		},
		// Test case - 3.
		// Fetching the entire Bucket and validating its contents.
		{
			isPut:                 true,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusOK,
			expectedDelRespStatus: http.StatusNoContent,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		// mock an HTTP request
		if testCase.isPut {
			reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result1 := reqTest(reqPutBucket)
			if result1.Code != testCase.expectedPutRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedPutRespStatus, result1.Code)
			}
		}

		reqDeleteBucket := utils.MustNewSignedV4Request(http.MethodDelete, "/testbucketdel", 0,
			nil, "s3", testCase.accessKey, testCase.secretKey, t)
		result := reqTest(reqDeleteBucket)
		if result.Code != testCase.expectedDelRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedDelRespStatus, result.Code)
		}
		reqListBucket := utils.MustNewSignedV4Request(http.MethodGet, "/", 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)

		resp1 := response.ListAllMyBucketsResult{}
		utils.XmlDecoder(reqTest(reqListBucket).Body, &resp1, reqListBucket.ContentLength)
		fmt.Printf("case:%v  list:%v\n", i+1, resp1)
	}

}

func TestS3ApiServer_BucketPolicyHandler(t *testing.T) {
	u := "/testbucketpoliy"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())

	p := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucket/*"}]}`
	reqPut := utils.MustNewSignedV4Request(http.MethodPut, u+"?policy", int64(len(p)), strings.NewReader(p),
		"s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("put:", reqTest(reqPut).Body.String())

	reqGet := utils.MustNewSignedV4Request(http.MethodGet, u+"?policy", 0, nil, "s3",
		DefaultTestAccessKey, DefaultTestSecretKey, t)
	resp1 := policy.Policy{}
	json.Unmarshal([]byte(reqTest(reqGet).Body.String()), &resp1)
	fmt.Println("get:", resp1)

	reqDel := utils.MustNewSignedV4Request(http.MethodDelete, u+"?policy", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("del:", reqTest(reqDel).Body.String())
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
/*func TestS3_PutBucketHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test22"

	req := testsign.MustNewSignedV4Request(http.MethodPut, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	//req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}*/
