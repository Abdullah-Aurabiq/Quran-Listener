package main

import (
	"database/sql"
	"fmt"
	"time"
)

// Define structures to hold the JSON data
type QuranResponse struct {
	Quran QuranData `json:"quran"`
}

type QuranData struct {
	UthmaniHafs map[string]Verse `json:"quran-uthmani-hafs"`
}

type Verse struct {
	ID    int    `json:"id"`
	Surah int    `json:"surah"`
	Ayah  int    `json:"ayah"`
	Verse string `json:"verse"`
}

type Courses struct {
	Name            string  `json:"name"`
	Detail          string  `json:"detail"`
	Price           float64 `json:"price"`
	DiscountedPrice float64 `json:"discountedprice"`
}

type Student struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	CourseName string `json:"course_name"`
	Password   string `json:"password"`
	Email      string `json:"email"`
}

type Teacher struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Age        int     `json:"age"`
	Gender     string  `json:"gender"`
	CourseName string  `json:"course_name"`
	Password   string  `json:"password"`
	Email      string  `json:"email"`
	StudentsID []uint8 `json:"student_id"`
}
type Class struct {
	ID         int    `json:"id"`
	TeacherID  int    `json:"teacher_id"`
	Start_Time string `json:"start_time"`
	End_Time   string `json:"end_time"`
	CourseName string `json:"course_name"`
}

// uint8 to time.Time

