package policy

import (
	"testing"
)

func TestPrincipal_Match(t *testing.T) {

	p := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be")

	if !p.Match("arn:aws:iam::123456789012:root") {
		t.Error("principal match failed")
	}
	if !p.Match("999999999999") {
		t.Error("principal match failed")
	}
	if !p.Match("CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be") {
		t.Error("principal match failed")
	}
	if p.Match("arn:aws:iam::123456789012:user") {
		t.Error("principal match failed")
	}
	if p.Match("arn:aws:iam::123456789012:role") {
		t.Error("principal match failed")
	}
	if p.Match("arn:aws:iam::123456789012:group") {
		t.Error("principal match failed")
	}
	if p.Match("arn:aws:iam::123456789012:policy") {
		t.Error("principal match failed")
	}
	if p.Match("arn:aws:iam::123456789012:root/user") {
		t.Error("principal match failed")
	}
	if p.Match("arn:aws:iam::123456789012:root/role") {
		t.Error("principal match failed")
	}

}
func TestPrincipal_IsValid(t *testing.T) {

	p := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be")

	if !p.IsValid() {
		t.Error("principal is invalid")
	}
	p = NewPrincipal()
	if p.IsValid() {
		t.Error("principal is invalid")
	}
	p = NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012:user")
	if !p.IsValid() {
		t.Error("principal is invalid")
	}
}
func TestPrincipal_Equals(t *testing.T) {

	p := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be")
	p2 := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be")
	if !p.Equals(p2) {
		t.Error("principal equals failed")
	}
	p2 = NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012:user")
	if p.Equals(p2) {
		t.Error("principal equals failed")
	}
	p2 = NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012:role")
	if p.Equals(p2) {
		t.Error("principal equals failed")
	}
	p2 = NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012:group")
}
func TestPrincipal_Intersection(t *testing.T) {
	p1 := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be")
	p2 := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012:user")
	p3 := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012:role")
	p4 := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012:group")
	p5 := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012")
	p6 := NewPrincipal("arn:aws:iam::123456789012:root", "999999999999", "CanonicalUser:79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be", "arn:aws:iam::123456789012", "arn:aws:iam::123456789012:user")
	r1 := p1.Intersection(p2)
	if r1.Equals(p2.AWS) {
		t.Error("principal intersection failed")
	}
	r2 := p2.Intersection(p3)
	if r2.Equals(p3.AWS) {
		t.Error("principal intersection failed")
	}
	r3 := p3.Intersection(p4)
	if r3.Equals(p4.AWS) {
		t.Error("principal intersection failed")
	}
	r4 := p4.Intersection(p5)
	if r4.Equals(p4.AWS) {
		t.Error("principal intersection failed")
	}
	r5 := p5.Intersection(p6)
	if !r5.Equals(p5.AWS) {
		t.Error("principal intersection failed")
	}
}
