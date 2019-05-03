package todo

import ("context"
	"errors"
	"github.com/rs/xid"
	"math/rand"
	"sync"
	"time"
)

//Package xid is a globally unique id generator suited for web scale

//Create interface with all functions
type TodoService interface{
	GetAllForUser(ctx context.Context, username string) ([]Todo, error)
	GetByID(ctx context.Context, id string) (Todo, error)
	Add(ctx context.Context, todo Todo) (Todo, error)
	Update(ctx context.Context, id string, todo Todo) error
	Delete(ctx context.Context, id string) error
}

var (
	ErrInconsistentIDs = errors.New("Inconsistent IDs")
	ErrNotFOund = errors.New("Not found")
)

func NewInMemoryTodoService() TodoService{
	s := &inmemService{
		m : map[string]Todo{},
	}
	rand.Seed(time.Now().UnixNano())
	return s
}

type inmemService struct{
	sync.RWMutex
	m map[string]Todo
}

func (s inmemService) GetAllForUser(ctx context.Context, username string) ([]Todo, error) {
	//lock writing data while reading
	s.RLock()

	//after everything is done the function allow writing data
	defer s.RUnlock()

	todos := make([]Todo,0,len(s.m))

	for _, todo := range s.m  {
		if todo.UserName == username{
			todos = append(todos, todo)
		}
	}

	//return all todos and nil for error
	return todos,nil
}

func (s inmemService) GetByID(ctx context.Context, id string) (Todo, error) {
	//Lock for for writing
	s.Lock()
	//Unlock for writing
	defer s.Unlock()

	if todo, ok := s.m[id]; ok{
		return todo,nil
	}

	//return empty object and nil
	return Todo{},nil
}

func (s inmemService) Add(ctx context.Context, todo Todo) (Todo, error) {
	//Lock for for writing
	s.Lock()
	//Unlock for writing
	defer s.Unlock()

	//Get unique ID
	todo.ID = xid.New().String()
	//Get time
	todo.CreatedOn = time.Now()

	s.m[todo.ID] = todo
	return todo,nil
}

func (s inmemService) Update(ctx context.Context, id string, todo Todo) error {
	s.Lock()

	defer s.Unlock()

	if(id != todo.ID){
		return ErrInconsistentIDs
	}

	if _,ok := s.m[id]; ok{
		return ErrNotFOund
	}

	s.m[todo.ID] = todo
	return nil
}

func (s inmemService) Delete(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if _,ok := s.m[id]; ok{
		return ErrNotFOund
	}

	delete(s.m, id)
	return nil;
}
