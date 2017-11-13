package models

import (
	"errors"
	"fmt"
)

//Course ..
type Course struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Professors  []string `json:"professors"`
	IsOffered   string   `json:"isOffered"`
	Term        string   `json:"term"`
}

//GetCourses ..
func GetCourses() ([]*Course, error) {
	stmt, err := db.Prepare("SELECT * FROM Courses")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]*Course, 0)
	for rows.Next() {
		course := new(Course)
		err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.Professors, &course.IsOffered, &course.Term)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	fmt.Println("Returning courses!")
	return courses, nil
}

//GetCourse ..
func GetCourse(id string) (*Course, error) {
	row, err := db.Query("SELECT * FROM Courses WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	var course = new(Course)
	for row.Next() {
		err = row.Scan(&course.ID, &course.Title, &course.Description, &course.Professors, &course.IsOffered, &course.Term)
		if err != nil {
			return nil, err
		}
	}
	if course.ID == 0 || course.Title == "" {
		err := errors.New("event with id: " + id + " not found")
		return nil, err
	}
	return course, nil
}

//AddCourse ..
func AddCourse(myCourse *Course) error {

	stmt, err := db.Prepare("INSERT INTO Courses (title, description, professors, isOffered, term) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	ok := myCourse.validate()
	if !ok {
		//return custom error message about validity of submitted struct
	}
	result, err := stmt.Exec(myCourse.Title, myCourse.Description, myCourse.Professors, myCourse.IsOffered, myCourse.Term)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("New event created with id: %d", id)
	return nil

}

//DeleteCourse ..
func DeleteCourse(id string) error {
	stmt, err := db.Prepare("DELETE FROM Courses WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateCourse ..
func UpdateCourse(myCourse *Course, id string) error {
	stmt, err := db.Prepare("UPDATE Courses SET (title=?, description=?, professors=?, isOffered=?, term=?)")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(myCourse.Title, myCourse.Description, myCourse.Professors, myCourse.IsOffered, myCourse.Term, id)
	if err != nil {
		return err
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Course updated: ", rowNum)
	return nil
}

//validation function for a course to add
func (course *Course) validate() bool {

	return true
}
