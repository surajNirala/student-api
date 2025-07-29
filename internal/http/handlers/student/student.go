package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/surajNirala/student-api/internal/models"
	"github.com/surajNirala/student-api/internal/utils/response"
)

func Create() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var student models.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			customMsg := fmt.Errorf("empty body")
			response.WriteJson(w, http.StatusBadRequest, response.GenerateError(customMsg))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenerateError(err))
			return
		}
		// Request Validation
		if err := validator.New().Struct(student); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validatorErr))
			return
		}
		data := make(map[string]string)
		data["Success"] = "OK"
		response.WriteJson(w, http.StatusCreated, data)
		slog.Info("Creating a student")
		w.Write([]byte("Welcome the student new style"))
	}
}
