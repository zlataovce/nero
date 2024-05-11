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

func (s *Server) Search(_ context.Context, request v2.SearchRequestObject) (v2.SearchResponseObject, error) {
	if request.Params.Type < 1 || request.Params.Type > 2 {
		return v2.Search400JSONResponse(v2.Error{Code: http.StatusBadRequest, Message: "invalid type"}), nil
	}

	needed := 20
	if request.Params.Amount != nil {
		needed = *request.Params.Amount
	}
	if needed > 20 { // clamp amount
		needed = 20
	}

	var res []*media.Media
	if request.Params.Category != nil {
		r, ok := s.repos[*request.Params.Category]
		if !ok {
			return v2.Search400JSONResponse(v2.Error{Code: http.StatusBadRequest, Message: "invalid category"}), nil
		}

		res = r.Find(request.Params.Query, media.Format(request.Params.Type), needed)
	} else {
		for _, r := range s.repos {
			res0 := r.Find(request.Params.Query, media.Format(request.Params.Type), needed)
			if needed < len(res0) {
				res0 = res0[:needed]
			}

			res = append(res, res0...)

			needed -= len(res0)
			if needed <= 0 {
				break
			}
		}
	}

	return &filesRes{items: res}, nil
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

func (fr *fileRes) VisitGetCategoryFileResponse(w http.ResponseWriter, r *http.Request) error {
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

type filesRes struct {
	items []*media.Media
}

func (fr *filesRes) VisitSearchResponse(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	u := &(*r.URL) // clean URL
	u.Fragment = ""
	u.RawQuery = ""

	return json.NewEncoder(w).Encode(v2.Search200JSONResponse{Results: wrapResults(u, fr.items)})
}

func (fr *filesRes) VisitGetCategoryFilesResponse(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	u := &(*r.URL) // clean URL
	u.Fragment = ""
	u.RawQuery = ""

	return json.NewEncoder(w).Encode(v2.GetCategoryFiles200JSONResponse{Results: wrapResults(u, fr.items)})
}

func wrapResults(base *url.URL, ms []*media.Media) []v2.Result {
	res := make([]v2.Result, len(ms))
	for i, m0 := range ms {
		res[i] = wrapResult(base, m0)
	}

	return res
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
