package models

import (
	"errors"
	"fmt"
)

//Experience ..
type Experience struct {
	ID          int    `db:"expID" json:"id"`
	Type        string `db:"expType" json:"type"`
	Description string `db:"expPost" json:"desc"`
	ImgURL      string `db:"expImg" json:"url"`
}

//GetExperiences ..
func GetExperiences() ([]*Experience, error) {
	stmt, err := db.Prepare("SELECT * FROM furmexp")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	experiences := make([]*Experience, 0)
	for rows.Next() {
		exp := new(Experience)
		err := rows.Scan(&exp.ID, &exp.Type, &exp.Description, &exp.ImgURL)
		if err != nil {
			return nil, err
		}
		experiences = append(experiences, exp)
	}
	fmt.Println("Returning Experiences!")
	return experiences, nil
}

//GetExperience ..
func GetExperience(id string) (*Experience, error) {
	row, err := db.Query("SELECT * FROM furmexp WHERE expID=?")
	if err != nil {
		return nil, err
	}
	var exp = new(Experience)
	for row.Next() {
		err = row.Scan(&exp.ID, &exp.Type, &exp.Description, &exp.ImgURL)
		if err != nil {
			return nil, err
		}
	}
	if exp.ID == 0 || exp.Type == "" {
		err := errors.New("experience with id: " + id + " not found")
		return nil, err
	}
	return exp, nil
}

//AddExperience ..
func AddExperience(myExperience *Experience) error {
	stmt, err := db.Prepare("INSERT INTO furmexp (expType, expPost, expImg) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	ok := myExperience.validate()
	if !ok {
		return nil
	}
	result, err := stmt.Exec(myExperience.Type, myExperience.Description, myExperience.ImgURL)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("New experience created with id: %d", id)
	return nil
}

//DeleteExperience ..
func DeleteExperience(id string) error {
	stmt, err := db.Prepare("DELETE FROM furmexp WHERE expID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateExperience ..
func UpdateExperience(myExperience *Experience, id string) error {

	stmt, err := db.Prepare("UPDATE furmexp SET expType=?, expPost=?, expImg=? WHERE expID=?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(myExperience.Type, myExperience.Description, myExperience.ImgURL, id)
	if err != nil {
		return err
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Experience updated: ", rowNum)

	return nil
}

func (*Experience) validate() bool {
	return true
}
