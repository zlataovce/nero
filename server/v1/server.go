package v1

import (
	"encoding/json"
	"fmt"
	"github.com/zlataovce/nero/internal/errors"
	"github.com/zlataovce/nero/repo"
	"github.com/zlataovce/nero/server/api"
	"github.com/zlataovce/nero/server/api/v1"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"net/http"
)

var (
	DefaultRequestErrorHandler api.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		e := v1.Error{Type: v1.BadRequest, Description: err.Error()}
		if err := json.NewEncoder(w).Encode(e); err != nil {
			_, _ = fmt.Fprintf(w, "{\"type\":\"%s\",\"description\":\"%s\"}", v1.InternalError, "failed to serialize error")
		}
	}

	DefaultResponseErrorHandler api.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")

		var (
			status = http.StatusInternalServerError
			type_  = v1.InternalError

			httpErr *api.HTTPError
		)
		if errors.As(err, &httpErr) {
			status = httpErr.Status

			if httpErr.Type != "" {
				type_ = v1.ErrorType(httpErr.Type)
			}
		}

		w.WriteHeader(status)

		e := v1.Error{Type: type_, Description: err.Error()}
		if err := json.NewEncoder(w).Encode(e); err != nil {
			_, _ = fmt.Fprintf(w, "{\"type\":\"%s\",\"description\":\"%s\"}", v1.InternalError, "failed to serialize error")
		}
	}
)

// Server is a REST server for the nero v1 API.
type Server struct {
	repos  map[string]*repo.Repository
	logger *zap.Logger
}

// NewServer creates a new server with pre-defined repositories.
func NewServer(repos []*repo.Repository, logger *zap.Logger) (*Server, error) {
	reposById := make(map[string]*repo.Repository, len(repos))
	for _, r := range repos {
		repoId := r.ID()
		if _, ok := reposById[repoId]; ok {
			return nil, fmt.Errorf("duplicate repository ID %s", repoId)
		}

		reposById[repoId] = r
	}

	return &Server{
		repos:  reposById,
		logger: logger,
	}, nil
}

// NewRouter creates a new nero v1 API router.
func NewRouter(handler v1.StrictServerInterface) http.Handler {
	h := v1.NewStrictHandlerWithOptions(handler, nil, v1.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  DefaultRequestErrorHandler,
		ResponseErrorHandlerFunc: DefaultResponseErrorHandler,
	})

	return v1.HandlerWithOptions(h, v1.ChiServerOptions{ErrorHandlerFunc: DefaultRequestErrorHandler})
}

// Repos returns all repositories available to the server.
func (s *Server) Repos() []*repo.Repository {
	return maps.Values(s.repos)
}
