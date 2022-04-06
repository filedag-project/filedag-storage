package credentials

import "strings"

// SignatureType is type of Authorization requested for a given HTTP request.
type SignatureType int

// Different types of supported signatures - default is SignatureV4 or SignatureDefault.
const (
	// SignatureDefault is always set to v4.
	SignatureDefault SignatureType = iota
	SignatureV4
	SignatureV2
	SignatureV4Streaming
	SignatureAnonymous // Anonymous signature signifies, no signature.
)

// IsV2 - is signature SignatureV2?
func (s SignatureType) IsV2() bool {
	return s == SignatureV2
}

// IsV4 - is signature SignatureV4?
func (s SignatureType) IsV4() bool {
	return s == SignatureV4 || s == SignatureDefault
}

// IsStreamingV4 - is signature SignatureV4Streaming?
func (s SignatureType) IsStreamingV4() bool {
	return s == SignatureV4Streaming
}

// IsAnonymous - is signature empty?
func (s SignatureType) IsAnonymous() bool {
	return s == SignatureAnonymous
}

// Stringer humanized version of signature type,
// strings returned here are case insensitive.
func (s SignatureType) String() string {
	if s.IsV2() {
		return "S3v2"
	} else if s.IsV4() {
		return "S3v4"
	} else if s.IsStreamingV4() {
		return "S3v4Streaming"
	}
	return "Anonymous"
}

func parseSignatureType(str string) SignatureType {
	if strings.EqualFold(str, "S3v4") {
		return SignatureV4
	} else if strings.EqualFold(str, "S3v2") {
		return SignatureV2
	} else if strings.EqualFold(str, "S3v4Streaming") {
		return SignatureV4Streaming
	}
	return SignatureAnonymous
}
