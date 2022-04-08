package controllers

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"io/ioutil"
	"net/http"
)

// ListBuckets bucket list
func (control *Control) ListBuckets(w http.ResponseWriter, r *http.Request) {
	var resp *models.ListBucketsResponse
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	resp, error := control.apiServer.GetListBucketsResponse(principal)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// AddBucket bucket add
func (control *Control) AddBucket(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		BucketName string `json:"bucket_name"`
		Location   string `json:"location"`
	}
	param := new(Params)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	json.Unmarshal(body, &param)
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	error := control.apiServer.GetCreateBucketResponse(principal, param.BucketName, param.Location, false)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, "")
}

// RemoveBucket bucket remove
func (control *Control) RemoveBucket(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		BucketName string `json:"bucket_name"`
		Location   string `json:"location"`
	}
	param := new(Params)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	json.Unmarshal(body, &param)
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	error := control.apiServer.GetDeleteBucketResponse(principal, param.BucketName)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, "")
}

// SetBucketPolicy set bucket policy
func (control *Control) SetBucketPolicy(w http.ResponseWriter, r *http.Request) {
	var setBucketPolicyParams *models.SetBucketPolicyParams
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	json.Unmarshal(body, &setBucketPolicyParams)
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	policy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::bucket1/*"}]}`
	setBucketPolicyParams.Definition = policy
	resp, error := control.apiServer.GetBucketSetPolicyResponse(principal, setBucketPolicyParams)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// GetBucketPolicy get bucket policy
func (control *Control) GetBucketPolicy(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket_name")
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	resp, error := control.apiServer.GetBucketPolicyResponse(principal, bucketName)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// RemoveBucketPolicy remove bucket policy
func (control *Control) RemoveBucketPolicy(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		BucketName string `json:"bucket_name"`
		Location   string `json:"location"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var params Params
	json.Unmarshal(body, &params)
	//name := r.FormValue("policy_name")
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	error := control.apiServer.RemoveBucketPolicyResponse(principal, params.BucketName)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, "")
}
