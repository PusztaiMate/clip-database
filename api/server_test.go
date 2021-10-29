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

	"github.com/PusztaiMate/clip-database/db"
	"github.com/PusztaiMate/clip-database/utils"
	"github.com/matryer/is"
)

const (
	clipEnpointUrl = "/clip"
)

func getServerTestLogger() *log.Logger {
	return log.New(os.Stdout, "[server-test-log] ", log.LstdFlags)
}

func TestGetClipEndpoint(t *testing.T) {
	store := db.NewInMemoryClipStore()
	logger := getServerTestLogger()
	server := NewServer(store, logger)

	t.Run("can retrieve a clip that is present in the db", func(t *testing.T) {
		is := is.New(t)
		args := db.AddClipParams{
			Subject:   utils.GetRandomString(6),
			VideoUrl:  "http://youtube.com/watch?v=1234asdf",
			StartTime: "00:10",
			EndTime:   "00:30",
			Tags:      []string{"tackle", "good"},
		}
		ctx := context.Background()
		addClipResult, err := store.AddClip(ctx, args)
		is.NoErr(err)
		is.True(addClipResult.Id > 0)

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
	})

	t.Run("not found (404) is returned if the clip is not in the db", func(t *testing.T) {
		is := is.New(t)
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, getClipUrl(9999), nil)

		server.ServeHTTP(responseRecorder, request)

		is.Equal(http.StatusNotFound, responseRecorder.Code)
	})

	t.Run("bad request (400) is returned if id is not integer", func(t *testing.T) {
		is := is.New(t)
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, getClipUrlForInvalid("hello"), nil)

		server.ServeHTTP(responseRecorder, request)

		is.Equal(http.StatusBadRequest, responseRecorder.Code)
	})
}

func TestAddClipEndpoint(t *testing.T) {
	store := db.NewInMemoryClipStore()
	logger := getServerTestLogger()
	server := NewServer(store, logger)

	type addClipRequest struct {
		Subject   string   `json:"subject,omitempty"`
		VideoUrl  string   `json:"video_url"`
		StartTime string   `json:"start_time"`
		EndTime   string   `json:"end_time"`
		Tags      []string `json:"tags"`
	}

	t.Run("can add a new clip to the db", func(t *testing.T) {
		is := is.New(t)

		recorder := httptest.NewRecorder()
		arg := addClipRequest{
			Subject:   utils.GetRandomString(10),
			VideoUrl:  utils.GetRandomString(30),
			StartTime: "01:00",
			EndTime:   "01:23",
			Tags:      []string{utils.GetRandomString(4), utils.GetRandomString(4), utils.GetRandomString(4)},
		}
		buffer := bytes.Buffer{}
		json.NewEncoder(&buffer).Encode(arg)
		request := httptest.NewRequest(http.MethodPost, addClipUrl(), &buffer)

		server.ServeHTTP(recorder, request)

		is.Equal(http.StatusCreated, recorder.Code)
	})

	t.Run("missing a required field", func(t *testing.T) {
		is := is.New(t)

		recorder := httptest.NewRecorder()
		arg := addClipRequest{
			Subject:   "",
			VideoUrl:  utils.GetRandomString(30),
			StartTime: "01:00",
			EndTime:   "01:23",
			Tags:      []string{utils.GetRandomString(4), utils.GetRandomString(4), utils.GetRandomString(4)},
		}
		buffer := bytes.Buffer{}
		json.NewEncoder(&buffer).Encode(arg)
		request := httptest.NewRequest(http.MethodPost, addClipUrl(), &buffer)

		server.ServeHTTP(recorder, request)

		is.Equal(http.StatusBadRequest, recorder.Code)
	})

}

func getClipUrl(id int64) string {
	return fmt.Sprintf("%s/%d", clipEnpointUrl, id)
}

// used for invalid url testing, kept here so update is faster
func getClipUrlForInvalid(s string) string {
	return fmt.Sprintf("%s/%s", clipEnpointUrl, s)
}

func addClipUrl() string {
	return clipEnpointUrl
}
