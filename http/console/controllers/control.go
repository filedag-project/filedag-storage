package controllers

import (
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/pkg/auth"
	"github.com/filedag-project/filedag-storage/http/console/restapi"
)

type Control struct {
	apiServer *restapi.ApiServer
}

func VerifyToken(token string) (*models.Principal, error) {
	// we are validating the session token by decrypting the claims inside, if the operation succeed that means the jwt
	// was generated and signed by us in the first place
	claims, err := auth.SessionTokenAuthenticate(token)
	if err != nil {
		restapi.LogInfo("Unable to validate the session token %s: %v", token, err)
		return nil, err
	}
	return &models.Principal{
		STSAccessKeyID:     claims.STSAccessKeyID,
		STSSecretAccessKey: claims.STSSecretAccessKey,
		STSSessionToken:    claims.STSSessionToken,
		AccountAccessKey:   claims.AccountAccessKey,
		Hm:                 claims.HideMenu,
	}, nil
}
