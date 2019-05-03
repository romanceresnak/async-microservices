package todo

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)
// Endpoints collects all endpoints which compose the Todo service
type TodoEndpoints struct {
	GetAllForUserEndPoint endpoint.Endpoint
	GetByIDEndpoint       endpoint.Endpoint
	AddEndpoint           endpoint.Endpoint
	UpdateEndpoint        endpoint.Endpoint
	DeleteEndpoint        endpoint.Endpoint
}

func MakeTodoEnpoints(s TodoService) TodoEndpoints{
	return TodoEndpoints{
		GetAllForUserEndPoint: MakeGetAllForUserEndpoint(s),
		GetByIDEndpoint:       MakeGetByIDEndpoint(s),
		AddEndpoint:           MakeAddEndpoint(s),
		UpdateEndpoint:        MakeUpdateEndpoint(s),
		DeleteEndpoint:        MakeDeleteEndpoint(s),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type GetAllForUserRequest struct{

}

type GetAllForuserResponse struct{
	Todos []Todo `json:todos`
}

func MakeGetAllForUserEndpoint(service TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		username := ctx.Value("username").(string)
		todos, err := service.GetAllForUser(ctx,username)
		return GetAllForuserResponse{todos},err
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	Todo Todo `json:todo`
}

func MakeGetByIDEndpoint(service TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetByIDRequest)
		todo, err := service.GetByID(ctx,req.ID)
		return GetByIDResponse{todo},err
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type AddRequest struct{
	Todo Todo
}

type AddResponse struct{
	Toto Todo `json:"todo"`
}

func MakeAddEndpoint(service TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddRequest)
		todo, err := service.Add(ctx,req.Todo)
		return AddResponse{todo},err
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type UpdateRequest struct{
	ID   string
	Todo Todo
}

type UpdateResponse struct{
}

func MakeUpdateEndpoint(service TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateRequest)
		err = service.Update(ctx,req.ID,req.Todo)
		return UpdateResponse{},err
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type DeleteRequest struct{
	ID string
}
type DeleteResponse struct{

}

func MakeDeleteEndpoint(service TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteRequest)
		err = service.Delete(ctx,req.ID)
		return DeleteResponse{},err
	}
}