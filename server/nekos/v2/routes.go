package v2

import (
	"context"
	"encoding/json"
	"github.com/cephxdev/nero/internal/errors"
	"github.com/cephxdev/nero/repo/media"
	"github.com/cephxdev/nero/repo/media/meta"
	"github.com/cephxdev/nero/server/api"
	"github.com/cephxdev/nero/server/api/nekos/v2"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type category struct {
	Format string `json:"format"`
}

func (s *Server) GetCategories(_ context.Context, _ v2.GetCategoriesRequestObject) (v2.GetCategoriesResponseObject, error) {
	res := make(v2.GetCategories200JSONResponse)
	for i := range s.repos {
		res[i] = category{Format: "gif"} // TODO: make an educated guess about the content
	}

	return res, nil
}

func (s *Server) Search(_ context.Context, _ v2.SearchRequestObject) (v2.SearchResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetCategoryFiles(_ context.Context, request v2.GetCategoryFilesRequestObject) (v2.GetCategoryFilesResponseObject, error) {
	r, ok := s.repos[request.Category]
	if !ok {
		return v2.GetCategoryFiles404JSONResponse(v2.Error{Code: http.StatusNotFound, Message: "category not found"}), nil
	}

	num := 1
	if request.Params.Amount != nil {
		num = *request.Params.Amount
	}
	if num > 20 { // clamp amount
		num = 20
	}

	return &filesRes{items: r.Random(num)}, nil
}

func (s *Server) GetCategoryFile(_ context.Context, request v2.GetCategoryFileRequestObject) (v2.GetCategoryFileResponseObject, error) {
	r, ok := s.repos[request.Category]
	if !ok {
		return v2.GetCategoryFile404JSONResponse(v2.Error{Code: http.StatusNotFound, Message: "category not found"}), nil
	}

	id, err := uuid.Parse(request.Filename)
	if err != nil {
		return v2.GetCategoryFile404JSONResponse(v2.Error{Code: http.StatusNotFound, Message: "file not found"}), nil
	}

	m := r.Get(id)
	if m == nil {
		return v2.GetCategoryFile404JSONResponse(v2.Error{Code: http.StatusNotFound, Message: "file not found"}), nil
	}

	return &fileRes{path: m.Path}, nil
}

type fileRes struct {
	path string
}

func (fr *fileRes) writeResponse(w http.ResponseWriter, r *http.Request) (err error) {
	f, err := os.Open(fr.path)
	if err != nil {
		return errors.Wrap(err, "failed to open media")
	}
	defer func() {
		if err0 := f.Close(); err0 != nil {
			err = multierr.Append(err, errors.Wrap(err0, "failed to close file"))
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to stat media")
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(fr.path))

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	return err
}

func (fr *fileRes) VisitGetCategoryFileResponse(w http.ResponseWriter, r *http.Request) error {
	return fr.writeResponse(w, r)
}

type filesRes struct {
	items []*media.Media
}

func (fr *filesRes) VisitGetCategoryFilesResponse(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	u := &(*r.URL)
	u.Fragment = ""
	u.RawQuery = ""

	res := make([]v2.Result, len(fr.items))
	for i, m0 := range fr.items {
		res[i] = wrapResult(u, m0)
	}

	return json.NewEncoder(w).Encode(v2.GetCategoryFiles200JSONResponse{Results: res})
}

func wrapResult(base *url.URL, m *media.Media) v2.Result {
	res := v2.Result{Url: base.JoinPath(filepath.Base(m.Path)).String()}

	switch data := m.Meta.(type) {
	case *meta.GenericMetadata:
		res.ArtistHref = api.MakeOptString(data.ArtistLink)
		res.ArtistName = api.MakeOptString(data.Artist)
		res.SourceUrl = api.MakeOptString(data.Source)
	case *meta.AnimeMetadata:
		res.AnimeName = api.MakeOptString(data.Name)
	}

	return res
}
