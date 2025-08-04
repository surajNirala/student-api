package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/surajNirala/student-api/internal/models"
	"github.com/surajNirala/student-api/internal/storage"
	"github.com/surajNirala/student-api/internal/utils/response"
)

func List(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := storage.StudentList()
		// fmt.Println("list", list)
		// slog.Info("err : ", err.Error())
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}
		response.WriteJson(w, http.StatusOK, list)
	}
}

func Create(storage storage.Storage) http.HandlerFunc {
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
		lastID, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}

		data := make(map[string]any)
		data["Success"] = "OK"
		data["Code"] = 201
		data["id"] = lastID
		response.WriteJson(w, http.StatusCreated, data)
	}
}

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}
		student, err := storage.GetStudentByID(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func UpdateByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var studentupdate models.Student
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&studentupdate)
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
		if err := validator.New().Struct(studentupdate); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validatorErr))
			return
		}
		_, err = storage.GetStudentByID(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		message, err := storage.UpdateStudentByID(studentupdate.Name, studentupdate.Email, studentupdate.Age, id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}

		data := make(map[string]any)
		data["Success"] = "OK"
		data["Code"] = 200
		data["id"] = id
		data["message"] = message
		response.WriteJson(w, http.StatusOK, data)
	}
}

func DeleteByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}
		student, err := storage.DeleteStudentByID(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func FileUpload10MB(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenerateError(err))
			return
		}
		defer file.Close()

		if header.Size > 10*1024*1024 { // 10 MB limit
			fmt.Println("Bad Request: File size exceeds limit", header.Size)
			response.WriteJson(w, http.StatusBadRequest, response.GenerateError(fmt.Errorf("file size exceeds limit")))
			return
		}
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		result, err := storage.StudentFileUpload10MB(header.Filename, fileBytes)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		data := make(map[string]any)
		data["Success"] = result
		data["Code"] = 200
		response.WriteJson(w, http.StatusOK, data)
	}
}

func LargeFileUpload(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenerateError(err))
			return
		}
		defer file.Close()

		result, err := storage.StudentLargeFileUpload(header.Filename, file)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		data := make(map[string]any)
		data["Success"] = result
		data["Code"] = 200
		response.WriteJson(w, http.StatusOK, data)
	}
}
