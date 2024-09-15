package main

import (
	"context"
	"fmt"
	"github.com/zlataovce/nero/config"
	"github.com/zlataovce/nero/internal/errors"
	"github.com/zlataovce/nero/repo"
	"github.com/zlataovce/nero/server"
	"github.com/urfave/cli/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

type httpServer struct {
	servers []*http.Server
	errChan chan error

	logger *zap.Logger
}

func (hs *httpServer) add(s *http.Server) {
	hs.servers = append(hs.servers, s)
	go func() {
		hs.logger.Info("listening for http requests", zap.String("addr", s.Addr))
		hs.errChan <- s.ListenAndServe()
	}()
}

func (hs *httpServer) shutdown(ctx context.Context) (err error) {
	for _, s := range hs.servers {
		err = multierr.Append(err, s.Shutdown(ctx))
	}

	return err
}

// handleServer handles the server sub-command.
func (ac *appContext) handleServer(cCtx *cli.Context) (err error) {
	cfg, err := config.ParseWithDefaults(cCtx.String("config"))
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	repos0 := make(map[string]*repo.Repository, len(cfg.Repos))
	for repoId, repoConfig := range cfg.Repos {
		if _, ok := repos0[repoId]; ok {
			return fmt.Errorf("duplicate repository ID %s, path %s", repoId, repoConfig.Path)
		}

		r, err := repo.NewFile(repoId, repoConfig.Path, repoConfig.LockPath, repoConfig.Meta, ac.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create repository")
		}

		repos0[repoId] = r

		ac.logger.Info(
			"registered repository",
			zap.String("repo", repoId),
			zap.String("path", repoConfig.Path),
		)
	}
	defer func() {
		for _, r := range repos0 {
			if err0 := r.Close(); err0 != nil {
				err = multierr.Append(err, errors.Wrap(err0, "failed to close repository"))
			}
		}
	}()

	var (
		repos   = maps.Values(repos0)
		httpSrv = &httpServer{
			errChan: make(chan error),
			logger:  ac.logger,
		}
	)
	if cfg.HTTP.Nero.Enabled() {
		handler, err := server.NewNeroRouter(repos, ac.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create nero api router")
		}

		httpSrv.add(&http.Server{Addr: cfg.HTTP.Nero.Host, Handler: handler})
	}
	if cfg.HTTP.Nekos.Enabled() {
		var baseURL *url.URL
		if cfg.HTTP.Nekos.BaseURL != "" {
			if baseURL, err = url.Parse(cfg.HTTP.Nekos.BaseURL); err != nil {
				return errors.Wrap(err, "failed to parse nekos api base url")
			}
		}

		handler, err := server.NewNekosRouter(repos, baseURL, ac.logger)
		if err != nil {
			return errors.Wrap(err, "failed to create nekos api router")
		}

		httpSrv.add(&http.Server{Addr: cfg.HTTP.Nekos.Host, Handler: handler})
	}

	ctx, stop := signal.NotifyContext(cCtx.Context, os.Interrupt)
	defer stop()

	select {
	case <-ctx.Done():
		ac.logger.Info("shutting down gracefully")
		if err = httpSrv.shutdown(ctx); err != nil {
			err = errors.Wrap(err, "failed to shutdown http server")
		}
	case err = <-httpSrv.errChan:
		err = errors.Wrap(err, "http server errored")
	}

	return err
}
