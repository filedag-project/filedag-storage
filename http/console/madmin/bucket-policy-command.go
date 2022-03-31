package madmin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// PutBucketPolicy .
func (adm *AdminClient) PutBucketPolicy(ctx context.Context, bucketName, policyStr string) error {
	byte, err := json.Marshal(policyStr)
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName + "?policy",
		content: byte,
	}
	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}

// GetBucketPolicy .
func (adm *AdminClient) GetBucketPolicy(ctx context.Context, bucketName string) error {
	//q := make(url.Values)
	//q.Set("policy", "")
	//reqData := requestData{
	//	relPath: adminAPIPrefix + bucketName,
	//	queryValues: q,
	//}
	//resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	//defer closeResponse(resp)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(resp)
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	//var response policy.Policy
	//err = xml.Unmarshal(body, &response)
	//fmt.Println(response)
	//if resp.StatusCode != http.StatusOK {
	//	return httpRespToErrorResponse(resp)
	//}
	//return nil

	q := make(url.Values)
	q.Set("policy", "")
	reqData := requestData{
		relPath:     adminAPIPrefix + bucketName,
		queryValues: q,
	}
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}

// RemoveBucketPolicy .
func (adm *AdminClient) RemoveBucketPolicy(ctx context.Context, bucketName string) error {
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName + "?policy",
	}
	resp, err := adm.executeMethod(ctx, http.MethodDelete, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}
