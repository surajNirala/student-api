package storage

import "github.com/surajNirala/student-api/internal/models"

type Storage interface {
	StudentList() ([]models.Student, error)
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (models.Student, error)
	UpdateStudentByID(name string, email string, age int, id int64) (string, error)
	DeleteStudentByID(id int64) (string, error)
}
