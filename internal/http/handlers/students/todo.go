package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/arishazmat/go-crud/internal/storage"
	"github.com/arishazmat/go-crud/internal/types"
	"github.com/arishazmat/go-crud/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var todo types.Todo

		err := json.NewDecoder(r.Body).Decode(&todo)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(todo); err != nil {
			validateErrors := err.(validator.ValidationErrors) // type casting
			response.WriteJson(w, http.StatusBadRequest, response.ValidatonError(validateErrors))
			return
		}

		_, err = storage.CreateTodo(todo.Title, todo.Description)
		slog.Info("Todo created successfully", slog.String("Id", fmt.Sprint(err)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": 0})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("GetById called", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		todo, err := storage.GetTodoById(intId)

		if err != nil {
			slog.Error("Failed to get todo", slog.String("id ", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, todo)
	}
}
