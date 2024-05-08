package main

import (
	"encoding/base64"
	"github.com/cephxdev/nero/internal/errors"
	"github.com/cephxdev/nero/server/api"
	"github.com/cephxdev/nero/server/api/v1"
	"github.com/urfave/cli/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
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

func (ac *appContext) handleUpload(cCtx *cli.Context, m *v1.ProtoMedia_Meta) (err error) {
	c, err := v1.NewClientWithResponses(cCtx.String("url"))
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	f, err := os.Open(filepath.Clean(cCtx.String("path")))
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer func() {
		if err0 := f.Close(); err0 != nil {
			err = multierr.Append(err, errors.Wrap(err0, "failed to close file"))
		}
	}()

	b, err := io.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	res, err := c.PostRepoWithResponse(cCtx.Context, cCtx.String("repo"), v1.ProtoMedia{
		Data: base64.StdEncoding.EncodeToString(b),
		Meta: m,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	ac.logger.Info("request completed", zap.ByteString("body", res.Body))
	return err
}
