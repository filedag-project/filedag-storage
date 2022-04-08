package restapi

import (
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/pkg/auth"
)

type ApiServer struct {
}

// login
func login(credentials ConsoleCredentialsI, sessionFeatures *auth.SessionFeatures) (*string, error) {
	// try to obtain consoleCredentials,
	tokens, err := credentials.Get()
	if err != nil {
		return nil, err
	}
	// if we made it here, the consoleCredentials work, generate a jwt with claims
	token, err := auth.NewEncryptedTokenForClient(&tokens, credentials.GetAccountAccessKey(), sessionFeatures)
	if err != nil {
		LogError("error authenticating user: %v", err)
		return nil, errInvalidCredentials
	}
	return &token, nil
}

// getConsoleCredentials will return ConsoleCredentials interface
func getConsoleCredentials(accessKey, secretKey string) (*ConsoleCredentials, error) {
	creds, err := NewConsoleCredentials(accessKey, secretKey, GetRegion())
	if err != nil {
		return nil, err
	}
	return &ConsoleCredentials{
		ConsoleCredentials: creds,
		AccountAccessKey:   accessKey,
	}, nil
}

// GetLoginResponse performs login() and serializes it to the handler's output
func (apiServer *ApiServer) GetLoginResponse(lr *models.LoginRequest) (*models.LoginResponse, *models.Error) {
	// prepare console credentials
	consoleCreds, err := getConsoleCredentials(*lr.AccessKey, *lr.SecretKey)
	if err != nil {
		return nil, prepareError(err, errInvalidCredentials, err)
	}
	sf := &auth.SessionFeatures{}
	if lr.Features != nil {
		sf.HideMenu = lr.Features.HideMenu
	}
	sessionID, err := login(consoleCreds, sf)
	if err != nil {
		return nil, prepareError(err, errInvalidCredentials, err)
	}
	// serialize output
	loginResponse := &models.LoginResponse{
		SessionID: *sessionID,
	}
	return loginResponse, nil
}
