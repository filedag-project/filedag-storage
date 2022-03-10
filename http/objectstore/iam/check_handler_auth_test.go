package iam

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestV2CheckRequestAuthType(t *testing.T) {
	var aSys AuthSys
	aSys.Init()
	req := mustNewSignedV2Request("GET", "http://127.0.0.1:9000", 0, nil, t)
	err := aSys.CheckRequestAuthType(context.Background(), req, s3action.ListAllMyBucketsAction, "test", "testobject")
	fmt.Println(api_errors.GetAPIError(err))
}
func TestV4CheckRequestAuthType(t *testing.T) {
	var aSys AuthSys
	aSys.Init()
	req := mustNewSignedV4Request("GET", "http://127.0.0.1:9000", 0, nil, t)
	err := aSys.CheckRequestAuthType(context.Background(), req, s3action.ListAllMyBucketsAction, "test", "testobject")
	fmt.Println(api_errors.GetAPIError(err))
}

var ignoredHeaders = map[string]bool{
	"Authorization":  true,
	"Content-Type":   true,
	"Content-Length": true,
	"User-Agent":     true,
}

// Sign given request using Signature V4.
func signRequestV4(req *http.Request, accessKey, secretKey string) error {
	// Get hashed payload.
	hashedPayload := req.Header.Get("x-amz-content-sha256")
	if hashedPayload == "" {
		return fmt.Errorf("Invalid hashed payload")
	}

	currTime := time.Now()

	// Set x-amz-date.
	req.Header.Set("x-amz-date", currTime.Format(iso8601Format))

	// Get header map.
	headerMap := make(map[string][]string)
	for k, vv := range req.Header {
		// If request header key is not in ignored headers, then add it.
		if _, ok := ignoredHeaders[http.CanonicalHeaderKey(k)]; !ok {
			headerMap[strings.ToLower(k)] = vv
		}
	}

	// Get header keys.
	headers := []string{"host"}
	for k := range headerMap {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	region := "us-east-1"

	// Get canonical headers.
	var buf bytes.Buffer
	for _, k := range headers {
		buf.WriteString(k)
		buf.WriteByte(':')
		switch {
		case k == "host":
			buf.WriteString(req.URL.Host)
			fallthrough
		default:
			for idx, v := range headerMap[k] {
				if idx > 0 {
					buf.WriteByte(',')
				}
				buf.WriteString(v)
			}
			buf.WriteByte('\n')
		}
	}
	canonicalHeaders := buf.String()

	// Get signed headers.
	signedHeaders := strings.Join(headers, ";")

	// Get canonical query string.
	req.URL.RawQuery = strings.Replace(req.URL.Query().Encode(), "+", "%20", -1)

	// Get canonical URI.
	canonicalURI := EncodePath(req.URL.Path)

	// Get canonical request.
	// canonicalRequest =
	//  <HTTPMethod>\n
	//  <CanonicalURI>\n
	//  <CanonicalQueryString>\n
	//  <CanonicalHeaders>\n
	//  <SignedHeaders>\n
	//  <HashedPayload>
	//
	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		hashedPayload,
	}, "\n")

	// Get scope.
	scope := strings.Join([]string{
		currTime.Format(yyyymmdd),
		region,
		"s3",
		"aws4_request",
	}, "/")

	stringToSign := "AWS4-HMAC-SHA256" + "\n" + currTime.Format(iso8601Format) + "\n"
	stringToSign = stringToSign + scope + "\n"
	stringToSign = stringToSign + getSHA256Hash([]byte(canonicalRequest))

	date := sumHMAC([]byte("AWS4"+secretKey), []byte(currTime.Format(yyyymmdd)))
	regionHMAC := sumHMAC(date, []byte(region))
	service := sumHMAC(regionHMAC, []byte("s3"))
	signingKey := sumHMAC(service, []byte("aws4_request"))

	signature := hex.EncodeToString(sumHMAC(signingKey, []byte(stringToSign)))

	// final Authorization header
	parts := []string{
		"AWS4-HMAC-SHA256" + " Credential=" + accessKey + "/" + scope,
		"SignedHeaders=" + signedHeaders,
		"Signature=" + signature,
	}
	author := strings.Join(parts, ", ")
	req.Header.Set("Authorization", author)

	return nil
}

