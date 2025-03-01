package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/arishazmat/go-crud/internal/config"
	"github.com/arishazmat/go-crud/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title text NOT NULL,
		description text NOT NULL
		)`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil

}

func (s *Sqlite) CreateTodo(title string, description string) (int64, error) {

	stmt, err := s.Db.Prepare(`INSERT INTO todos (title, description) VALUES (?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(title, description)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil

}

func (s *Sqlite) GetTodoById(id int64) (types.Todo, error) {
	stmt, err := s.Db.Prepare("SELECT id, title, description from todos WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Todo{}, err
	}

	defer stmt.Close()

	var todo types.Todo

	err = stmt.QueryRow(id).Scan(&todo.Id, &todo.Title, &todo.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Todo{}, fmt.Errorf("No Todo found with id %s", fmt.Sprint(id))
		}
		return types.Todo{}, fmt.Errorf("Error while fetching todos", err)
	}

	return todo, nil
}

func (s *Sqlite) GetTodos() ([]types.Todo, error) {
	stmt, err := s.Db.Prepare("SELECT id, title, description from todos")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []types.Todo

	for rows.Next() {
		var todo types.Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Description)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}
