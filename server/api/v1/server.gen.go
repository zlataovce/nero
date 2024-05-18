// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /repos/{repo})
	PostRepo(w http.ResponseWriter, r *http.Request, repo string, params PostRepoParams)

	// (DELETE /repos/{repo}/{id})
	DeleteRepoId(w http.ResponseWriter, r *http.Request, repo string, id openapi_types.UUID, params DeleteRepoIdParams)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (POST /repos/{repo})
func (_ Unimplemented) PostRepo(w http.ResponseWriter, r *http.Request, repo string, params PostRepoParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (DELETE /repos/{repo}/{id})
func (_ Unimplemented) DeleteRepoId(w http.ResponseWriter, r *http.Request, repo string, id openapi_types.UUID, params DeleteRepoIdParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostRepo operation middleware
func (siw *ServerInterfaceWrapper) PostRepo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "repo" -------------
	var repo string

	err = runtime.BindStyledParameterWithOptions("simple", "repo", chi.URLParam(r, "repo"), &repo, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "repo", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PostRepoParams

	headers := r.Header

	// ------------- Optional header parameter "X-Nero-Key" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Nero-Key")]; found {
		var XNeroKey string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-Nero-Key", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "X-Nero-Key", valueList[0], &XNeroKey, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-Nero-Key", Err: err})
			return
		}

		params.XNeroKey = &XNeroKey

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostRepo(w, r, repo, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteRepoId operation middleware
func (siw *ServerInterfaceWrapper) DeleteRepoId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "repo" -------------
	var repo string

	err = runtime.BindStyledParameterWithOptions("simple", "repo", chi.URLParam(r, "repo"), &repo, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "repo", Err: err})
		return
	}

	// ------------- Path parameter "id" -------------
	var id openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteRepoIdParams

	headers := r.Header

	// ------------- Optional header parameter "X-Nero-Key" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Nero-Key")]; found {
		var XNeroKey string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-Nero-Key", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "X-Nero-Key", valueList[0], &XNeroKey, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-Nero-Key", Err: err})
			return
		}

		params.XNeroKey = &XNeroKey

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteRepoId(w, r, repo, id, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/repos/{repo}", wrapper.PostRepo)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/repos/{repo}/{id}", wrapper.DeleteRepoId)
	})

	return r
}

type PostRepoRequestObject struct {
	Repo   string `json:"repo"`
	Params PostRepoParams
	Body   *PostRepoJSONRequestBody
}

type PostRepoResponseObject interface {
	VisitPostRepoResponse(w http.ResponseWriter, r *http.Request) error
}

type PostRepo200JSONResponse Media

func (response PostRepo200JSONResponse) VisitPostRepoResponse(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostRepo400JSONResponse Error

func (response PostRepo400JSONResponse) VisitPostRepoResponse(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostRepo401JSONResponse Error

func (response PostRepo401JSONResponse) VisitPostRepoResponse(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type DeleteRepoIdRequestObject struct {
	Repo   string             `json:"repo"`
	Id     openapi_types.UUID `json:"id"`
	Params DeleteRepoIdParams
}

type DeleteRepoIdResponseObject interface {
	VisitDeleteRepoIdResponse(w http.ResponseWriter, r *http.Request) error
}

type DeleteRepoId200JSONResponse Media

func (response DeleteRepoId200JSONResponse) VisitDeleteRepoIdResponse(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type DeleteRepoId400JSONResponse Error

func (response DeleteRepoId400JSONResponse) VisitDeleteRepoIdResponse(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type DeleteRepoId401JSONResponse Error

func (response DeleteRepoId401JSONResponse) VisitDeleteRepoIdResponse(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /repos/{repo})
	PostRepo(ctx context.Context, request PostRepoRequestObject) (PostRepoResponseObject, error)

	// (DELETE /repos/{repo}/{id})
	DeleteRepoId(ctx context.Context, request DeleteRepoIdRequestObject) (DeleteRepoIdResponseObject, error)
}
type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// PostRepo operation middleware
func (sh *strictHandler) PostRepo(w http.ResponseWriter, r *http.Request, repo string, params PostRepoParams) {
	var request PostRepoRequestObject

	request.Repo = repo
	request.Params = params

	var body PostRepoJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostRepo(ctx, request.(PostRepoRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostRepo")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostRepoResponseObject); ok {
		if err := validResponse.VisitPostRepoResponse(w, r); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// DeleteRepoId operation middleware
func (sh *strictHandler) DeleteRepoId(w http.ResponseWriter, r *http.Request, repo string, id openapi_types.UUID, params DeleteRepoIdParams) {
	var request DeleteRepoIdRequestObject

	request.Repo = repo
	request.Id = id
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteRepoId(ctx, request.(DeleteRepoIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteRepoId")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(DeleteRepoIdResponseObject); ok {
		if err := validResponse.VisitDeleteRepoIdResponse(w, r); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
