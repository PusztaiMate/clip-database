package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/PusztaiMate/clip-database/db"
	"github.com/matryer/is"
)

const (
	ADD_CLIP = "addClip"
	GET_CLIP = "getClip"
)

type SpyStore struct {
	CallCount map[string]int
}

func NewSpyStore() *SpyStore {
	cc := make(map[string]int)
	return &SpyStore{CallCount: cc}
}

func (ss *SpyStore) AddClip(ctx context.Context, args db.AddClipParams) (db.AddClipResults, error) {
	ss.CallCount[ADD_CLIP]++

	return db.AddClipResults{}, nil
}

func (ss *SpyStore) GetClip(ctx context.Context, id int64) (*db.GetClipResult, error) {
	ss.CallCount[GET_CLIP]++

	return &db.GetClipResult{}, nil
}

type FaultyStore struct {
}

func (fs *FaultyStore) AddClip(ctx context.Context, args db.AddClipParams) (db.AddClipResults, error) {
	return db.AddClipResults{}, errors.New("error while saving clip to db")
}

func (fs *FaultyStore) GetClip(ctx context.Context, id int64) (*db.GetClipResult, error) {
	return &db.GetClipResult{}, errors.New("error while retrieving clip from db")
}

func TestAddClip(t *testing.T) {
	t.Run("service calls the store to save the clip", func(t *testing.T) {
		is := is.New(t)
		// new store and service, so we can be sure that this is the first call
		store := NewSpyStore()
		service := NewClipperService(store)
		is.Equal(0, store.CallCount[ADD_CLIP])

		_, err := service.AddClip(context.Background(), AddClipParams{})
		is.NoErr(err)
		is.Equal(1, store.CallCount[ADD_CLIP])
	})

	t.Run("error returned when db returns with error", func(t *testing.T) {
		is := is.New(t)
		store := &FaultyStore{}
		service := NewClipperService(store)

		_, err := service.AddClip(context.Background(), AddClipParams{})
		is.True(err != nil)
		is.True(strings.Contains(err.Error(), "error while saving clip to db"))
	})
}

func TestGetClip(t *testing.T) {
	t.Run("service calls the store to retrieve the clip", func(t *testing.T) {
		is := is.New(t)
		// new store and service, so we can be sure that this is the first call
		store := NewSpyStore()
		service := NewClipperService(store)
		is.Equal(0, store.CallCount[ADD_CLIP])

		// using id=0, because so far no check is implemented, so this should fail later on
		// as the service layer should catch this error
		service.GetClip(context.Background(), 0)
		is.Equal(1, store.CallCount[GET_CLIP])
	})

	t.Run("error returned when db returns with error", func(t *testing.T) {
		is := is.New(t)
		store := &FaultyStore{}
		service := NewClipperService(store)

		_, err := service.GetClip(context.Background(), 0)
		is.True(err != nil)
		is.True(strings.Contains(err.Error(), "error while retrieving clip from db"))
	})
}
