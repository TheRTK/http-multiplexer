package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TheRTK/http-multiplexer/internal/app"
	"github.com/TheRTK/http-multiplexer/internal/config"
	"github.com/TheRTK/http-multiplexer/internal/services/request"
)

type Server struct {
	httpServer *http.Server

	requestService request.IRequestService
	cfg            *config.Config
}

func New(opts ...app.ServerOption) *Server {
	srv := &Server{}

	for _, option := range opts {
		option(srv)
	}

	return srv
}

func (s *Server) GetAppOptions() []app.Option {
	return []app.Option{
		app.ServerConnector(s),
	}
}

func (s *Server) SetRequestService(rs request.IRequestService) {
	s.requestService = rs
}

func (s *Server) GetRequestService() request.IRequestService {
	return s.requestService
}

func (s *Server) SetConfig(cfg *config.Config) {
	s.cfg = cfg
}

func (s *Server) GetConfig() *config.Config {
	return s.cfg
}


func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:           ":" + s.cfg.PortHTTP,
		ReadTimeout:    time.Second * 15,
		WriteTimeout:   time.Second * 15,
		MaxHeaderBytes: 1 << 20, // 1MB
		Handler:        NewHandler(s.GetAppOptions, s.cfg.RequestLimitCount),
	}

	fmt.Printf("Starting server at port %s\n", s.cfg.PortHTTP)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func extractRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

type ResponseWithMessage struct {
	Message string `json:"message"`
}

func writeResponseJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	response, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(""))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(response)
}
