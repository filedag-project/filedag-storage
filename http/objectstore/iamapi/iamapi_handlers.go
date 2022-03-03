package iamapi

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api/s3resp"
	"net/http"
)

//GetUserList get all user
func (iama *IamApiServer) GetUserList(w http.ResponseWriter, r *http.Request) {
	var resp ListUsersResponse
	resp.ListUsersResult.Users = iam.GlobalIAMSys.GetUserList(context.Background())
	s3resp.WriteXMLResponse(w, r, http.StatusOK, resp)
}
