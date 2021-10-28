package api

import (
	"bytes"
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
	NEEDS_DB = "NEEDS_DB"
)

func getComponentTestLogger() *log.Logger {
	return log.New(os.Stdout, "[component-test-log]", log.LstdFlags)
}

func TestCreateAndRetrieveSingleClip(t *testing.T) {
	// We create a new clip via a POST request with clip data and then retrieve it using the ID via a GET request
	// create new clip
	is := is.New(t)

	store := NewInMemoryClipStore()
	logger := getComponentTestLogger()
	s := NewServer(store, logger)

	payloadJson := `{
		"subject": "John Doe",
		"tags": ["good", "tackle"],
		"video_url": "http://youtube.com/watch?v=thisistheid123456,
		"start_time": "00:17",
		"end_time": "01:43"
	}`
	payload := bytes.NewBufferString(payloadJson)
	createRequest := httptest.NewRequest(http.MethodPost, "/clip", payload)
	createResponseRecorder := httptest.NewRecorder()

	s.ServeHTTP(createResponseRecorder, createRequest)

	var createResponse GetClipResponse
	is.Equal(201, createResponseRecorder.Code)
	err := json.NewDecoder(createResponseRecorder.Body).Decode(&createResponse)
	is.NoErr(err)
	is.True(createResponse.Id > 0)
	is.Equal("application/json", createResponseRecorder.Result().Header.Get("Content-Type"))

	// retrieve it
	retrieveRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/clip/%d", createResponse.Id), nil)
	retrieveResponseRecorder := httptest.NewRecorder()

	s.ServeHTTP(retrieveResponseRecorder, retrieveRequest)

	var retrieveResponse AddClipResponse

	err = json.NewDecoder(retrieveResponseRecorder.Body).Decode(&retrieveResponse)
	is.NoErr(err)
	is.Equal(200, retrieveResponseRecorder.Code)
	is.True(0 < retrieveResponse.Id)
}
