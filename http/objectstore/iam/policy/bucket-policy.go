package policy

// IPolicySys - policy subsystem.
type IPolicySys struct{}

//GlobalPolicySys policy system
var GlobalPolicySys = NewPolicySys()

// Get returns stored bucket policy
func (sys *IPolicySys) Get(bucket, accessKey string) (*Policy, error) {
	return globalBucketMetadataSys.GetPolicyConfig(bucket, accessKey)
}

// Set returns stored bucket policy
func (sys *IPolicySys) Set(bucket, accessKey string) error {
	return globalBucketMetadataSys.Set(bucket, accessKey, newBucketMetadata(bucket))
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IPolicySys) IsAllowed(args Args) bool {
	p, err := sys.Get(args.BucketName, args.AccountName)
	if err != nil {
		return false
	} else {
		return p.IsAllowed(args)
	}
}

// NewPolicySys - creates new policy system.
func NewPolicySys() *IPolicySys {
	return &IPolicySys{}
}
