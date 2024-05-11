package main

import (
	"encoding/base64"
	"fmt"
	"github.com/cephxdev/nero/internal/errors"
	"github.com/cephxdev/nero/server/api"
	"github.com/cephxdev/nero/server/api/v1"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// handleUploadGeneric handles the upload generic sub-command.
func (ac *appContext) handleUploadGeneric(cCtx *cli.Context) error {
	m := &v1.ProtoMedia_Meta{}
	_ = m.FromGenericMetadata(v1.GenericMetadata{
		Artist:     api.MakeOptString(cCtx.String("artist")),
		ArtistLink: api.MakeOptString(cCtx.String("artist-link")),
		Source:     api.MakeOptString(cCtx.String("source")),
	})

	return ac.handleUpload(cCtx, m)
}

// handleUploadAnime handles the upload anime sub-command.
func (ac *appContext) handleUploadAnime(cCtx *cli.Context) error {
	m := &v1.ProtoMedia_Meta{}
	_ = m.FromAnimeMetadata(v1.AnimeMetadata{
		Name: api.MakeOptString(cCtx.String("name")),
	})

	return ac.handleUpload(cCtx, m)
}

func (ac *appContext) handleUpload(cCtx *cli.Context, m *v1.ProtoMedia_Meta) error {
	c, err := v1.NewClientWithResponses(cCtx.String("url"))
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	var (
		path = cCtx.String("path")
		data io.ReadCloser
	)
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		ac.logger.Info("treating path as remote url", zap.String("path", path))

		res, err := http.Get(path)
		if err != nil {
			return errors.Wrap(err, "failed to get remote url")
		}

		if res.StatusCode > 399 {
			return fmt.Errorf("remote url request returned error status code %d", res.StatusCode)
		}

		data = res.Body
	} else {
		if data, err = os.Open(filepath.Clean(path)); err != nil {
			return errors.Wrap(err, "failed to open file")
		}
	}

	b, err := io.ReadAll(data)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	if err = data.Close(); err != nil {
		return errors.Wrap(err, "failed to close data stream")
	}

	res, err := c.PostRepoWithResponse(cCtx.Context, cCtx.String("repo"), v1.ProtoMedia{
		Data: base64.StdEncoding.EncodeToString(b),
		Meta: m,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	code := res.StatusCode()
	if code > 399 {
		ac.logger.Error(
			"request completed with errors",
			zap.String("status", res.Status()),
			zap.Int("code", code),
			zap.ByteString("body", res.Body),
		)

		// error out to force an error exit code
		return fmt.Errorf("request completed with error status code %d", code)
	}

	ac.logger.Info("request completed", zap.ByteString("body", res.Body))
	return nil
}
