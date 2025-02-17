package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/arishazmat/go-crud/internal/types"
	"github.com/arishazmat/go-crud/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
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

		slog.Info("Create new todo")

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "true"})
	}
}
