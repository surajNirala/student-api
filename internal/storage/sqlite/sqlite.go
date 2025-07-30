package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/surajNirala/student-api/internal/config"
	"github.com/surajNirala/student-api/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	createQuery := `CREATE TABLE IF NOT EXISTS students (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT,
						email TEXT,
						age INT,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRECT_TIMESTAMP
					)`
	_, err = db.Exec(createQuery)
	if err != nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) StudentList() ([]models.Student, error) {
	var list []models.Student
	stmt, err := s.Db.Prepare("SELECT id,name,email,age,created_at,updated_at FROM students ORDER BY id DESC")
	if err != nil {
		return list, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age, &student.CreatedAt, &student.UpdatedAt)
		if err != nil {
			return nil, err
		}

		list = append(list, student)
	}
	return list, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return lastID, nil
}

func (s *Sqlite) GetStudentByID(id int64) (models.Student, error) {
	var student models.Student
	stmt, err := s.Db.Prepare("SELECT id,name,email,age,created_at,updated_at FROM students WHERE id = ?")
	if err != nil {
		return student, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&student.Id, &student.Name, &student.Email, &student.Age, &student.CreatedAt, &student.UpdatedAt)
	if err != nil {
		return student, err
	}
	return student, nil
}

func (s *Sqlite) DeleteStudentByID(id int64) (string, error) {
	msg := ""
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return msg, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return msg, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return msg, err
	}
	if rowsAffected == 0 {
		return msg, fmt.Errorf("no student found with id %d", id)
	}
	return fmt.Sprintf("student with id %d deleted successfully", id), nil
}

func (s *Sqlite) UpdateStudentByID(name string, email string, age int, id int64) (string, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("no student found with id %d", id)
	}

	return fmt.Sprintf("student with id %d updated successfully", id), nil
}
