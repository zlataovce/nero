package server

import (
	"github.com/cephxdev/nero/internal/errors"
	"github.com/cephxdev/nero/repo"
	"github.com/cephxdev/nero/server/nekos/v2"
	"github.com/cephxdev/nero/server/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"net/http"
)

var corsOpts = cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,
	MaxAge:           300,
}

// NewNeroRouter creates a new nero API router.
func NewNeroRouter(repos []*repo.Repository, logger *zap.Logger) (http.Handler, error) {
	srv, err := v1.NewServer(repos, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create nero v1 api handler")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  zap.NewStdLog(logger),
		NoColor: true,
	}))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(corsOpts))
	r.Mount("/api/v1", v1.NewRouter(srv))

	return r, nil
}

// NewNekosRouter creates a new nekos API router.
func NewNekosRouter(repos []*repo.Repository, logger *zap.Logger) (http.Handler, error) {
	srv, err := v2.NewServer(repos, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create nekos v2 api handler")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  zap.NewStdLog(logger),
		NoColor: true,
	}))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(corsOpts))
	r.Mount("/api/v2", v2.NewRouter(srv))

	return r, nil
}
