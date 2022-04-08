package policy

// PolicySys - policy subsystem.
type PolicySys struct{}

//GlobalPolicySys policy system
var GlobalPolicySys = NewPolicySys()

// Get returns stored bucket policy
func (sys *PolicySys) Get(bucket string) (*Policy, error) {
	return globalBucketMetadataSys.GetPolicyConfig(bucket)
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *PolicySys) IsAllowed(args Args) bool {
	p, err := sys.Get(args.BucketName)
	if err != nil {
		return false
	} else {
		return p.IsAllowed(args)
	}
}

// NewPolicySys - creates new policy system.
func NewPolicySys() *PolicySys {
	return &PolicySys{}
}
