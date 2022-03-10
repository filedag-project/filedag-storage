package iam

import (
	"bytes"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestStsAPIHandlers_AssumeRole(t *testing.T) {
	body := bytes.NewReader([]byte("Version=2011-06-15&Action=AssumeRole"))
	req := mustNewSignedV4RequestSts(http.MethodPost, "http://127.0.0.1:9985/", 0, body, t)
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(err)
	fmt.Printf("resp%+v", resp.Body)
}
func mustNewSignedV4RequestSts(method string, urlStr string, contentLength int64, body io.ReadSeeker, t *testing.T) *http.Request {
	req := mustNewRequest(method, urlStr, contentLength, body, t)
	cred := &auth.Credentials{AccessKey: "test", SecretKey: "test"}
	if err := signRequestV4Sts(req, cred.AccessKey, cred.SecretKey); err != nil {
		t.Fatalf("Unable to inititalized new signed http request %s", err)
	}
	return req
}

// Sign given request using Signature V4.
func signRequestV4Sts(req *http.Request, accessKey, secretKey string) error {
	// Get hashed payload.
	hashedPayload := getContentSha256Cksum(req, ServiceSTS)
	fmt.Println(hashedPayload)
	currTime := time.Now()

	// Set x-amz-date.
	req.Header.Set("x-amz-date", currTime.Format(iso8601Format))
	req.Header.Set(consts.ContentType, "application/x-www-form-urlencoded")
	// Query string.
	// final Authorization header
	// Get header keys.
	// Get header map.
	headerMap := make(map[string][]string)
	for k, vv := range req.Header {
		// If request header key is not in ignored headers, then add it.
		if _, ok := ignoredHeaders[http.CanonicalHeaderKey(k)]; !ok {
			headerMap[strings.ToLower(k)] = vv
		}
	}
	headers := []string{"host"}
	for k := range headerMap {
		headers = append(headers, k)
	}
	sort.Strings(headers)

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
	headerMap["host"] = append(headerMap["host"], req.URL.Host)

	// Get signed headers.
	signedHeaders := strings.Join(headers, ";")
	queryStr := req.Form.Encode()
	region := "us-east-1"
	// Get scope.
	scope := strings.Join([]string{
		currTime.Format(yyyymmdd),
		region,
		string(ServiceSTS),
		"aws4_request",
	}, "/")
	// Get canonical request.
	fmt.Printf("headerMap:%+v,hashedPayload:%v,queryStr:%v,req.URL.Path:%v,req.Method:%v\n", headerMap, hashedPayload, queryStr, req.URL.Path, req.Method)
	canonicalRequest := getCanonicalRequest(headerMap, hashedPayload, queryStr, req.URL.Path, req.Method)
	// Get string to sign from canonical request.
	fmt.Printf("canonicalRequest:%v,currTime:%v,scope:%v\n", canonicalRequest, currTime, scope)
	stringToSign := getStringToSign(canonicalRequest, currTime, scope)
	fmt.Println("stringToSign", stringToSign)

	// Get hmac signing key.
	signingKey := getSigningKey(secretKey, currTime, region, ServiceSTS)
	fmt.Println("signingKey", signingKey)

	// Calculate signature.
	newSignature := getSignature(signingKey, stringToSign)
	fmt.Println(newSignature)

	parts := []string{
		"AWS4-HMAC-SHA256" + " Credential=" + accessKey + "/" + scope,
		"SignedHeaders=" + signedHeaders,
		"Signature=" + newSignature,
	}
	author := strings.Join(parts, ", ")
	req.Header.Set("Authorization", author)

	return nil
}
