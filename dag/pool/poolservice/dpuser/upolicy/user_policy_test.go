package upolicy

import "testing"

func TestDagPoolPolicy_Allow(t *testing.T) {
	var (
		a = OnlyRead
		b = OnlyWrite
		c = ReadWrite
	)
	testCases := []struct {
		pol    DagPoolPolicy
		status bool
	}{
		{OnlyRead, true},
		{ReadWrite, false},
		{OnlyWrite, false},
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
		{OnlyRead, false},
		{ReadWrite, false},
		{OnlyWrite, true},
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
		{OnlyRead, true},
		{ReadWrite, true},
		{OnlyWrite, true},
	}
	for i, test := range testCases {
		if c.Allow(test.pol) != test.status {
			t.Fatalf("ReadWrite case %v ,fail", i)
		}

	}
}