// Sign given request using Signature V4.
func signRequestV2(req *http.Request) error {
	creds := auth.GetDefaultActiveCred()

	policy := "policy"
	req.Header.Set("Awsaccesskeyid", creds.AccessKey)
	req.Header.Set("Policy", policy)
	stringToSign := getStringToSignV2(req.Method, "", "", req.Header, "")
	req.Header.Set("Signature", calculateSignatureV2(stringToSign, creds.SecretKey))
	req.Header.Set("Authorization", "AWS "+creds.AccessKey+":"+calculateSignatureV2(stringToSign, creds.SecretKey))

	return nil
}

// EncodePath encode the strings from UTF-8 byte representations to HTML hex escape sequences
//
// This is necessary since regular url.Parse() and url.Encode() functions do not support UTF-8
// non english characters cannot be parsed due to the nature in which url.Encode() is written
//
// This function on the other hand is a direct replacement for url.Encode() technique to support
// pretty much every UTF-8 character.
func EncodePath(pathName string) string {
	if reservedObjectNames.MatchString(pathName) {
		return pathName
	}
	var encodedPathname string
	for _, s := range pathName {
		if 'A' <= s && s <= 'Z' || 'a' <= s && s <= 'z' || '0' <= s && s <= '9' { // §2.3 Unreserved characters (mark)
			encodedPathname = encodedPathname + string(s)
			continue
		}
		switch s {
		case '-', '_', '.', '~', '/': // §2.3 Unreserved characters (mark)
			encodedPathname = encodedPathname + string(s)
			continue
		default:
			len := utf8.RuneLen(s)
			if len < 0 {
				// if utf8 cannot convert return the same string as is
				return pathName
			}
			u := make([]byte, len)
			utf8.EncodeRune(u, s)
			for _, r := range u {
				hex := hex.EncodeToString([]byte{r})
				encodedPathname = encodedPathname + "%" + strings.ToUpper(hex)
			}
		}
	}
	return encodedPathname
}
func mustNewSignedV2Request(method string, urlStr string, contentLength int64, body io.ReadSeeker, t *testing.T) *http.Request {
	req := mustNewRequest(method, urlStr, contentLength, body, t)
	if err := signRequestV2(req); err != nil {
		t.Fatalf("Unable to inititalized new signed http request %s", err)
	}
	return req
}
func mustNewSignedV4Request(method string, urlStr string, contentLength int64, body io.ReadSeeker, t *testing.T) *http.Request {
	req := mustNewRequest(method, urlStr, contentLength, body, t)
	cred := &auth.Credentials{AccessKey: "test", SecretKey: "test"}
	if err := signRequestV4(req, cred.AccessKey, cred.SecretKey); err != nil {
		t.Fatalf("Unable to inititalized new signed http request %s", err)
	}
	return req
}

// Provides a fully populated http request instance, fails otherwise.
func mustNewRequest(method string, urlStr string, contentLength int64, body io.ReadSeeker, t *testing.T) *http.Request {
	req, err := newTestRequest(method, urlStr, contentLength, body)
	if err != nil {
		t.Fatalf("Unable to initialize new http request %s", err)
	}
	return req
}

// Returns new HTTP request object.
func newTestRequest(method, urlStr string, contentLength int64, body io.ReadSeeker) (*http.Request, error) {
	if method == "" {
		method = "POST"
	}

	// Save for subsequent use
	var hashedPayload string
	var md5Base64 string
	switch {
	case body == nil:
		hashedPayload = getSHA256Hash([]byte{})
	default:
		payloadBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, err
		}
		hashedPayload = getSHA256Hash(payloadBytes)
		md5Base64 = getMD5HashBase64(payloadBytes)
	}
	// Seek back to beginning.
	if body != nil {
		body.Seek(0, 0)
	} else {
		body = bytes.NewReader([]byte(""))
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if md5Base64 != "" {
		req.Header.Set("Content-Md5", md5Base64)
	}
	req.Header.Set("x-amz-content-sha256", hashedPayload)

	// Add Content-Length
	req.ContentLength = contentLength

	return req, nil
}

// getSHA256Hash returns SHA-256 hash in hex encoding of given data.
func getSHA256Hash(data []byte) string {
	return hex.EncodeToString(getSHA256Sum(data))
}

// getMD5HashBase64 returns MD5 hash in base64 encoding of given data.
func getMD5HashBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(getMD5Sum(data))
}

// getSHA256Hash returns SHA-256 sum of given data.
func getSHA256Sum(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// getMD5Sum returns MD5 sum of given data.
func getMD5Sum(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}
