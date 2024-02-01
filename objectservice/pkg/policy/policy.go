package policy

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy/condition"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"golang.org/x/xerrors"
	"io"
)

// defaultVersion - default policy version as per AWS S3 specification.
const defaultVersion = "2012-10-17"

// Policy - iam bucket iamp.
type Policy struct {
	ID         ID `json:"ID"`
	Version    string
	Statements []Statement `json:"Statement"`
}

//PolicyDocument policy document
type PolicyDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// Merge merges two policies documents and drop
// duplicate statements if any.
func (p *PolicyDocument) Merge(input PolicyDocument) PolicyDocument {
	var mergedPolicy PolicyDocument
	for _, st := range p.Statement {
		mergedPolicy.Statement = append(mergedPolicy.Statement, st.Clone())
	}
	for _, st := range input.Statement {
		mergedPolicy.Statement = append(mergedPolicy.Statement, st.Clone())
	}
	mergedPolicy.dropDuplicateStatements()
	return mergedPolicy
}
func (p *PolicyDocument) dropDuplicateStatements() {
redo:
	for i := range p.Statement {
		for _, statement := range p.Statement[i+1:] {
			if !p.Statement[i].Equals(statement) {
				continue
			}
			p.Statement = append(p.Statement[:i], p.Statement[i+1:]...)
			goto redo
		}
	}
}

func (p *PolicyDocument) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (p *Policy) IsAllowed(args auth.Args) bool {
	// Check all deny statements. If any one statement denies, return false.
	for _, statement := range p.Statements {
		if statement.Effect == Deny {
			if !statement.IsAllowed(args) {
				return false
			}
		}
	}

	// For owner, it allowed by default.
	if args.IsOwner {
		// TODO: if it's admin, give it read-only permission?
		// return true
	}

	// Check all allow statements. If anyone statement allows, return true.
	for _, statement := range p.Statements {
		if statement.Effect == Allow {
			if statement.IsAllowed(args) {
				return true
			}
		}
	}

	return false
}

// ParseConfig - parses data in given reader to Policy.
func ParseConfig(reader io.Reader, bucketName string) (*Policy, error) {
	var policy Policy

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&policy); err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	err := policy.Validate(bucketName)
	return &policy, err
}

// Validate - validates all statements are for given bucket or not.
func (p *Policy) Validate(bucketName string) error {
	if err := p.isValid(); err != nil {
		return err
	}

	for _, statement := range p.Statements {
		if err := statement.Validate(bucketName); err != nil {
			return err
		}
	}

	return nil
}

// Merge merges two policies documents and drop
// duplicate statements if any.
func (p *Policy) Merge(input Policy) Policy {
	var mergedPolicy Policy
	for _, st := range p.Statements {
		mergedPolicy.Statements = append(mergedPolicy.Statements, st.Clone())
	}
	for _, st := range input.Statements {
		mergedPolicy.Statements = append(mergedPolicy.Statements, st.Clone())
	}
	mergedPolicy.dropDuplicateStatements()
	return mergedPolicy
}
func (p *Policy) dropDuplicateStatements() {
redo:
	for i := range p.Statements {
		for _, statement := range p.Statements[i+1:] {
			if !p.Statements[i].Equals(statement) {
				continue
			}
			p.Statements = append(p.Statements[:i], p.Statements[i+1:]...)
			goto redo
		}
	}
}

// Equals returns true if the two policies are identical
func (p *Policy) Equals(policy Policy) bool {
	if p.ID != policy.ID {
		return false
	}
	if len(p.Statements) != len(policy.Statements) {
		return false
	}
	for i, st := range policy.Statements {
		if !p.Statements[i].Equals(st) {
			return false
		}
	}
	return true
}

// IsEmpty - returns whether policy is empty or not.
func (p *Policy) IsEmpty() bool {
	return len(p.Statements) == 0
}

