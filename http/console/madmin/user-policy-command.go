package madmin

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// PutUserPolicy .
func (adm *AdminClient) PutUserPolicy(ctx context.Context, userName, policyStr string) error {
	queryValues := url.Values{}
	queryValues.Set("userName", userName)
	queryValues.Set("policyName", "read")
	queryValues.Set("policyDocument", policyStr)
	reqData := requestData{
		relPath:     adminAPIPrefix + "admin/v1/put-user-policy",
		queryValues: queryValues,
	}
	resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}

// GetUserPolicy .
func (adm *AdminClient) GetUserPolicy(ctx context.Context, userName string) error {
	queryValues := make(url.Values)
	queryValues.Set("userName", userName)
	queryValues.Set("policyName", "read")
	reqData := requestData{
		relPath:     adminAPIPrefix + "admin/v1/get-user-policy",
		queryValues: queryValues,
	}
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var response GetUserPolicyResponse
	err = xml.Unmarshal(body, &response)
	fmt.Println(response)
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}

// RemoveUserPolicy .
func (adm *AdminClient) RemoveUserPolicy(ctx context.Context, userName string) error {
	queryValues := make(url.Values)
	queryValues.Set("userName", userName)
	queryValues.Set("policyName", "read")
	reqData := requestData{
		relPath:     adminAPIPrefix + "admin/v1/remove-user-policy",
		queryValues: queryValues,
	}
	resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}
