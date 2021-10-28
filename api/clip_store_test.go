package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func addRandomClip(t *testing.T, store ClipStore) AddClipResults {
	t.Helper()
	is := is.New(t)

	ctx := context.Background()
	arg := AddClipParams{
		Subject:   getRandomString(8),
		VideoUrl:  fmt.Sprintf("http://youtube.com/watch?v=%s", getRandomString(8)),
		StartTime: "01:00",
		EndTime:   "03:00",
		Tags:      []string{getRandomString(6), getRandomString(6)},
	}

	result, err := store.AddClip(ctx, arg)
	is.NoErr(err)
	is.Equal(arg.Subject, result.Subject)
	is.True(0 < result.Id)
	is.Equal(arg.StartTime, result.StartTime)
	is.Equal(arg.EndTime, result.EndTime)
	is.Equal(arg.VideoUrl, result.VideoUrl)

	return result
}

func TestInMemoryClipStoreAdd(t *testing.T) {
	t.Run("can add a single clip", func(t *testing.T) {
		store := NewInMemoryClipStore()
		addRandomClip(t, store)
	})
}

func TestInMemoryClipStoreGet(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryClipStore()

	t.Run("can retrieve clip added earlier", func(t *testing.T) {
		is := is.New(t)
		clip1 := addRandomClip(t, store)

		getClipResult, err := store.GetClip(ctx, clip1.Id)
		is.NoErr(err)
		is.Equal(clip1.Id, getClipResult.Id)
	})

	t.Run("error given if the id is not found", func(t *testing.T) {
		is := is.New(t)

		getClipResult, err := store.GetClip(ctx, 9999999)
		is.True(err != nil)
		is.True(getClipResult == nil)
	})
}
