package v2

import (
	"encoding/json"
	"fmt"
	"github.com/cephxdev/nero/repo"
	"github.com/cephxdev/nero/server/api"
	"github.com/cephxdev/nero/server/api/nekos/v2"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"net/http"
	"net/url"
)

var (
	DefaultRequestErrorHandler api.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		e := v2.Error{Code: http.StatusBadRequest, Message: err.Error()}
		if err := json.NewEncoder(w).Encode(e); err != nil {
			_, _ = fmt.Fprintf(w, "{\"code\":\"%d\",\"message\":\"%s\"}", http.StatusInternalServerError, "failed to serialize error")
		}
	}

	DefaultResponseErrorHandler api.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		e := v2.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		if err := json.NewEncoder(w).Encode(e); err != nil {
			_, _ = fmt.Fprintf(w, "{\"code\":\"%d\",\"message\":\"%s\"}", http.StatusInternalServerError, "failed to serialize error")
		}
	}
)

// Server is a REST server for the nekos v2 API.
type Server struct {
	repos   map[string]*repo.Repository
	baseURL *url.URL
	logger  *zap.Logger
}

// NewServer creates a new server with pre-defined repositories.
func NewServer(repos []*repo.Repository, baseURL *url.URL, logger *zap.Logger) (*Server, error) {
	reposById := make(map[string]*repo.Repository, len(repos))
	for _, r := range repos {
		repoId := r.ID()
		if _, ok := reposById[repoId]; ok {
			return nil, fmt.Errorf("duplicate repository ID %s", repoId)
		}

		reposById[repoId] = r
	}

	return &Server{
		repos:   reposById,
		baseURL: baseURL,
		logger:  logger,
	}, nil
}

// NewRouter creates a new nekos v2 API router.
func NewRouter(handler v2.StrictServerInterface) http.Handler {
	h := v2.NewStrictHandlerWithOptions(handler, nil, v2.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  DefaultRequestErrorHandler,
		ResponseErrorHandlerFunc: DefaultResponseErrorHandler,
	})

	return v2.HandlerWithOptions(h, v2.ChiServerOptions{ErrorHandlerFunc: DefaultRequestErrorHandler})
}

// Repos returns all repositories available to the server.
func (s *Server) Repos() []*repo.Repository {
	return maps.Values(s.repos)
}

// BaseURL returns the base URL of the server.
func (s *Server) BaseURL() *url.URL {
	return s.baseURL
}
