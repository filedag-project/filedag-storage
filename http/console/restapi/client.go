package restapi

import (
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/credentials"
	"github.com/filedag-project/filedag-storage/http/console/models"
	xjwt "github.com/filedag-project/filedag-storage/http/console/pkg/token"
	"net/url"
	"path"
	"strings"
)

// ConsoleCredentialsI interface with all functions to be implemented
// by mock when testing, it should include all needed consoleCredentials.Login api calls
// that are used within this project.
type ConsoleCredentialsI interface {
	Get() (credentials.Value, error)
	Expire()
	GetAccountAccessKey() string
}

// Interface implementation
type ConsoleCredentials struct {
	ConsoleCredentials *credentials.Credentials
	AccountAccessKey   string
}

func (c ConsoleCredentials) GetAccountAccessKey() string {
	return c.AccountAccessKey
}

// Get implements *Login.Get()
func (c ConsoleCredentials) Get() (credentials.Value, error) {
	return c.ConsoleCredentials.Get()
}

// Expire implements *Login.Expire()
func (c ConsoleCredentials) Expire() {
	c.ConsoleCredentials.Expire()
}

// consoleSTSAssumeRole it's a STSAssumeRole wrapper, in general
// there's no need to use this struct anywhere else in the project, it's only required
// for passing a custom *http.Client to *credentials.STSAssumeRole
type consoleSTSAssumeRole struct {
	stsAssumeRole *credentials.STSAssumeRole
}

func (s consoleSTSAssumeRole) Retrieve() (credentials.Value, error) {
	return s.stsAssumeRole.Retrieve()
}

func (s consoleSTSAssumeRole) IsExpired() bool {
	return s.stsAssumeRole.IsExpired()
}

func NewConsoleCredentials(accessKey, secretKey, location string) (*credentials.Credentials, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("credentials endpoint, access and secret key are mandatory for AssumeRoleSTS")
	}
	opts := credentials.STSAssumeRoleOptions{
		AccessKey:       accessKey,
		SecretKey:       secretKey,
		Location:        location,
		DurationSeconds: int(xjwt.GetConsoleSTSDuration().Seconds()),
	}
	stsAssumeRole := &credentials.STSAssumeRole{
		Client:      GetConsoleHTTPClient(),
		STSEndpoint: getServer(),
		Options:     opts,
	}
	consoleSTSWrapper := consoleSTSAssumeRole{stsAssumeRole: stsAssumeRole}
	return credentials.New(consoleSTSWrapper), nil
}

// getConsoleCredentialsFromSession returns the *consoleCredentials.Login associated to the
// provided session token, this is useful for running the Expire() or IsExpired() operations
func getConsoleCredentialsFromSession(claims *models.Principal) *credentials.Credentials {
	return credentials.NewStaticV4(claims.STSAccessKeyID, claims.STSSecretAccessKey, claims.STSSessionToken)
}

// computeObjectURLWithoutEncode returns a MinIO url containing the object filename without encoding
func computeObjectURLWithoutEncode(bucketName, prefix string) (string, error) {
	endpoint := getServer()
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("the provided endpoint is invalid")
	}
	objectURL := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	if strings.TrimSpace(bucketName) != "" {
		objectURL = path.Join(objectURL, bucketName)
	}
	if strings.TrimSpace(prefix) != "" {
		objectURL = pathJoinFinalSlash(objectURL, prefix)
	}

	objectURL = fmt.Sprintf("%s://%s", u.Scheme, objectURL)
	return objectURL, nil
}

// pathJoinFinalSlash - like path.Join() but retains trailing slashSeparator of the last element
func pathJoinFinalSlash(elem ...string) string {
	if len(elem) > 0 {
		if strings.HasSuffix(elem[len(elem)-1], SlashSeparator) {
			return path.Join(elem...) + SlashSeparator
		}
	}
	return path.Join(elem...)
}
