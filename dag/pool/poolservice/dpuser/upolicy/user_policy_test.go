package upolicy

import "testing"

func TestDagPoolPolicy_Allow(t *testing.T) {
	var (
		a = ReadOnly
		b = WriteOnly
		c = ReadWrite
	)
	testCases := []struct {
		pol    DagPoolPolicy
		status bool
	}{
		{ReadOnly, true},
		{ReadWrite, false},
		{WriteOnly, false},
	}
	for i, test := range testCases {
		if a.Allow(test.pol) != test.status {
			t.Fatalf("OnlyRead case %v ,fail", i)
		}
	}
	testCases = []struct {
		pol    DagPoolPolicy
		status bool
	}{
		{ReadOnly, false},
		{ReadWrite, false},
		{WriteOnly, true},
	}
	for i, test := range testCases {
		if b.Allow(test.pol) != test.status {
			t.Fatalf("OnlyWrite case %v ,fail", i)
		}
	}
	testCases = []struct {
		pol    DagPoolPolicy
		status bool
	}{
		{ReadOnly, true},
		{ReadWrite, true},
		{WriteOnly, true},
	}
	for i, test := range testCases {
		if c.Allow(test.pol) != test.status {
			t.Fatalf("ReadWrite case %v ,fail", i)
		}

	}
}
