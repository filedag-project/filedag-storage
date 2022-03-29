package madmin

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// PutObject .
func (adm *AdminClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64) error {

	queryValues := url.Values{}
	contentLength := strconv.FormatInt(objectSize, 10)
	queryValues.Set("contentLength", contentLength)
	payloadBytes, err := io.ReadAll(reader)
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName + adminAPIPrefix + objectName,
		content: payloadBytes,
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

// GetObject .
func (adm *AdminClient) GetObject(ctx context.Context, bucketName, objectName string) error {
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName + adminAPIPrefix + objectName,
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

// RemoveObject .
func (adm *AdminClient) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName + adminAPIPrefix + objectName,
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

// CopyObject .
func (adm *AdminClient) CopyObject(ctx context.Context, bucketName, objectName string) error {
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName + adminAPIPrefix + objectName,
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

// HeadObject .
func (adm *AdminClient) HeadObject(ctx context.Context, bucketName, objectName string) error {
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName + adminAPIPrefix + objectName,
	}
	resp, err := adm.executeMethod(ctx, http.MethodHead, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}
	return nil
}

// ListObject .
func (adm *AdminClient) ListObject(ctx context.Context, bucketName string) error {
	reqData := requestData{
		relPath: adminAPIPrefix + bucketName,
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
