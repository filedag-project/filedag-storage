package madmin

import (
	"context"
	"io"
	"net/http"
)

// CreatePolicy .
func (adm *AdminClient) CreatePolicy(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64) error {
	//queryValues := url.Values{}
	reqData := requestData{
		relPath: adminAPIPrefix + "creat-policy",
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

// RemovePolicy .
func (adm *AdminClient) RemovePolicy(ctx context.Context, policy, objectName string, reader io.Reader, objectSize int64) error {
	//queryValues := url.Values{}
	reqData := requestData{
		relPath: adminAPIPrefix + "remove-policy",
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