func Uint8ToTime(u []uint8) (time.Time, error) {
	timeString := string(u)
	layout := time.RFC1123Z // Adjust according to your datetime format
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

// GetSuraNames retrieves the names of all the Suras

func GetClasses(db *sql.DB) ([]Class, error) {
	query := "SELECT c.start_time, c.id, c.end_time, c.course_name, t.id as teacher_id FROM classes c JOIN teachers t ON c.teacher_id = t.id"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var classes []Class
	for rows.Next() {
		var c Class
		var startTime []uint8
		var endTime []uint8
		err := rows.Scan(&startTime, &c.ID, &endTime, &c.CourseName, &c.TeacherID)
		if err != nil {
			return nil, err
		}
		c.Start_Time = string(startTime)
		c.End_Time = string(endTime)
		classes = append(classes, c)
	}
	return classes, nil
}

func (class *Class) CreateClass(db *sql.DB) error {
	teacherID := class.TeacherID
	courseName := class.CourseName
	// layout := "2006-01-02 15:04:05"
	startTime := class.Start_Time
	endTime := class.End_Time

	_, err := db.Exec("INSERT INTO classes (teacher_id, course_name, start_time, end_time) VALUES (?,?,?,?)", teacherID, courseName, startTime, endTime)
	if err != nil {
		return err
	}

	return nil

}

// 	teacherEmail := r.FormValue("teacher_email")
// 	courseName := r.FormValue("course_name")
// 	startTime := r.FormValue("start_time")
// 	endTime := r.FormValue("end_time")

// 	// Get teacher ID from email
// 	var teacherID int
// 	err := app.DB.QueryRow("SELECT id FROM teachers WHERE email = ?", teacherEmail).Scan(&teacherID)
// 	if err != nil {
// 		sendError(w, http.StatusInternalServerError, "Failed to get teacher ID")
// 		return
// 	}

// 	// Create new class
// 	_, err = app.DB.Exec("INSERT INTO classes (teacher_id, course_name, start_time, end_time) VALUES (?, ?, ?, ?)", teacherID, courseName, startTime, endTime)
// 	if err != nil {
// 		sendError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	sendResponse(w, http.StatusOK, map[string]string{"message": "Class created!"}, "", "text/html; charset=UTF-8")
// }

// func (app *App) MyClasses(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "session")
// 	email := session.Values["email"].(string)

// 	var teacherID int
// 	err := app.DB.QueryRow("SELECT id FROM teachers WHERE email = ?", email).Scan(&teacherID)
// 	if err != nil {
// 		sendError(w, http.StatusInternalServerError, "Failed to get teacher ID")
// 		return
// 	}

// 	rows, err := app.DB.Query("SELECT id, course_name, start_time, end_time FROM classes WHERE teacher_id = ?", teacherID)
// 	if err != nil {
// 		sendError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	defer rows.Close()

// 	var classes []Class
// 	for rows.Next() {
// 		var c Class
// 		err := rows.Scan(&c.ID, &c.CourseName, &c.Start_Time, &c.End_Time)
// 		if err != nil {
// 			sendError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		classes = append(classes, c)
// 	}

// 	sendResponse(w, http.StatusOK, classes, "templates/my_classes.html", "text/html; charset=UTF-8")
// }

func (st *Student) GetStudent(db *sql.DB) (Student, error) {
	query := fmt.Sprintf("SELECT ID, Name, Age, Gender, IFNULL(CourseName, 'Did not started Yet') AS CourseName, email, password FROM students WHERE email='%v' AND Name='%v' Limit 1", st.Email, st.Name)
	row := db.QueryRow(query)
	err := row.Scan(&st.ID, &st.Name, &st.Age, &st.Gender, &st.CourseName, &st.Email, &st.Password)
	return *st, err
}

func (t *Teacher) GetTeacher(db *sql.DB) (Teacher, error) {
	query := fmt.Sprintf("SELECT ID, Name, Age, Gender, IFNULL(CourseName, 'Did not started Yet') AS CourseName, email, password, StudentsID FROM teachers WHERE email='%v' AND Name='%v' Limit 1", t.Email, t.Name)
	row := db.QueryRow(query)
	err := row.Scan(&t.ID, &t.Name, &t.Age, &t.Gender, &t.CourseName, &t.Email, &t.Password, &t.StudentsID)
	return *t, err
}
func (st *Student) LoginStudent(db *sql.DB) error {
	query := fmt.Sprintf("SELECT ID, Name, Age, Gender, IFNULL(CourseName, 'Did not started Yet') AS CourseName, email, password FROM students WHERE email='%v' and password='%v' LIMIT 1", st.Email, st.Password)
	row := db.QueryRow(query)
	err := row.Scan(&st.ID, &st.Name, &st.Age, &st.Gender, &st.CourseName, &st.Email, &st.Password)
	return err
}
func (t *Teacher) LoginTeacher(db *sql.DB) error {
	query := fmt.Sprintf("SELECT ID, Name, Age, Gender, IFNULL(CourseName, 'Did not started Yet') AS CourseName, email, password, StudentsID FROM teachers WHERE email='%v' and password='%v' LIMIT 1", t.Email, t.Password)
	row := db.QueryRow(query)
	err := row.Scan(&t.ID, &t.Name, &t.Age, &t.Gender, &t.CourseName, &t.Email, &t.Password, &t.StudentsID)
	return err
}

func (st *Student) CreateStudentAccount(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO students (name, age, gender, password, email) VALUES ('%v',%v,'%v','%v','%v')", st.Name, st.Age, st.Gender, st.Password, st.Email)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	st.ID = int(id)
	return nil
}

func GetAllCourses(db *sql.DB) ([]Courses, error) {
	query := "SELECT course_name, course_detail, course_price, course_discounted_price FROM courses"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	courses := []Courses{}
	for rows.Next() {
		var c Courses
		err := rows.Scan(&c.Name, &c.Detail, &c.Price, &c.DiscountedPrice)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func GetAllStudents(db *sql.DB) ([]Student, error) {
	query := "SELECT ID, Name, Age, Gender, IFNULL(CourseName, 'Did not started Yet') AS CourseName, password, email FROM students"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	students := []Student{}
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Gender, &s.CourseName, &s.Password, &s.Email)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

// func (t *Teacher) GetMyStudents(db *sql.DB) (Teacher, error) {
// 	query := fmt.Sprintf("SELECT StudentsID FROM teachers Where email='%v' and password='%v' Limit 1", t.Email, t.Password)
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	teacher := []Teacher{}
// 	for rows.Next() {
// 		var s Teacher
// 		err := rows.Scan(&t.StudentID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		teacher = append(teacher, s)
// 	}
// 	return teacher, nil
// }
