package db

import (
	"context"
	"fmt"
	"sync"
)

type ClipStore interface {
	AddClip(ctx context.Context, args AddClipParams) (AddClipResults, error)
	GetClip(ctx context.Context, id int64) (*GetClipResult, error)
}

type AddClipParams struct {
	Subject   string   `json:"subject"`
	VideoUrl  string   `json:"url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}

type AddClipResults struct {
	Id        int64    `json:"id"`
	Subject   string   `json:"subject"`
	VideoUrl  string   `json:"url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}

type GetClipResult struct {
	Id        int64    `json:"id"`
	Subject   string   `json:"subject"`
	VideoUrl  string   `json:"url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}

type Clip struct {
	Id        int64    `json:"id"`
	Subject   string   `json:"subject"`
	VideoUrl  string   `json:"url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}

type InMemoryClipStore struct {
	idLock    *sync.Mutex
	currentId int64
	clips     []Clip
}

func NewInMemoryClipStore() *InMemoryClipStore {
	currId := 1
	lock := &sync.Mutex{}
	clips := make([]Clip, 0)

	return &InMemoryClipStore{idLock: lock, currentId: int64(currId), clips: clips}
}

func (s *InMemoryClipStore) getNextId() int64 {
	s.idLock.Lock()
	defer s.idLock.Unlock()

	id := s.currentId
	s.currentId++

	return id
}

func (s *InMemoryClipStore) AddClip(ctx context.Context, args AddClipParams) (AddClipResults, error) {
	c := Clip{
		Id:        s.getNextId(),
		Subject:   args.Subject,
		VideoUrl:  args.VideoUrl,
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
		Tags:      args.Tags,
	}
	s.clips = append(s.clips, c)

	clipResult := AddClipResults{
		Id:        c.Id,
		Subject:   c.Subject,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
		VideoUrl:  c.VideoUrl,
	}
	return clipResult, nil
}

func (s *InMemoryClipStore) GetClip(ctx context.Context, id int64) (*GetClipResult, error) {
	for _, c := range s.clips {
		if c.Id == id {
			res := &GetClipResult{
				Id:        c.Id,
				Subject:   c.Subject,
				VideoUrl:  c.VideoUrl,
				StartTime: c.StartTime,
				EndTime:   c.EndTime,
				Tags:      c.Tags,
			}
			return res, nil
		}
	}
	return nil, fmt.Errorf("clip with id '%d' not found", id)
}
