package service

import (
	"context"
	"fmt"

	"github.com/PusztaiMate/clip-database/db"
)

type ClipperService struct {
	store db.ClipStore
}

func NewClipperService(store db.ClipStore) *ClipperService {
	return &ClipperService{store: store}
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

func (cs *ClipperService) AddClip(ctx context.Context, arg AddClipParams) (AddClipResults, error) {
	result, err := cs.store.AddClip(ctx, db.AddClipParams{
		Subject:   arg.Subject,
		VideoUrl:  arg.VideoUrl,
		StartTime: arg.StartTime,
		EndTime:   arg.EndTime,
		Tags:      arg.Tags,
	})

	if err != nil {
		return AddClipResults{}, fmt.Errorf("could not add clip to db: %v", err)
	}

	r := AddClipResults{
		Id:        result.Id,
		Subject:   result.Subject,
		VideoUrl:  result.VideoUrl,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
		Tags:      result.Tags,
	}

	return r, nil
}

func (cs *ClipperService) GetClip(ctx context.Context, id int64) (GetClipResult, error) {
	result, err := cs.store.GetClip(ctx, id)
	if err != nil {
		return GetClipResult{}, fmt.Errorf("could not retrieve clip with id '%d' from db: '%v'", id, err)
	}

	r := GetClipResult{
		Id:        result.Id,
		Subject:   result.Subject,
		VideoUrl:  result.VideoUrl,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
		Tags:      result.Tags,
	}

	return r, nil
}
