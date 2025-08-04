package mysql

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/surajNirala/student-api/internal/config"
	"github.com/surajNirala/student-api/internal/models"
)

type MySQL struct {
	Db *sql.DB
}

func MysqlConnect(cfg *config.Config) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &MySQL{Db: db}, nil
}

func (m *MySQL) StudentList() ([]models.Student, error) {
	var list []models.Student
	stmt, err := m.Db.Prepare("SELECT id,name,email,age,created_at,updated_at FROM students ORDER BY id DESC")
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

func (m *MySQL) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := m.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")
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

func (m *MySQL) GetStudentByID(id int64) (models.Student, error) {
	var student models.Student
	stmt, err := m.Db.Prepare("SELECT id,name,email,age,created_at,updated_at FROM students WHERE id = ?")
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

func (m *MySQL) DeleteStudentByID(id int64) (string, error) {
	msg := ""
	stmt, err := m.Db.Prepare("DELETE FROM students WHERE id = ?")
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

func (m *MySQL) UpdateStudentByID(name string, email string, age int, id int64) (string, error) {
	stmt, err := m.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
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

func (m *MySQL) StudentFileUpload10MB(fileName string, fileData []byte) (string, error) {
	uploadDir := "uploads"
	// Ensure the upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}
	filePath := filepath.Join(uploadDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(fileData)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("File %s uploaded successfully", fileName), nil
}

func (m *MySQL) StudentLargeFileUpload(fileName string, fileReader io.Reader) (string, error) {
	uploadDir := "uploads"
	// Ensure the upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}
	filePath := filepath.Join(uploadDir, fileName)
	dstFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, fileReader)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Large File %s uploaded successfully", fileName), nil
}