// isValid - checks if Policy is valid or not.
func (p *Policy) isValid() error {

	for _, statement := range p.Statements {
		if err := statement.IsValid(); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalJSON - decodes JSON data to Iamp.
func (p *Policy) UnmarshalJSON(data []byte) error {
	// subtype to avoid recursive call to UnmarshalJSON()
	type subPolicy Policy
	var sp subPolicy
	if err := json.Unmarshal(data, &sp); err != nil {
		return err
	}

	po := Policy(sp)
	p.dropDuplicateStatements()
	*p = po
	return nil
}

//CreateAnonReadOnlyBucketPolicy creates a bucket policy for anonymous read-only access
func CreateAnonReadOnlyBucketPolicy(bucketName string) *Policy {
	return &Policy{
		Version: defaultVersion,
		Statements: []Statement{
			NewStatement(
				"",
				Allow,
				NewPrincipal("*"),
				s3action.NewActionSet(s3action.GetBucketLocationAction, s3action.ListBucketAction),
				NewResourceSet(NewResource(bucketName, "")),
				condition.NewConFunctions(),
			),
		},
	}
}

//CreateAnonWriteOnlyBucketPolicy creates a policy that allows anonymous users to write to a bucket
func CreateAnonWriteOnlyBucketPolicy(bucketName string) *Policy {
	return &Policy{
		Version: defaultVersion,
		Statements: []Statement{
			NewStatement(
				"",
				Allow,
				NewPrincipal("*"),
				s3action.NewActionSet(
					s3action.GetBucketLocationAction,
					s3action.ListBucketMultipartUploadsAction,
				),
				NewResourceSet(NewResource(bucketName, "")),
				condition.NewConFunctions(),
			),
		},
	}
}

//CreateAnonReadOnlyObjectPolicy - create a policy for anonymous read only access to an object
func CreateAnonReadOnlyObjectPolicy(bucketName, prefix string) *Policy {
	return &Policy{
		Version: defaultVersion,
		Statements: []Statement{
			NewStatement(
				"",
				Allow,
				NewPrincipal("*"),
				s3action.NewActionSet(s3action.GetObjectAction),
				NewResourceSet(NewResource(bucketName, prefix)),
				condition.NewConFunctions(),
			),
		},
	}
}

//CreateAnonWriteOnlyObjectPolicy creates a policy that allows anonymous users to upload objects to a bucket
func CreateAnonWriteOnlyObjectPolicy(bucketName, prefix string) *Policy {
	return &Policy{
		Version: defaultVersion,
		Statements: []Statement{
			NewStatement(
				"",
				Allow,
				NewPrincipal("*"),
				s3action.NewActionSet(
					s3action.AbortMultipartUploadAction,
					s3action.DeleteObjectAction,
					s3action.ListMultipartUploadPartsAction,
					s3action.PutObjectAction,
				),
				NewResourceSet(NewResource(bucketName, prefix)),
				condition.NewConFunctions(),
			),
		},
	}
}

//CreateUserPolicy create user policy according action and bucket
func CreateUserPolicy(accessKey string, actions []s3action.Action, bucketName string) *Policy {
	return &Policy{
		Version: defaultVersion,
		Statements: []Statement{
			NewStatement(
				"",
				Allow,
				NewPrincipal(accessKey),
				s3action.NewActionSet(
					actions...,
				),
				NewResourceSet(NewResource(bucketName, "*")),
				condition.NewConFunctions(),
			),
		},
	}
}

//CreateUserBucketPolicy create user policy according accessKey and bucket
func CreateUserBucketPolicy(bucketName, accessKey string) *Policy {
	return &Policy{
		Version: defaultVersion,
		Statements: []Statement{
			NewStatement(
				"",
				Allow,
				NewPrincipal(accessKey),
				s3action.NewActionSet(
					s3action.AllActions,
				),
				NewResourceSet(NewResource(bucketName, "*")),
				condition.NewConFunctions(),
			),
		},
	}
}
