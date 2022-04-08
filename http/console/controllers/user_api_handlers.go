package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"io/ioutil"
	"net/http"
)

// UserLogin user login
func (control *Control) UserLogin(w http.ResponseWriter, r *http.Request) {
	var resp *models.LoginResponse
	var loginRequest *models.LoginRequest
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Println(error)
	}
	json.Unmarshal(body, &loginRequest)
	resp, err := control.apiServer.GetLoginResponse(loginRequest)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}
