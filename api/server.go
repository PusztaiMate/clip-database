package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PusztaiMate/clip-database/db"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store     db.ClipStore
	router    *chi.Mux
	logger    *log.Logger
	validator *validator.Validate
}

func NewServer(store db.ClipStore, logger *log.Logger) *Server {
	r := chi.NewRouter()
	v := validator.New()

	server := Server{
		store:     store,
		logger:    logger,
		router:    r,
		validator: v,
	}

	r.Get("/clip/{clipId}", server.getSingleClip)
	r.Post("/clip", server.addClip)

	return &server
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) getSingleClip(rw http.ResponseWriter, r *http.Request) {
	clipId, err := getClipIdFromRequest(r)
	if err != nil {
		s.writeJsonMessage(rw, "the provided id is not valid", http.StatusBadRequest)
		return
	}

	c, err := s.store.GetClip(r.Context(), clipId)
	if err != nil {
		s.writeJsonMessage(rw, fmt.Sprintf("player with id '%d' not found", clipId), http.StatusNotFound)
		return
	}

	clipResponse := GetClipResponse{
		Id:        c.Id,
		Subject:   c.Subject,
		Url:       c.VideoUrl,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
		Tags:      c.Tags,
	}

	_ = s.writeJsonOutput(rw, clipResponse, http.StatusOK)
}

type AddClipRequest struct {
	Subject   string   `json:"subject" validate:"required"`
	Url       string   `json:"video_url" validate:"required"`
	StartTime string   `json:"start_time" validate:"required"`
	EndTime   string   `json:"end_time" validate:"required"`
	Tags      []string `json:"tags" validate:"required"`
}

type AddClipResponse struct {
	Id int64 `json:"id"`
}

func (s *Server) addClip(rw http.ResponseWriter, r *http.Request) {
	var addClipRequest AddClipRequest
	json.NewDecoder(r.Body).Decode(&addClipRequest)

	err := s.validator.Struct(addClipRequest)
	if err != nil {
		s.writeJsonMessage(rw, fmt.Sprintf("invalid payload: %s", err), http.StatusBadRequest)
		return
	}

	c, _ := s.store.AddClip(r.Context(), db.AddClipParams{
		Subject:   addClipRequest.Subject,
		VideoUrl:  addClipRequest.Url,
		StartTime: addClipRequest.StartTime,
		EndTime:   addClipRequest.EndTime,
		Tags:      addClipRequest.Tags,
	})

	response := AddClipResponse{
		Id: c.Id,
	}
	s.writeJsonOutput(rw, response, http.StatusCreated)
}

func getClipIdFromRequest(r *http.Request) (int64, error) {
	clipIdAsString := chi.URLParam(r, "clipId")
	clipId, err := strconv.Atoi(clipIdAsString)
	if err != nil {
		return 0, err
	}
	return int64(clipId), nil
}

func (s *Server) writeJsonOutput(rw http.ResponseWriter, v interface{}, code int) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	return json.NewEncoder(rw).Encode(v)
}

func (s *Server) writeJsonMessage(rw http.ResponseWriter, message string, code int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	type m struct {
		Message string `json:"message"`
	}

	msg := m{message}

	err := json.NewEncoder(rw).Encode(msg)
	if err != nil {
		s.logger.Print("could not log message: ", message)
	}
}

type GetClipResponse struct {
	Id        int64    `json:"id"`
	Subject   string   `json:"subject"`
	Url       string   `json:"url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}
