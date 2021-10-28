package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"
)

const (
	clipEnpointUrl = "/clip"
)

func getServerTestLogger() *log.Logger {
	return log.New(os.Stdout, "[server-test-log] ", log.LstdFlags)
}

func TestGetClipEndpoint(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryClipStore()
	args := AddClipParams{
		Subject:   "John Doe " + getRandomString(6),
		VideoUrl:  "http://youtube.com/watch?v=1234asdf",
		StartTime: "00:10",
		EndTime:   "00:30",
		Tags:      []string{"tackle", "good"},
	}
	ctx := context.Background()
	addClipResult, err := store.AddClip(ctx, args)

	is.NoErr(err)
	is.True(addClipResult.Id > 0)

	server := NewServer(store, getServerTestLogger())
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, getClipUrl(addClipResult.Id), nil)
	server.ServeHTTP(responseRecorder, request)

	var respone GetClipResponse
	err = json.NewDecoder(responseRecorder.Body).Decode(&respone)
	is.NoErr(err)
	is.Equal(args.Subject, respone.Subject)
	is.Equal(args.VideoUrl, respone.Url)
	is.Equal(args.StartTime, respone.StartTime)
	is.Equal(args.EndTime, respone.EndTime)
	is.Equal(args.Tags, respone.Tags)
}

func TestAddClipEndpoint(t *testing.T) {
	store := NewInMemoryClipStore()
	logger := getServerTestLogger()
	server := NewServer(store, logger)

	t.Run("can add a new clip to the db", func(t *testing.T) {
		is := is.New(t)

		recorder := httptest.NewRecorder()
		type addClipRequest struct {
			Subject   string   `json:"subject"`
			VideoUrl  string   `json:"url"`
			StartTime string   `json:"start_time"`
			EndTime   string   `json:"end_time"`
			Tags      []string `json:"tags"`
		}
		arg := addClipRequest{
			Subject:   getRandomString(10),
			VideoUrl:  getRandomString(30),
			StartTime: "01:00",
			EndTime:   "01:23",
			Tags:      []string{getRandomString(4), getRandomString(4), getRandomString(4)},
		}
		buffer := bytes.Buffer{}
		json.NewEncoder(&buffer).Encode(arg)
		request := httptest.NewRequest(http.MethodPost, clipEnpointUrl, &buffer)

		server.ServeHTTP(recorder, request)

		is.Equal(http.StatusCreated, recorder.Code)
	})
}

func getClipUrl(id int64) string {
	return fmt.Sprintf("%s/%d", clipEnpointUrl, id)
}
