package restapi

// list of all console environment constants
const (
	// Constants for common configuration
	ConsoleMinIOServer   = "CONSOLE_MINIO_SERVER"
	ConsoleSubnetProxy   = "CONSOLE_SUBNET_PROXY"
	ConsoleMinIORegion   = "CONSOLE_MINIO_REGION"
	ConsoleHostname      = "CONSOLE_HOSTNAME"
	ConsolePort          = "CONSOLE_PORT"
	ConsoleTLSHostname   = "CONSOLE_TLS_HOSTNAME"
	ConsoleTLSPort       = "CONSOLE_TLS_PORT"
	ConsoleSubnetLicense = "CONSOLE_SUBNET_LICENSE"
	MinIOConfigEnvFile   = "MINIO_CONFIG_ENV_FILE"
	MinIOSubnetLicense   = "MINIO_SUBNET_LICENSE"

	// Constants for Secure middleware
	ConsoleSecureAllowedHosts                    = "CONSOLE_SECURE_ALLOWED_HOSTS"
	ConsoleSecureAllowedHostsAreRegex            = "CONSOLE_SECURE_ALLOWED_HOSTS_ARE_REGEX"
	ConsoleSecureFrameDeny                       = "CONSOLE_SECURE_FRAME_DENY"
	ConsoleSecureContentTypeNoSniff              = "CONSOLE_SECURE_CONTENT_TYPE_NO_SNIFF"
	ConsoleSecureBrowserXSSFilter                = "CONSOLE_SECURE_BROWSER_XSS_FILTER"
	ConsoleSecureContentSecurityPolicy           = "CONSOLE_SECURE_CONTENT_SECURITY_POLICY"
	ConsoleSecureContentSecurityPolicyReportOnly = "CONSOLE_SECURE_CONTENT_SECURITY_POLICY_REPORT_ONLY"
	ConsoleSecureHostsProxyHeaders               = "CONSOLE_SECURE_HOSTS_PROXY_HEADERS"
	ConsoleSecureSTSSeconds                      = "CONSOLE_SECURE_STS_SECONDS"
	ConsoleSecureSTSIncludeSubdomains            = "CONSOLE_SECURE_STS_INCLUDE_SUB_DOMAINS"
	ConsoleSecureSTSPreload                      = "CONSOLE_SECURE_STS_PRELOAD"
	ConsoleSecureTLSRedirect                     = "CONSOLE_SECURE_TLS_REDIRECT"
	ConsoleSecureTLSHost                         = "CONSOLE_SECURE_TLS_HOST"
	ConsoleSecureTLSTemporaryRedirect            = "CONSOLE_SECURE_TLS_TEMPORARY_REDIRECT"
	ConsoleSecureForceSTSHeader                  = "CONSOLE_SECURE_FORCE_STS_HEADER"
	ConsoleSecurePublicKey                       = "CONSOLE_SECURE_PUBLIC_KEY"
	ConsoleSecureReferrerPolicy                  = "CONSOLE_SECURE_REFERRER_POLICY"
	ConsoleSecureFeaturePolicy                   = "CONSOLE_SECURE_FEATURE_POLICY"
	ConsoleSecureExpectCTHeader                  = "CONSOLE_SECURE_EXPECT_CT_HEADER"
	PrometheusURL                                = "CONSOLE_PROMETHEUS_URL"
	PrometheusJobID                              = "CONSOLE_PROMETHEUS_JOB_ID"
	ConsoleLogQueryURL                           = "CONSOLE_LOG_QUERY_URL"
	ConsoleLogQueryAuthToken                     = "CONSOLE_LOG_QUERY_AUTH_TOKEN"
	LogSearchQueryAuthToken                      = "LOGSEARCH_QUERY_AUTH_TOKEN"
	SlashSeparator                               = "/"
)
