package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	jwtgo "github.com/golang-jwt/jwt/v4"
	logging "github.com/ipfs/go-log/v2"
	"strconv"
	"strings"
	"time"
)

var log = logging.Logger("auth")

// Args - arguments to policy to check whether it is allowed
type Args struct {
	AccountName string              `json:"account"`
	Groups      []string            `json:"groups"`
	Action      s3action.Action     `json:"action"`
	Conditions  map[string][]string `json:"conditions"`
	BucketName  string              `json:"bucket"`
	IsOwner     bool                `json:"owner"`
	ObjectName  string              `json:"object"`
}

// Common errors generated for access and secret key validation.
var (
	errInvalidAccessKeyLength = fmt.Errorf("access key length should be between %d and %d", accessKeyMinLen, accessKeyMaxLen)
	errInvalidSecretKeyLength = fmt.Errorf("secret key length should be between %d and %d", secretKeyMinLen, secretKeyMaxLen)
)
var timeSentinel = time.Unix(0, 0).UTC()

// errInvalidDuration invalid token expiry
var errInvalidDuration = errors.New("invalid token expiry")

// Default access and secret keys.
const (
	DefaultAccessKey = "filedagadmin"
	DefaultSecretKey = "filedagadmin"
)
const (
	// Minimum length for  access key.
	accessKeyMinLen = 3

	// Maximum length for  access key.
	// There is no max length enforcement for access keys
	accessKeyMaxLen = 20

	// Minimum length for  secret key for both server and gateway mode.
	secretKeyMinLen = 8

	// Maximum secret key length , this
	// is used when autogenerating new credentials.
	// There is no max length enforcement for secret keys
	secretKeyMaxLen = 40

	// Alphanumeric table used for generating access keys.
	alphaNumericTable = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Total length of the alpha numeric table.
	alphaNumericTableLen = byte(len(alphaNumericTable))
)

//GetDefaultActiveCred Get DefaultActiveCred
func GetDefaultActiveCred() Credentials {
	return Credentials{
		AccessKey: DefaultAccessKey,
		SecretKey: DefaultSecretKey,
	}
}

const (
	// AccountOn indicates that credentials are enabled
	AccountOn = "on"
	// AccountOff indicates that credentials are disabled
	AccountOff = "off"
)

// Credentials holds access and secret keys.
type Credentials struct {
	AccessKey    string    `xml:"AccessKeyId" json:"accessKey,omitempty"`
	SecretKey    string    `xml:"SecretAccessKey" json:"secretKey,omitempty"`
	CreateTime   time.Time `xml:"CreateTime" json:"createTime,omitempty"`
	Expiration   time.Time `xml:"Expiration" json:"expiration,omitempty"`
	SessionToken string    `xml:"SessionToken" json:"sessionToken"`
	Status       string    `xml:"-" json:"status,omitempty"`
	ParentUser   string    `xml:"-" json:"parentUser,omitempty"`
}

// generateCredentials - creates randomly generated credentials of maximum
// allowed length.
func generateCredentials() (accessKey, secretKey string, err error) {
	// Generate access key.
	keyBytes, err := credentialsReadBytes(accessKeyMaxLen)
	if err != nil {
		return "", "", err
	}
	for i := 0; i < accessKeyMaxLen; i++ {
		keyBytes[i] = alphaNumericTable[keyBytes[i]%alphaNumericTableLen]
	}
	accessKey = string(keyBytes)

	// Generate secret key.
	keyBytes, err = credentialsReadBytes(secretKeyMaxLen)
	if err != nil {
		return "", "", err
	}

	secretKey = strings.ReplaceAll(string([]byte(base64.StdEncoding.EncodeToString(keyBytes))[:secretKeyMaxLen]),
		"/", "+")

	return accessKey, secretKey, nil
}

// CreateCredentials Error is returned if given access key or secret key are invalid length.
func CreateCredentials(accessKey, secretKey string) (cred Credentials, err error) {
	if !IsAccessKeyValid(accessKey) {
		return cred, errInvalidAccessKeyLength
	}
	if !IsSecretKeyValid(secretKey) {
		return cred, errInvalidSecretKeyLength
	}
	cred.AccessKey = accessKey
	cred.SecretKey = secretKey
	cred.Expiration = timeSentinel
	cred.Status = AccountOn
	return cred, nil
}

