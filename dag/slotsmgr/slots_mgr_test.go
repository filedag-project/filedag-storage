package slotsmgr

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewHashSlot(t *testing.T) {
	slots := NewSlotsManager()
	slots.SetRange(SlotPair{Start: 1, End: 1000}, true)
	slots.SetRange(SlotPair{Start: 2000, End: 2100}, true)
	slots.SetRange(SlotPair{Start: 2001, End: 2101}, false)

	pairs := slots.ToSlotPair()
	t.Log(pairs)
}

type SlotPairEx struct {
	SlotPair
	Value bool
}

func TestHashSlot_ToSlotPair(t *testing.T) {
	testCases := []struct {
		name              string
		srcSlotPairs      []SlotPairEx
		expectedSlotPairs []SlotPair
	}{
		{
			name: "test1",
			srcSlotPairs: []SlotPairEx{
				{SlotPair{Start: 0, End: 200}, true},
			},
			expectedSlotPairs: []SlotPair{
				{Start: 0, End: 200},
			},
		},
		{
			name: "test2",
			srcSlotPairs: []SlotPairEx{
				{SlotPair{Start: 0, End: 200}, true},
				{SlotPair{Start: 77, End: 180}, false},
			},
			expectedSlotPairs: []SlotPair{
				{Start: 0, End: 76},
				{Start: 181, End: 200},
			},
		},
		{
			name: "test3",
			srcSlotPairs: []SlotPairEx{
				{SlotPair{Start: 0, End: 200}, true},
				{SlotPair{Start: 10001, End: 12000}, true},
			},
			expectedSlotPairs: []SlotPair{
				{Start: 0, End: 200},
				{Start: 10001, End: 12000},
			},
		},
		{
			name: "test4",
			srcSlotPairs: []SlotPairEx{
				{SlotPair{Start: 0, End: 16383}, true},
				{SlotPair{Start: 5000, End: 6342}, false},
				{SlotPair{Start: 8000, End: 13000}, false},
			},
			expectedSlotPairs: []SlotPair{
				{Start: 0, End: 4999},
				{Start: 6343, End: 7999},
				{Start: 13001, End: 16383},
			},
		},
		{
			name: "test5",
			srcSlotPairs: []SlotPairEx{
				{SlotPair{Start: 0, End: 234}, true},
				{SlotPair{Start: 366, End: 5678}, true},
				{SlotPair{Start: 8000, End: 13000}, true},
				{SlotPair{Start: 5700, End: 5801}, true},
				{SlotPair{Start: 5820, End: 7342}, true},
			},
			expectedSlotPairs: []SlotPair{
				{Start: 0, End: 234},
				{Start: 366, End: 5678},
				{Start: 8000, End: 13000},
				{Start: 5700, End: 5801},
				{Start: 5820, End: 7342},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(s *testing.T) {
			slots := NewSlotsManager()
			for _, pair := range tc.srcSlotPairs {
				if err := slots.SetRange(pair.SlotPair, pair.Value); err != nil {
					t.Fatal(err)
				}
			}

			pairs := slots.ToSlotPair()
			require.Equal(s, len(tc.expectedSlotPairs), len(pairs))
			for _, pair := range pairs {
				require.Contains(s, tc.expectedSlotPairs, pair)
			}
		})
	}
}

func TestHashSlot_Count(t *testing.T) {
	slots := NewSlotsManager()
	slots.SetRange(SlotPair{Start: 1, End: 1000}, true)
	slots.SetRange(SlotPair{Start: 2000, End: 2100}, true)
	require.Equal(t, uint64(1101), slots.Count())
}
