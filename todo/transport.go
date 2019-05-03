package todo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

var ErrMissingParam = errors.New("Missing parameter")

func MakeHTTPHandler(endpoints TodoEndpoints) http.Handler{

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.DefaultCompress)

	todoRouter := chi.NewRouter()

	todoRouter.Get("/", httptransport.NewServer(
		endpoints.GetAllForUserEndPoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Get("/{id}", httptransport.NewServer(
		endpoints.GetByIDEndpoint,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Post("/", httptransport.NewServer(
		endpoints.AddEndpoint,
		decodeAddRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Put("/{id}", httptransport.NewServer(
		endpoints.UpdateEndpoint,
		decodeUpdateRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	todoRouter.Delete("/{id}", httptransport.NewServer(
		endpoints.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Mount("/todos",todoRouter)

	return r
}

func decodeGetRequest(i context.Context, request2 *http.Request) (request interface{}, err error) {
	return GetAllForUserRequest{},err
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func decodeGetByIDRequest(i context.Context, request2 *http.Request) (request interface{}, err error) {
	id := chi.URLParam(request2,"id")
	if(id == ""){
		return nil, ErrMissingParam
	}
	return GetByIDRequest{id},err
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func decodeAddRequest(i context.Context, request2 *http.Request) (request interface{}, err error) {
	var todo Todo
	err = render.Decode(request2,&todo)
	if(err != nil){
		return nil,err
	}
	return AddRequest{todo},err
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func decodeUpdateRequest(i context.Context, request2 *http.Request) (request interface{}, err error) {
	id := chi.URLParam(request2,"id")
	if(id == ""){
		return nil,ErrMissingParam
	}
	var todo Todo
	err = render.Decode(request2,&todo)
	if(err != nil){
		return nil,err
	}
	return UpdateRequest{id,todo},err

}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func decodeDeleteRequest(i context.Context, request2 *http.Request) (request interface{}, err error) {
	id := chi.URLParam(request2, "id")
	if(id == ""){
		return nil,ErrMissingParam
	}
	return DeleteRequest{"id"},err
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func encodeResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	if err,ok := response.(error); ok && err != nil{
		encodeError(ctx,err,writer)
		return nil
	}
	writer.Header().Set("Content-type","application/json; charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}


func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func codeFrom(err error) int{
	switch err {
	case ErrInconsistentIDs, ErrMissingParam:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}