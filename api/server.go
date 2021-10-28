package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Server struct {
	store  ClipStore
	router *chi.Mux
	logger *log.Logger
}

func NewServer(store ClipStore, logger *log.Logger) *Server {
	router := chi.NewRouter()

	server := Server{
		store:  store,
		logger: logger,
		router: router,
	}

	router.Get("/clip/{clipId}", server.getSingleClip)
	router.Post("/clip", server.addClip)

	return &server
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) getSingleClip(rw http.ResponseWriter, r *http.Request) {
	clipId, _ := getClipIdFromRequest(r)

	c, _ := s.store.GetClip(r.Context(), clipId)

	clipResponse := GetClipResponse{
		Id:        c.Id,
		Subject:   c.Subject,
		Url:       c.VideoUrl,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
		Tags:      c.Tags,
	}

	_ = writeJsonOutput(rw, clipResponse, http.StatusOK)
}

type AddClipRequest struct {
	Subject   string   `json:"subject"`
	Url       string   `json:"url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}

type AddClipResponse struct {
	Id int64 `json:"id"`
}

func (s *Server) addClip(rw http.ResponseWriter, r *http.Request) {
	var addClipRequest AddClipRequest
	json.NewDecoder(r.Body).Decode(&addClipRequest)

	c, _ := s.store.AddClip(r.Context(), AddClipParams{
		Subject:   addClipRequest.Subject,
		VideoUrl:  addClipRequest.Url,
		StartTime: addClipRequest.StartTime,
		EndTime:   addClipRequest.EndTime,
		Tags:      addClipRequest.Tags,
	})

	response := AddClipResponse{
		Id: c.Id,
	}
	writeJsonOutput(rw, response, http.StatusCreated)
}

func getClipIdFromRequest(r *http.Request) (int64, error) {
	clipIdAsString := chi.URLParam(r, "clipId")
	clipId, err := strconv.Atoi(clipIdAsString)
	if err != nil {
		return 0, err
	}
	return int64(clipId), nil
}

func writeJsonOutput(rw http.ResponseWriter, v interface{}, code int) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	return json.NewEncoder(rw).Encode(v)
}

type GetClipResponse struct {
	Id        int64    `json:"id"`
	Subject   string   `json:"subject"`
	Url       string   `json:"url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}
