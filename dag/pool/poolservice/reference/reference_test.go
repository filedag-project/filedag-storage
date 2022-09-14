package reference

import (
	"context"
	"errors"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCacheSet(t *testing.T) {
	db, err := uleveldb.OpenDb(t.TempDir())
	require.NoError(t, err)
	cset := NewCacheSet(db)
	testKeys := []string{
		"a",
		"b354646rt23sdsfddfx",
		"cffdf",
		"132d132sds",
	}

	for _, key := range testKeys {
		err := cset.Add(key)
		require.NoError(t, err)
	}

	for _, key := range testKeys {
		has, err := cset.Has(key)
		require.NoError(t, err)
		require.True(t, has)
	}

	has, err := cset.Has("not_exist")
	require.NoError(t, err)
	require.False(t, has)

	keysCh, err := cset.AllKeysChan(context.TODO())
	require.NoError(t, err)
	for key := range keysCh {
		require.Contains(t, testKeys, key)
	}
	for _, key := range testKeys {
		err = cset.Remove(key)
		require.NoError(t, err)
	}
}

func TestRefCounter(t *testing.T) {
	db, err := uleveldb.OpenDb(t.TempDir())
	require.NoError(t, err)
	cset := NewCacheSet(db)
	counter := NewRefCounter(db, cset)
	testKeys := []string{
		"a",
		"b354646rt23sdsfddfx",
		"cffdf",
		"132d132sds",
	}

	for i, key := range testKeys {
		if i == 0 {
			expError := errors.New("test error")
			err := counter.IncrOrCreate(key, func() error {
				return expError
			})
			require.EqualError(t, err, expError.Error())
		}
		err := counter.IncrOrCreate(key, func() error {
			return nil
		})
		require.NoError(t, err)
	}

	for _, key := range testKeys {
		err := counter.Incr(key)
		require.NoError(t, err)
	}

	for _, key := range testKeys {
		num, err := counter.Get(key)
		require.NoError(t, err)
		require.Equal(t, int64(2), num, "key=%s", key)
	}

	for _, key := range testKeys {
		err := counter.Decr(key)
		require.NoError(t, err)
	}

	for _, key := range testKeys {
		num, err := counter.Get(key)
		require.NoError(t, err)
		require.Equal(t, int64(1), num)
	}

	for _, key := range testKeys {
		has, err := counter.Has(key)
		require.NoError(t, err)
		require.True(t, has)
	}

	has, err := counter.Has("not_exist")
	require.NoError(t, err)
	require.False(t, has)

	keysCh, err := counter.AllKeysChan(context.TODO(), 1)
	require.NoError(t, err)
	for key := range keysCh {
		require.Contains(t, testKeys, key)
	}

	for _, key := range testKeys {
		err := counter.Decr(key)
		require.NoError(t, err)
	}

	extKey := "remove key"
	err = counter.Incr(extKey)
	require.NoError(t, err)
	err = counter.Remove(extKey, false)
	require.Error(t, err)
	err = counter.Remove(extKey, true)
	require.NoError(t, err)

	keysCh2, err := cset.AllKeysChan(context.TODO())
	require.NoError(t, err)
	testKeys = append(testKeys, extKey)
	for key := range keysCh2 {
		require.Contains(t, testKeys, key)
	}
}
