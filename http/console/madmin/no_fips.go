//go:build !fips
// +build !fips

package madmin

// FIPSEnabled returns true if and only if FIPS 140-2 support
// is enabled.
//
// FIPS 140-2 requires that only specifc cryptographic
// primitives, like AES or SHA-256, are used and that
// those primitives are implemented by a FIPS 140-2
// certified cryptographic module.
func FIPSEnabled() bool { return false }
