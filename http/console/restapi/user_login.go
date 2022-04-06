package restapi

import (
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/pkg/auth"
)

//func registerLoginHandlers(api *operations.ConsoleAPI) {
//	// GET login strategy
//	api.UserAPILoginDetailHandler = user_api.LoginDetailHandlerFunc(func(params user_api.LoginDetailParams) middleware.Responder {
//		loginDetails, err := getLoginDetailsResponse(params.HTTPRequest)
//		if err != nil {
//			return user_api.NewLoginDetailDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewLoginDetailOK().WithPayload(loginDetails)
//	})
//	// POST login using user credentials
//	api.UserAPILoginHandler = user_api.LoginHandlerFunc(func(params user_api.LoginParams) middleware.Responder {
//		loginResponse, err := getLoginResponse(params.Body)
//		if err != nil {
//			return user_api.NewLoginDefault(int(err.Code)).WithPayload(err)
//		}
//		// Custom response writer to set the session cookies
//		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
//			cookie := NewSessionCookieForConsole(loginResponse.SessionID)
//			http.SetCookie(w, &cookie)
//			user_api.NewLoginNoContent().WriteResponse(w, p)
//		})
//	})
//	// POST login using external IDP
//	api.UserAPILoginOauth2AuthHandler = user_api.LoginOauth2AuthHandlerFunc(func(params user_api.LoginOauth2AuthParams) middleware.Responder {
//		loginResponse, err := getLoginOauth2AuthResponse(params.HTTPRequest, params.Body)
//		if err != nil {
//			return user_api.NewLoginOauth2AuthDefault(int(err.Code)).WithPayload(err)
//		}
//		// Custom response writer to set the session cookies
//		return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
//			cookie := NewSessionCookieForConsole(loginResponse.SessionID)
//			http.SetCookie(w, &cookie)
//			user_api.NewLoginOauth2AuthNoContent().WriteResponse(w, p)
//		})
//	})
//}

// login performs a check of ConsoleCredentials against MinIO, generates some claims and returns the jwt
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
	creds, err := NewConsoleCredentials(accessKey, secretKey, GetMinIORegion())
	if err != nil {
		return nil, err
	}
	return &ConsoleCredentials{
		ConsoleCredentials: creds,
		AccountAccessKey:   accessKey,
	}, nil
}

// getLoginResponse performs login() and serializes it to the handler's output
func getLoginResponse(lr *models.LoginRequest) (*models.LoginResponse, *models.Error) {
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
