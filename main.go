package servicetodoapp

import (
	"net/http"
	"servicetodoapp/todo"
)

func main(){
	service := todo.NewInMemoryTodoService()

	endpoints := todo.MakeTodoEnpoints(service)

	err := http.ListenAndServe(":8000", todo.MakeHTTPHandler(endpoints))
	if err != nil {
		panic(err)
	}
}