package policy

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/set"
	"golang.org/x/xerrors"
	"path"
	"strings"
)

// resourceARNPrefix - resource ARN prefix as per AWS S3 specification.
const resourceARNPrefix = "arn:aws:s3:::"

//The Resource element specifies the object or objects that the statement covers. Statements must include either a Resource or a NotResource element.
// An entity that users can work with in AWS, such as an EC2 instance, an Amazon DynamoDB table, an Amazon S3 bucket, an IAM user, or an AWS OpsWorks stack.
//Resource - resource in policy statement.
//"Resource": "arn:aws:iam::account-ID-without-hyphens:user/accounting/*"
type Resource struct {
	BucketName string
	Pattern    string
}

// IsValid - checks whether Resource is valid or not.
func (r Resource) IsValid() bool {
	return r.BucketName != "" && r.Pattern != ""
}

// Match - matches object name with resource pattern, including specific conditionals.
func (r Resource) Match(resource string, conditionValues map[string][]string) bool {
	pattern := r.Pattern
	for _, key := range CommonKeys {
		// Empty values are not supported for policy variables.
		if rvalues, ok := conditionValues[key.Name()]; ok && rvalues[0] != "" {
			pattern = strings.Replace(pattern, key.VarName(), rvalues[0], -1)
		}
	}
	if cp := path.Clean(resource); cp != "." && cp == pattern {
		return true
	}
	return set.Match(pattern, resource)
}

// MarshalJSON - encodes Resource to JSON data.
func (r Resource) MarshalJSON() ([]byte, error) {
	if !r.IsValid() {
		return nil, xerrors.Errorf("invalid resource %v", r)
	}

	return json.Marshal(r.String())
}

func (r Resource) String() string {
	return resourceARNPrefix + r.Pattern
}

// UnmarshalJSON - decodes JSON data to Resource.
func (r *Resource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsedResource, err := parseResource(s)
	if err != nil {
		return err
	}

	*r = parsedResource

	return nil
}

// Validate - validates Resource is for given bucket or not.
func (r Resource) Validate(bucketName string) error {
	if !r.IsValid() {
		return xerrors.Errorf("invalid resource")
	}

	if !set.Match(r.BucketName, bucketName) {
		return xerrors.Errorf("bucket name does not match")
	}

	return nil
}

// parseResource - parses string to Resource.
func parseResource(s string) (Resource, error) {
	if !strings.HasPrefix(s, resourceARNPrefix) {
		return Resource{}, xerrors.Errorf("invalid resource '%v'", s)
	}

	pattern := strings.TrimPrefix(s, resourceARNPrefix)
	tokens := strings.SplitN(pattern, "/", 2)
	bucketName := tokens[0]
	if bucketName == "" {
		return Resource{}, xerrors.Errorf("invalid resource format '%v'", s)
	}

	return Resource{
		BucketName: bucketName,
		Pattern:    pattern,
	}, nil
}

// NewResource - creates new resource.
func NewResource(bucketName, keyName string) Resource {
	pattern := bucketName
	if keyName != "" {
		if !strings.HasPrefix(keyName, "/") {
			pattern += "/"
		}

		pattern += keyName
	}

	return Resource{
		BucketName: bucketName,
		Pattern:    pattern,
	}
}
