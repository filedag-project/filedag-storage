package restapi

//
//import (
//	"bytes"
//	"crypto/tls"
//	"fmt"
//	models "github.com/filedag-project/filedag-storage/http/console/model"
//	"github.com/filedag-project/filedag-storage/http/console/restapi/operations"
//	"io"
//	"io/fs"
//	"log"
//	"net"
//	"net/http"
//	"path/filepath"
//	"strings"
//	"sync"
//	"time"
//
//	"github.com/klauspost/compress/gzhttp"
//
//	portal_ui "github.com/minio/console/portal-ui"
//	"github.com/minio/pkg/env"
//	"github.com/minio/pkg/mimedb"
//
//	"github.com/go-openapi/errors"
//	"github.com/go-openapi/swag"
//	"github.com/minio/console/pkg/auth"
//	"github.com/unrolled/secure"
//)
//
////go:generate swagger generate server --target ../../console --name Console --spec ../swagger.yml
//
//var additionalServerFlags = struct {
//	CertsDir string `long:"certs-dir" description:"path to certs directory" env:"CONSOLE_CERTS_DIR"`
//}{}
//
//var subPath = "/"
//var subPathOnce sync.Once
//
//func configureFlags(api *operations.ConsoleAPI) {
//	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
//		{
//			ShortDescription: "additional server flags",
//			Options:          &additionalServerFlags,
//		},
//	}
//}
//
//func configureAPI(api *operations.ConsoleAPI) http.Handler {
//	// Applies when the "x-token" header is set
//	api.KeyAuth = func(token string, scopes []string) (*models.Principal, error) {
//		// we are validating the session token by decrypting the claims inside, if the operation succeed that means the jwt
//		// was generated and signed by us in the first place
//		claims, err := auth.SessionTokenAuthenticate(token)
//		if err != nil {
//			api.Logger("Unable to validate the session token %s: %v", token, err)
//			return nil, errors.New(401, "incorrect api key auth")
//		}
//		return &models.Principal{
//			STSAccessKeyID:     claims.STSAccessKeyID,
//			STSSecretAccessKey: claims.STSSecretAccessKey,
//			STSSessionToken:    claims.STSSessionToken,
//			AccountAccessKey:   claims.AccountAccessKey,
//			Hm:                 claims.HideMenu,
//		}, nil
//	}
//
//	// Register login handlers
//	registerLoginHandlers(api)
//	// Register logout handlers
//	registerLogoutHandlers(api)
//	// Register bucket handlers
//	registerBucketsHandlers(api)
//	// Register all users handlers
//	registerUsersHandlers(api)
//	// Register groups handlers
//	registerGroupsHandlers(api)
//	// Register policies handlers
//	registersPoliciesHandler(api)
//	// Register configurations handlers
//	registerConfigHandlers(api)
//	// Register bucket events handlers
//	registerBucketEventsHandlers(api)
//	// Register bucket lifecycle handlers
//	registerBucketsLifecycleHandlers(api)
//	// Register service handlers
//	registerServiceHandlers(api)
//	// Register profiling handlers
//	registerProfilingHandler(api)
//	// Register session handlers
//	registerSessionHandlers(api)
//	// Register admin info handlers
//	registerAdminInfoHandlers(api)
//	// Register admin arns handlers
//	registerAdminArnsHandlers(api)
//	// Register admin notification endpoints handlers
//	registerAdminNotificationEndpointsHandlers(api)
//	// Register admin Service Account Handlers
//	registerServiceAccountsHandlers(api)
//	// Register admin remote buckets
//	registerAdminBucketRemoteHandlers(api)
//	// Register admin log search
//	registerLogSearchHandlers(api)
//	// Register admin subnet handlers
//	registerSubnetHandlers(api)
//	// Register Account handlers
//	registerAdminTiersHandlers(api)
//	//Register Inspect Handler
//	registerInspectHandler(api)
//
//	// Operator Console
//
//	// Register Object's Handlers
//	registerObjectsHandlers(api)
//	// Register Bucket Quota's Handlers
//	registerBucketQuotaHandlers(api)
//	// Register Account handlers
//	registerAccountHandlers(api)
//
//	api.PreServerShutdown = func() {}
//
//	api.ServerShutdown = func() {}
//
//	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
//}
//
//// The TLS configuration before HTTPS server starts.
//func configureTLS(tlsConfig *tls.Config) {
//	tlsConfig.RootCAs = GlobalRootCAs
//	tlsConfig.GetCertificate = GlobalTLSCertsManager.GetCertificate
//}
//
//// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
//// The middleware executes after routing but before authentication, binding and validation
//func setupMiddlewares(handler http.Handler) http.Handler {
//	return handler
//}
//
//// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
//// So this is a good place to plug in a panic handling middleware, logging and metrics
//func setupGlobalMiddleware(handler http.Handler) http.Handler {
//	// handle cookie or authorization header for session
//	next := AuthenticationMiddleware(handler)
//	// serve static files
//	next = FileServerMiddleware(next)
//
//	sslHostFn := secure.SSLHostFunc(func(host string) string {
//		h, _, err := net.SplitHostPort(host)
//		if err != nil {
//			return host
//		}
//		return net.JoinHostPort(h, TLSPort)
//	})
//
//	// Secure middleware, this middleware wrap all the previous handlers and add
//	// HTTP security headers
//	secureOptions := secure.Options{
//		AllowedHosts:                    GetSecureAllowedHosts(),
//		AllowedHostsAreRegex:            GetSecureAllowedHostsAreRegex(),
//		HostsProxyHeaders:               GetSecureHostsProxyHeaders(),
//		SSLRedirect:                     GetTLSRedirect() == "on" && len(GlobalPublicCerts) > 0,
//		SSLHostFunc:                     &sslHostFn,
//		SSLHost:                         GetSecureTLSHost(),
//		STSSeconds:                      GetSecureSTSSeconds(),
//		STSIncludeSubdomains:            GetSecureSTSIncludeSubdomains(),
//		STSPreload:                      GetSecureSTSPreload(),
//		SSLTemporaryRedirect:            false,
//		ForceSTSHeader:                  GetSecureForceSTSHeader(),
//		FrameDeny:                       GetSecureFrameDeny(),
//		ContentTypeNosniff:              GetSecureContentTypeNonSniff(),
//		BrowserXssFilter:                GetSecureBrowserXSSFilter(),
//		ContentSecurityPolicy:           GetSecureContentSecurityPolicy(),
//		ContentSecurityPolicyReportOnly: GetSecureContentSecurityPolicyReportOnly(),
//		PublicKey:                       GetSecurePublicKey(),
//		ReferrerPolicy:                  GetSecureReferrerPolicy(),
//		FeaturePolicy:                   GetSecureFeaturePolicy(),
//		ExpectCTHeader:                  GetSecureExpectCTHeader(),
//		IsDevelopment:                   false,
//	}
//	secureMiddleware := secure.New(secureOptions)
//	next = secureMiddleware.Handler(next)
//	gnext := gzhttp.GzipHandler(next)
//	return RejectS3Middleware(gnext)
//}
//
//// RejectS3Middleware will reject requests that have AWS S3 specific headers.
//func RejectS3Middleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if len(r.Header.Get("X-Amz-Content-Sha256")) > 0 ||
//			len(r.Header.Get("X-Amz-Date")) > 0 ||
//			strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256") ||
//			r.URL.Query().Get("AWSAccessKeyId") != "" {
//			w.WriteHeader(http.StatusForbidden)
//			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
//<Error>
//<Code>AccessDenied</Code>
//<Message>S3 API Request made to Console port. S3 Requests should be sent to API port.</Message>
//<RequestId>0</RequestId>
//</Error>
//`))
//			return
//		}
//		next.ServeHTTP(w, r)
//	})
//}
//
//func AuthenticationMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token, err := auth.GetTokenFromRequest(r)
//		if err != nil && err != auth.ErrNoAuthToken {
//			http.Error(w, err.Error(), http.StatusUnauthorized)
//			return
//		}
//		// All handlers handle appropriately to return errors
//		// based on their swagger rules, we do not need to
//		// additionally return error here, let the next ServeHTTPs
//		// handle it appropriately.
//		if token != "" {
//			r.Header.Add("Authorization", "Bearer "+token)
//		}
//		next.ServeHTTP(w, r)
//	})
//}
//
//// FileServerMiddleware serves files from the static folder
//func FileServerMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Server", globalAppName) // do not add version information
//		switch {
//		case strings.HasPrefix(r.URL.Path, "/ws"):
//			serveWS(w, r)
//		case strings.HasPrefix(r.URL.Path, "/api"):
//			next.ServeHTTP(w, r)
//		default:
//			buildFs, err := fs.Sub(portal_ui.GetStaticAssets(), "build")
//			if err != nil {
//				panic(err)
//			}
//			wrapHandlerSinglePageApplication(http.FileServer(http.FS(buildFs))).ServeHTTP(w, r)
//		}
//	})
//}
//
//type notFoundRedirectRespWr struct {
//	http.ResponseWriter // We embed http.ResponseWriter
//	status              int
//}
//
//func (w *notFoundRedirectRespWr) WriteHeader(status int) {
//	w.status = status // Store the status for our own use
//	if status != http.StatusNotFound {
//		w.ResponseWriter.WriteHeader(status)
//	}
//}
//
//func (w *notFoundRedirectRespWr) Write(p []byte) (int, error) {
//	if w.status != http.StatusNotFound {
//		return w.ResponseWriter.Write(p)
//	}
//	return len(p), nil // Lie that we successfully wrote it
//}
//
//func handleSPA(w http.ResponseWriter, r *http.Request) {
//	basePath := "/"
//	// For SPA mode we will replace root base with a sub path if configured unless we received cp=y and cpb=/NEW/BASE
//	if v := r.URL.Query().Get("cp"); v == "y" {
//		if base := r.URL.Query().Get("cpb"); base != "" {
//			// make sure the subpath has a trailing slash
//			if !strings.HasSuffix(base, "/") {
//				base = fmt.Sprintf("%s/", base)
//			}
//			basePath = base
//		}
//	}
//
//	indexPage, err := portal_ui.GetStaticAssets().Open("build/index.html")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	indexPageBytes, err := io.ReadAll(indexPage)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// if we have a seeded basePath. This should override CONSOLE_SUBPATH every time, thus the `if else`
//	if basePath != "/" {
//		indexPageBytes = replaceBaseInIndex(indexPageBytes, basePath)
//		// if we have a custom subpath replace it in
//	} else if getSubPath() != "/" {
//		indexPageBytes = replaceBaseInIndex(indexPageBytes, getSubPath())
//	}
//
//	mimeType := mimedb.TypeByExtension(filepath.Ext(r.URL.Path))
//
//	if mimeType == "application/octet-stream" {
//		mimeType = "text/html"
//	}
//
//	w.Header().Set("Content-Type", mimeType)
//	http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader(indexPageBytes))
//}
//
//// wrapHandlerSinglePageApplication handles a http.FileServer returning a 404 and overrides it with index.html
//func wrapHandlerSinglePageApplication(h http.Handler) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if r.URL.Path == "/" {
//			handleSPA(w, r)
//			return
//		}
//
//		w.Header().Set("Content-Type", mimedb.TypeByExtension(filepath.Ext(r.URL.Path)))
//		nfw := &notFoundRedirectRespWr{ResponseWriter: w}
//		h.ServeHTTP(nfw, r)
//		if nfw.status == http.StatusNotFound {
//			handleSPA(w, r)
//		}
//	}
//}
//
//type nullWriter struct{}
//
//func (lw nullWriter) Write(b []byte) (int, error) {
//	return len(b), nil
//}
//
//// As soon as server is initialized but not run yet, this function will be called.
//// If you need to modify a config, store server instance to stop it individually later, this is the place.
//// This function can be called multiple times, depending on the number of serving schemes.
//// scheme value will be set accordingly: "http", "https" or "unix"
//func configureServer(s *http.Server, _, _ string) {
//	// Turn-off random logging by Go net/http
//	s.ErrorLog = log.New(&nullWriter{}, "", 0)
//}
//
//func getSubPath() string {
//	subPathOnce.Do(func() {
//		if v := env.Get("CONSOLE_SUBPATH", ""); v != "" {
//			// Replace all unnecessary `\` to `/`
//			// also add pro-actively at the end.
//			subPath = filepath.Clean(filepath.ToSlash(v)) + SlashSeparator
//			if !strings.HasPrefix(subPath, SlashSeparator) {
//				subPath = SlashSeparator + subPath
//			}
//		}
//	})
//	return subPath
//}
//
//func replaceBaseInIndex(indexPageBytes []byte, basePath string) []byte {
//	indexPageStr := string(indexPageBytes)
//	newBase := fmt.Sprintf("<base href=\"%s\"/>", basePath)
//	indexPageStr = strings.Replace(indexPageStr, "<base href=\"/\"/>", newBase, 1)
//	indexPageBytes = []byte(indexPageStr)
//	return indexPageBytes
//}
