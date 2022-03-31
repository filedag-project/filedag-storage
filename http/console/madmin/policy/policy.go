package policy

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"io"
)

// Policy - iam bucket iamp.
type Policy struct {
	ID         ID `json:"ID,omitempty"`
	Version    string
	Statements []Statement `json:"Statement"`
}
type PolicyDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}
type Policies struct {
	Policies map[string]PolicyDocument `json:"policies"`
}

func (p PolicyDocument) String() string {
	b, _ := json.Marshal(p)
	return string(b)
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
func (p Policy) Validate(bucketName string) error {
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
func (p Policy) IsEmpty() bool {
	return len(p.Statements) == 0
}

// isValid - checks if Policy is valid or not.
func (p Policy) isValid() error {

	for _, statement := range p.Statements {
		if err := statement.IsValid(); err != nil {
			return err
		}
	}
	return nil
}
