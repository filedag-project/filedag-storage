package controllers

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"io/ioutil"
	"net/http"
)

// ListObjects object list
func (control *Control) ListObjects(w http.ResponseWriter, r *http.Request) {
	var listObjectsParams *models.ListObjectsParams
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	json.Unmarshal(body, &listObjectsParams)
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	resp, error := control.apiServer.GetListObjectsResponse(principal, *listObjectsParams)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}
