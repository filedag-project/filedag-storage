package iamapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"net/http"
	"net/url"
	"testing"
)

func TestIamApiServer_AccountInfos(t *testing.T) {
	// test cases with inputs and expected result for UserInfo.
	testCases := []struct {
		name               string
		credAccessKey      string
		credSecretKey      string
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user get himself info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user get himself info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
	}
	u := "http://127.0.0.1:9985/admin/v1/user-infos"
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.credAccessKey, func(t *testing.T) {
			//user info
			userinfoReq := utils.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result1 := reqTest(userinfoReq)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result1.Code)
			}
			fmt.Println(result1.Body.String())
		})

	}
}
func TestStorePoolStats(t *testing.T) {
	testCases := []struct {
		name               string
		credAccessKey      string
		credSecretKey      string
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user get  info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user get  info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
	}
	u := "http://127.0.0.1:9985/admin/v1/store-pool-stats"
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.credAccessKey, func(t *testing.T) {
			//store-pool-stats
			poolStatsReq := utils.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result1 := reqTest(poolStatsReq)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result1.Code)
			}
			//fmt.Println(result1.Body.String())
		})

	}

}
func TestIamApiServer_request_overview(t *testing.T) {
	// test cases with inputs and expected result for UserInfo.
	testCases := []struct {
		name               string
		credAccessKey      string
		credSecretKey      string
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user get himself info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "root user get himself info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
	}
	u := "http://127.0.0.1:9985/admin/v1/request-overview"
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.credAccessKey, func(t *testing.T) {
			//user info
			userinfoReq := utils.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result1 := reqTest(userinfoReq)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result1.Code)
			}
			fmt.Println(result1.Body.String())
		})

	}
}

//change password
func TestIamApiServer_ChangePassword(t *testing.T) {
	changePassSuccessAccessKey := "changePassSuccess"
	changePassSuccessSecretKey := "changePassSuccess"
	himselfChangeSuccessAccessKey := "himselfChangeSuccess"
	himselfChangeSuccessSecretKey := "himselfChangeSuccess"

	thePassToChange := "thePassToChange"
	reqPutUserChangePassSuccessUrl := addUserUrl(changePassSuccessAccessKey, changePassSuccessSecretKey, defaultCap)
	reqPutUserChangePassSuccess := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserChangePassSuccessUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultChangePassSucces := reqTest(reqPutUserChangePassSuccess)
	if resultChangePassSucces.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultChangePassSucces.Code, resultChangePassSucces.Body.String())
	}
	reqPutUserHimselfUrl := addUserUrl(himselfChangeSuccessAccessKey, himselfChangeSuccessSecretKey, defaultCap)
	reqPutUserHimself := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserHimselfUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultHimself := reqTest(reqPutUserHimself)
	if resultHimself.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultHimself.Code, resultHimself.Body.String())
	}

	testCases := []struct {
		name          string
		credAccessKey string
		credSecretKey string
		oldSecretKey  string
		accessKey     string
		pass          string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user change a user pass",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          changePassSuccessAccessKey,
			oldSecretKey:       changePassSuccessSecretKey,
			expectedRespStatus: http.StatusOK,
			pass:               thePassToChange,
		},
		{
			name:               "root user change root pass",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			oldSecretKey:       DefaultTestSecretKey,
			accessKey:          DefaultTestAccessKey,
			expectedRespStatus: http.StatusConflict,
			pass:               thePassToChange,
		},
		{
			name:               "root user change a non-exist user pass",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          userNonExistAccessKey,
			oldSecretKey:       userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change himself pass",
			credAccessKey:      himselfChangeSuccessAccessKey,
			credSecretKey:      himselfChangeSuccessSecretKey,
			accessKey:          himselfChangeSuccessAccessKey,
			oldSecretKey:       himselfChangeSuccessSecretKey,
			expectedRespStatus: http.StatusOK,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change other user pass",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			oldSecretKey:       otherUserSecretKey,
			accessKey:          otherUserAccessKey,
			expectedRespStatus: http.StatusForbidden,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change a non-exist user pass",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			oldSecretKey:       userNonExistSecretKey,
			accessKey:          userNonExistAccessKey,
			expectedRespStatus: http.StatusConflict,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change user err pass",
			credAccessKey:      himselfChangeSuccessAccessKey,
			credSecretKey:      thePassToChange,
			oldSecretKey:       thePassToChange,
			accessKey:          himselfChangeSuccessAccessKey,
			expectedRespStatus: http.StatusBadRequest,
			pass:               "dj",
		},
	}
	changePassUrl := "http://127.0.0.1:9985/console/v1/change-password?"
	for _, testCase := range testCases {
		// mock an HTTP request
		//change password
		t.Run(testCase.name, func(t *testing.T) {
			urlValues := make(url.Values)
			urlValues.Set(newSecretKey, testCase.pass)
			urlValues.Set(accessKey, testCase.accessKey)
			urlValues.Set(oldSecretKey, testCase.oldSecretKey)
			//urlValues.Set("status", string(iam.AccountDisabled
			reqChange := utils.MustNewSignedV4Request(http.MethodPost, changePassUrl+urlValues.Encode(), 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result := reqTest(reqChange)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})
	}
}

