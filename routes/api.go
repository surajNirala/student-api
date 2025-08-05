package routes

import (
	"net/http"
	"time"

	"github.com/surajNirala/student-api/internal/http/handlers/student"
	"github.com/surajNirala/student-api/internal/storage"
)

func RouteLoad(router *http.ServeMux, storage storage.Storage) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Student API " + time.Now().Format(time.RFC3339)))
	})
	router.HandleFunc("GET /api/students", student.List(storage))
	router.HandleFunc("POST /api/students", student.Create(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateByID(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteByID(storage))

	router.HandleFunc("GET /api/students1", student.List(storage))

	router.HandleFunc("POST /api/students/file-upload", student.FileUpload10MB(storage))
	router.HandleFunc("POST /api/students/large-file-upload", student.LargeFileUpload(storage))
}
