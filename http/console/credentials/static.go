package credentials

// A Static is a set of credentials which are set programmatically,
// and will never expire.
type Static struct {
	Value
}

// NewStaticV2 returns a pointer to a new Credentials object
// wrapping a static credentials value provider, signature is
// set to v2. If access and secret are not specified then
// regardless of signature type set it Value will return
// as anonymous.
func NewStaticV2(id, secret, token string) *Credentials {
	return NewStatic(id, secret, token, SignatureV2)
}

// NewStaticV4 is similar to NewStaticV2 with similar considerations.
func NewStaticV4(id, secret, token string) *Credentials {
	return NewStatic(id, secret, token, SignatureV4)
}

// NewStatic returns a pointer to a new Credentials object
// wrapping a static credentials value provider.
func NewStatic(id, secret, token string, signerType SignatureType) *Credentials {
	return New(&Static{
		Value: Value{
			AccessKeyID:     id,
			SecretAccessKey: secret,
			SessionToken:    token,
			SignerType:      signerType,
		},
	})
}

// Retrieve returns the static credentials.
func (s *Static) Retrieve() (Value, error) {
	if s.AccessKeyID == "" || s.SecretAccessKey == "" {
		// Anonymous is not an error
		return Value{SignerType: SignatureAnonymous}, nil
	}
	return s.Value, nil
}

// IsExpired returns if the credentials are expired.
//
// For Static, the credentials never expired.
func (s *Static) IsExpired() bool {
	return false
}