// set status
func TestIamApiServer_SetStatus(t *testing.T) {
	offUserAccessKey := "offUser"
	offUserSecretKey := "offUser1234"
	otherUserOffAccessKey := "otherUserOffUser"
	otherUserOffSecretKey := "otherUserOffUser1234"
	reqPutUserOffUrl := addUserUrl(offUserAccessKey, offUserSecretKey, defaultCap)
	reqPutUserOff := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserOffUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultOff := reqTest(reqPutUserOff)
	if resultOff.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultOff.Code, resultOff.Body.String())
	}
	reqPutUserOtherOffUrl := addUserUrl(otherUserOffAccessKey, otherUserOffSecretKey, defaultCap)
	reqPutUserOtherOff := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserOtherOffUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultOtherOff := reqTest(reqPutUserOtherOff)
	if resultOtherOff.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultOtherOff.Code, resultOtherOff.Body.String())
	}
	testCases := []struct {
		name          string
		credAccessKey string
		credSecretKey string
		accessKey     string
		secretKey     string
		status        string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user set a user off",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusOK,
			status:             "off",
		},
		{
			name:               "root user set a user on",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusOK,
			status:             "on",
		},
		{
			name:               "root user set a non exist user off",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
			status:             "off",
		},
		{
			name:               "user set himself off",
			credAccessKey:      offUserAccessKey,
			credSecretKey:      offUserSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusOK,
			status:             "off",
		},
		{
			// after user set status off only root user can set status on
			name:               "user set himself on",
			credAccessKey:      offUserAccessKey,
			credSecretKey:      offUserSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusForbidden,
			status:             "on",
		},
		{
			name:               "user set a non exist user on",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
			status:             "on",
		},
		{
			name:               "user set other user off",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          otherUserOffAccessKey,
			secretKey:          otherUserOffSecretKey,
			expectedRespStatus: http.StatusForbidden,
			status:             "off",
		},
		{
			name:               "user set other user on",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          otherUserOffAccessKey,
			secretKey:          otherUserOffSecretKey,
			expectedRespStatus: http.StatusForbidden,
			status:             "on",
		},
	}
	setStatusUrl := "http://127.0.0.1:9985/admin/v1/update-accessKey_status?"
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			urlValues := make(url.Values)
			urlValues.Set(accessKey, testCase.accessKey)
			urlValues.Set(accountStatus, testCase.status)
			reqSetStatus := utils.MustNewSignedV4Request(http.MethodPost, setStatusUrl+urlValues.Encode(), 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result := reqTest(reqSetStatus)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})

	}
}
