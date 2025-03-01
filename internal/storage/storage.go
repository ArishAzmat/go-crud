package storage

import "github.com/arishazmat/go-crud/internal/types"

type Storage interface {
	CreateTodo(title string, description string) (int64, error)
	GetTodoById(id int64) (types.Todo, error)
	GetTodos() ([]types.Todo, error)
}
