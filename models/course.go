package models

import (
	"errors"
	"fmt"
)

//Course ..
type Course struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Department      string `json:"dept"`
	Department2     string `json:"dept2"`
	Description     string `json:"descrip"`
	Faculty         string `json:"faculty"`
	Faculty2        string `json:"faculty2"`
	CincDescription string `json:"cincComp"`
	// Approved        int    `json:"approved"`
}

//GetCourses ..
func GetCourses(admin bool) ([]*Course, error) {

	query := ""
	if admin {
		query = "SELECT * FROM furmcourse WHERE adApproval=0"
	} else {
		query = "SELECT * FROM furmcourse WHERE adApproval <>0"
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blank int
	courses := make([]*Course, 0)
	for rows.Next() {
		course := new(Course)
		err := rows.Scan(&course.ID, &course.Title, &course.Department, &course.Department2, &course.Description, &course.Faculty, &course.Faculty2, &course.CincDescription, &blank)
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
	row, err := db.Query("SELECT * FROM furmcourse WHERE cID=?", id)
	if err != nil {
		return nil, err
	}
	var blank int
	var course = new(Course)
	for row.Next() {
		err = row.Scan(&course.ID, &course.Title, &course.Department, &course.Department2, &course.Description, &course.Faculty, &course.Faculty2, &course.CincDescription, &blank)
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

	stmt, err := db.Prepare("INSERT INTO furmcourse (title, dept, dept2, description, faculty, faculty2, cinc_description, adApproval) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	ok := myCourse.validate()
	if !ok {
		//return custom error message about validity of submitted struct
	}
	result, err := stmt.Exec(myCourse.Title, myCourse.Department, myCourse.Department2, myCourse.Description, myCourse.Faculty, myCourse.Faculty2, myCourse.CincDescription, 0)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("New course created with id: %d", id)
	return nil

}

//DeleteCourse ..
func DeleteCourse(id string) error {
	stmt, err := db.Prepare("DELETE FROM furmcourse WHERE cID=?")
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
	stmt, err := db.Prepare("UPDATE furmcourse SET title=?, dept=?, dept2=?, description=?, faculty=?, faculty2=?, cinc_description=? WHERE cID=?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(myCourse.Title, myCourse.Department, myCourse.Department2, myCourse.Description, myCourse.Faculty, myCourse.Faculty2, myCourse.CincDescription, id)
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

//ApproveCourse sets the adApproval to true for the specified course
func ApproveCourse(id string) error {
	_, err := db.Query("UPDATE furmcourse SET adApproval=? WHERE cID=?", 1, id)
	if err != nil {
		return err
	}
	return nil
}

//validation function for a Course struct

func (course *Course) validate() bool {

	return true
}
