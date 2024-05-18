package v1

import (
	"context"
	"encoding/base64"
	"github.com/cephxdev/nero/internal/errors"
	"github.com/cephxdev/nero/repo"
	"github.com/cephxdev/nero/repo/media"
	"github.com/cephxdev/nero/repo/media/meta"
	"github.com/cephxdev/nero/server/api"
	"github.com/cephxdev/nero/server/api/v1"
	"net/http"
)

var (
	unauthorizedError = &api.HTTPError{
		Err:    errors.New("wrong or missing key"),
		Status: http.StatusUnauthorized,
		Type:   string(v1.Unauthorized),
	}
)

func (s *Server) PostRepo(_ context.Context, request v1.PostRepoRequestObject) (v1.PostRepoResponseObject, error) {
	r, ok := s.repos[request.Repo]
	if !ok {
		return v1.PostRepo400JSONResponse(v1.Error{Type: v1.NotFound, Description: "unknown repository"}), nil
	}

	if !checkKey(r, api.MakeString(request.Params.XNeroKey)) {
		return nil, unauthorizedError
	}

	var m meta.Metadata
	if request.Body.Meta != nil {
		m0, err := request.Body.Meta.ValueByDiscriminator()
		if err != nil {
			return nil, err
		}

		m = unwrapMetadata(m0)
	}

	d, err := base64.StdEncoding.DecodeString(request.Body.Data)
	if err != nil {
		return v1.PostRepo400JSONResponse(v1.Error{Type: v1.BadRequest, Description: "failed to decode data"}), nil
	}

	m0, err := r.Create(d, m)
	if err != nil {
		return nil, err
	}

	m1, err := wrapMedia(m0)
	if err != nil {
		return nil, err
	}

	return v1.PostRepo200JSONResponse(m1), nil
}

func (s *Server) DeleteRepoId(_ context.Context, request v1.DeleteRepoIdRequestObject) (v1.DeleteRepoIdResponseObject, error) {
	r, ok := s.repos[request.Repo]
	if !ok {
		return v1.DeleteRepoId400JSONResponse(v1.Error{Type: v1.NotFound, Description: "unknown repository"}), nil
	}

	if !checkKey(r, api.MakeString(request.Params.XNeroKey)) {
		return nil, unauthorizedError
	}

	m := r.Get(request.Id)
	if m == nil {
		return v1.DeleteRepoId400JSONResponse(v1.Error{Type: v1.NotFound, Description: "unknown item id"}), nil
	}

	if err := r.Remove(request.Id); err != nil {
		return nil, err
	}

	m0, err := wrapMedia(m)
	if err != nil {
		return nil, err
	}

	return v1.DeleteRepoId200JSONResponse(m0), nil
}

func wrapMedia(m *media.Media) (v1.Media, error) {
	var (
		m0  = &v1.Media_Meta{}
		err error
	)
	switch v := wrapMetadata(m.Meta).(type) {
	case v1.GenericMetadata:
		err = m0.FromGenericMetadata(v)
	case v1.AnimeMetadata:
		err = m0.FromAnimeMetadata(v)
	}

	if err != nil {
		return v1.Media{}, err
	}

	return v1.Media{
		Format: wrapFormat(m.Format),
		Id:     m.ID,
		Meta:   m0,
	}, nil
}

func wrapFormat(f media.Format) v1.MediaFormat {
	switch f {
	case media.FormatImage:
		return v1.Image
	case media.FormatAnimatedImage:
		return v1.AnimatedImage
	default:
		return v1.Unknown
	}
}

func unwrapMetadata(v interface{}) meta.Metadata {
	switch m := v.(type) {
	case v1.GenericMetadata:
		return &meta.GenericMetadata{
			Source:     api.MakeString(m.Source),
			Artist:     api.MakeString(m.Artist),
			ArtistLink: api.MakeString(m.ArtistLink),
		}
	case v1.AnimeMetadata:
		return &meta.AnimeMetadata{
			Name: api.MakeString(m.Name),
		}
	}

	return nil
}

func wrapMetadata(v meta.Metadata) interface{} {
	switch m := v.(type) {
	case *meta.GenericMetadata:
		return v1.GenericMetadata{
			Source:     api.MakeOptString(m.Source),
			Artist:     api.MakeOptString(m.Artist),
			ArtistLink: api.MakeOptString(m.ArtistLink),
		}
	case *meta.AnimeMetadata:
		return v1.AnimeMetadata{
			Name: api.MakeOptString(m.Name),
		}
	}

	return nil
}

func checkKey(r *repo.Repository, key string) bool {
	if expectedKey, ok := r.Meta().Value(repo.AuthKey); ok {
		return key == expectedKey
	}
	return true // no required key, no authentication needed
}
