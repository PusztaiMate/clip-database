package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PusztaiMate/clip-database/service"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	service   *service.ClipperService
	router    *chi.Mux
	logger    *log.Logger
	validator *validator.Validate
}

func NewServer(service *service.ClipperService, logger *log.Logger) *Server {
	r := chi.NewRouter()
	v := validator.New()

	server := Server{
		service:   service,
		logger:    logger,
		router:    r,
		validator: v,
	}

	r.Get("/clip/{clipId}", server.handlerGetSingleClip)
	r.Post("/clip", server.handlerAddClip)

	return &server
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

type getClipResponse struct {
	Id        int64    `json:"id"`
	Subject   string   `json:"subject"`
	VideoUrl  string   `json:"video_url"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Tags      []string `json:"tags"`
}

func (s *Server) handlerGetSingleClip(rw http.ResponseWriter, r *http.Request) {
	clipId, err := getClipIdFromRequest(r)
	if err != nil {
		s.writeJsonMessage(rw, "the provided id is not valid", http.StatusBadRequest)
		return
	}

	c, err := s.service.GetClip(r.Context(), clipId)
	if err != nil {
		s.writeJsonMessage(rw, fmt.Sprintf("player with id '%d' not found", clipId), http.StatusNotFound)
		return
	}

	clipResponse := getClipResponse{
		Id:        c.Id,
		Subject:   c.Subject,
		VideoUrl:  c.VideoUrl,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
		Tags:      c.Tags,
	}

	_ = s.writeJsonOutput(rw, clipResponse, http.StatusOK)
}

type addClipRequest struct {
	Subject   string   `json:"subject" validate:"required"`
	VideoUrl  string   `json:"video_url" validate:"required"`
	StartTime string   `json:"start_time" validate:"required"`
	EndTime   string   `json:"end_time" validate:"required"`
	Tags      []string `json:"tags" validate:"required"`
}

type addClipResponse struct {
	Id int64 `json:"id"`
}

func (s *Server) handlerAddClip(rw http.ResponseWriter, r *http.Request) {
	var addClipRequest addClipRequest
	json.NewDecoder(r.Body).Decode(&addClipRequest)

	err := s.validator.Struct(addClipRequest)
	if err != nil {
		s.writeJsonMessage(rw, fmt.Sprintf("invalid payload: %s", err), http.StatusBadRequest)
		return
	}

	c, _ := s.service.AddClip(r.Context(), service.AddClipParams{
		Subject:   addClipRequest.Subject,
		VideoUrl:  addClipRequest.VideoUrl,
		StartTime: addClipRequest.StartTime,
		EndTime:   addClipRequest.EndTime,
		Tags:      addClipRequest.Tags,
	})

	response := addClipResponse{
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
