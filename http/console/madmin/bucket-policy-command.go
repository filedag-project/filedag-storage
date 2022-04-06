package madmin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/madmin/policy"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// PutBucketPolicy .
func (adm *AdminClient) PutBucketPolicy(ctx context.Context, bucketName, policyStr string) error {
	payloadBytes, err := io.ReadAll(strings.NewReader(policyStr))
	q := make(url.Values)
	q.Set("policy", "")
	reqData := requestData{
		relPath:     adminAPIPrefix + bucketName,
		content:     payloadBytes,
		queryValues: q,
	}
	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)
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

// GetBucketPolicy .
func (adm *AdminClient) GetBucketPolicy(ctx context.Context, bucketName string) (*policy.Policy, error) {
	var response *policy.Policy
	q := make(url.Values)
	q.Set("policy", "")
	reqData := requestData{
		relPath:     adminAPIPrefix + bucketName,
		queryValues: q,
	}
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	defer closeResponse(resp)
	if err != nil {
		return response, err
	}
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	err = json.Unmarshal(body, &response)
	fmt.Println(response)
	if resp.StatusCode != http.StatusOK {
		return response, httpRespToErrorResponse(resp)
	}
	return response, nil
}

// RemoveBucketPolicy .
func (adm *AdminClient) RemoveBucketPolicy(ctx context.Context, bucketName string) error {
	q := make(url.Values)
	q.Set("policy", "")
	reqData := requestData{
		relPath:     adminAPIPrefix + bucketName,
		queryValues: q,
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