// CreateNewCredentialsWithMetadata - creates new credentials using the specified access & secret keys
// and generate a session token if a secret token is provided.
func CreateNewCredentialsWithMetadata(accessKey, secretKey string, m map[string]interface{}, tokenSecret string) (cred Credentials, err error) {
	if len(accessKey) < accessKeyMinLen || len(accessKey) > accessKeyMaxLen {
		return Credentials{}, errInvalidAccessKeyLength
	}

	if len(secretKey) < secretKeyMinLen || len(secretKey) > secretKeyMaxLen {
		return Credentials{}, errInvalidSecretKeyLength
	}

	cred.AccessKey = accessKey
	cred.SecretKey = secretKey
	cred.Status = AccountOn
	cred.CreateTime = time.Now().UTC()
	if tokenSecret == "" {
		cred.Expiration = timeSentinel
		return cred, nil
	}

	expiry, err := expToInt64(m["exp"])
	if err != nil {
		return cred, err
	}
	cred.Expiration = time.Unix(expiry, 0).UTC()

	cred.SessionToken, err = jWTSignWithAccessKey(cred.AccessKey, m, tokenSecret)
	if err != nil {
		return cred, err
	}

	return cred, nil
}

// GetNewCredentialsWithMetadata generates and returns new credential with expiry.
func GetNewCredentialsWithMetadata(m map[string]interface{}, tokenSecret string) (Credentials, error) {
	accessKey, secretKey, err := generateCredentials()
	if err != nil {
		return Credentials{}, err
	}
	return CreateNewCredentialsWithMetadata(accessKey, secretKey, m, tokenSecret)
}

// jWTSignWithAccessKey - generates a session token.
func jWTSignWithAccessKey(accessKey string, m map[string]interface{}, tokenSecret string) (string, error) {
	m["accessKey"] = accessKey
	jwt := jwtgo.NewWithClaims(jwtgo.SigningMethodHS512, jwtgo.MapClaims(m))
	return jwt.SignedString([]byte(tokenSecret))
}

// expToInt64 - convert input interface value to int64.
func expToInt64(expI interface{}) (expAt int64, err error) {
	switch exp := expI.(type) {
	case string:
		expAt, err = strconv.ParseInt(exp, 10, 64)
	case float64:
		expAt, err = int64(exp), nil
	case int64:
		expAt, err = exp, nil
	case int:
		expAt, err = int64(exp), nil
	case uint64:
		expAt, err = int64(exp), nil
	case uint:
		expAt, err = int64(exp), nil
	case json.Number:
		expAt, err = exp.Int64()
	case time.Duration:
		expAt, err = time.Now().UTC().Add(exp).Unix(), nil
	case nil:
		expAt, err = 0, nil
	default:
		expAt, err = 0, errInvalidDuration
	}
	if expAt < 0 {
		return 0, errInvalidDuration
	}
	return expAt, err
}
func credentialsReadBytes(size int) (data []byte, err error) {
	data = make([]byte, size)
	var n int
	if n, err = rand.Read(data); err != nil {
		return nil, err
	} else if n != size {
		return nil, fmt.Errorf("not enough data. Expected to read: %v bytes, got: %v bytes", size, n)
	}
	return data, nil
}

// IsAccessKeyValid - validate access key for right length.
func IsAccessKeyValid(accessKey string) bool {
	return len(accessKey) >= accessKeyMinLen
}

// IsSecretKeyValid - validate secret key for right length.
func IsSecretKeyValid(secretKey string) bool {
	return len(secretKey) >= secretKeyMinLen
}

// IsValid - returns whether credential is valid or not.
func (cred *Credentials) IsValid() bool {
	// Verify credentials if its enabled or not set.
	if cred.Status == AccountOff {
		return false
	}
	return IsAccessKeyValid(cred.AccessKey) && IsSecretKeyValid(cred.SecretKey) && !cred.IsExpired()
}

// IsExpired - returns whether Credential is expired or not.
func (cred *Credentials) IsExpired() bool {
	if cred.Expiration.IsZero() || cred.Expiration.Equal(timeSentinel) {
		return false
	}

	return cred.Expiration.Before(time.Now().UTC())
}

// IsTemp - returns whether credential is temporary or not.
func (cred *Credentials) IsTemp() bool {
	return cred.SessionToken != "" && !cred.Expiration.IsZero() && !cred.Expiration.Equal(timeSentinel)
}
