package models

import (
	"errors"
	"fmt"
)

var (
	//ErrNoUserFound ..
	ErrNoUserFound = errors.New("No user found")
)

//TODO: modify this struct to match Joel's user table

//User ..
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//TODO: rewrite this function

//AddUser ..
func AddUser(user *User) error {
	// _, err := cba.bucket.Insert(getUserKey(newID), user, 0)
	stmt, err := db.Prepare("INSERT INTO furminf (uEmail, pass) VALUES (?, ?)")
	if err != nil {
	}
	result, err := stmt.Exec(user.Email, user.Password)
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

//TODO: rewrite this method

//FindUser ..
func FindUser(userEmail string) (*User, error) {
	row, err := db.Query("SELECT * FROM furminf WHERE uEmail=?", userEmail)
	if err != nil {
		return nil, err
	}
	var user = new(User)
	var blank string
	for row.Next() {
		//err = row.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Location, &event.Start, &event.End)
		err = row.Scan(&blank, &user.Password, &blank, &blank, &user.Email, &blank, &blank, &user.ID)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 || user.Email == "" {
		err := errors.New("user not found")
		return nil, err
	}
	return user, nil
}
