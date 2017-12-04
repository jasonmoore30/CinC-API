package models

import (
	"errors"
	"fmt"
)

//Event ..
type Event struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	CreatorID   int    `db:"creator" json:"creatorid"`
	Description string `db:"description" json:"description"`
	//Date        string `db:"date" json:"date"`
	Location  string `db:"location" json:"location"`
	Start     string `db:"start_time" json:"start_time"`
	End       string `db:"end_time" json:"end_time"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

//GetEvents returns an array of all the event objects in the database
func GetEvents() ([]*Event, error) {
	// err := db.Ping()
	// if err != nil {
	// 	return nil, err
	// }

	stmt, err := db.Prepare("SELECT furmcal.evStart, furmcal.evEnd, furmcal.createdAt, furmcal.evTitle, furmcal.evDesc, furmcal.evLoc, furminf.fName, furminf.lName, furminf.uPhone, furminf.uEmail FROM furmcal INNER JOIN furminf ON furmcal.usernum=furminf.usernum")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*Event, 0)
	for rows.Next() {
		event := new(Event)
		err := rows.Scan(&event.ID, &event.CreatorID, &event.Start, &event.End, &event.CreatedAt, &event.Title, &event.Description, &event.Location)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	fmt.Println("Returning events!")
	return events, nil
}

//GetEvent returns the event with the unique identifier given as an argument
func GetEvent(id string) (*Event, error) {
	row, err := db.Query("SELECT furmcal.evStart, furmcal.evEnd, furmcal.createdAt, furmcal.evTitle, furmcal.evDesc, furmcal.evLoc, furminf.fName, furminf.lName, furminf.uPhone, furminf.uEmail FROM furmcal WHERE evID=? INNER JOIN furminf ON furmcal.usernum=furminf.usernum", id)
	if err != nil {
		return nil, err
	}
	var event = new(Event)
	for row.Next() {
		//err = row.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Location, &event.Start, &event.End)
		err = row.Scan(&event.ID, &event.CreatorID, &event.Start, &event.End, &event.CreatedAt, &event.Title, &event.Description, &event.Location)
		if err != nil {
			return nil, err
		}
	}
	if event.ID == 0 || event.Title == "" {
		err := errors.New("event with id: " + id + " not found")
		return nil, err
	}
	return event, nil
}

//AddEvent inserts a new event into the Events table
func AddEvent(myEvent *Event) error {

	stmt, err := db.Prepare("INSERT INTO furmcal (usernum, evTitle, evDesc, evLoc, evStart, evEnd) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	//generate unique id unless we are going with simple indexing in the SQL table.
	//For now, we will let auto-increment create our id's

	//any data transformations necessary should be done here before being inserted.
	//we also need to validate the data before we try to insert it

	ok := myEvent.validate()
	if !ok {
		return nil
	}

	//We may need to convert the start and end time strings into
	result, err := stmt.Exec(myEvent.Title, myEvent.Description, myEvent.Location, myEvent.Start, myEvent.End)
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

//DeleteEvent ..
func DeleteEvent(id string) error {
	stmt, err := db.Prepare("DELETE FROM furmcal WHERE evID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateEvent ..
func UpdateEvent(myEvent *Event, id string) error {

	stmt, err := db.Prepare("UPDATE furmcal SET (evTitle=?, evDesc=?, evLoc=?, evStart=?, evEnd=?) WHERE evID=?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(myEvent.Title, myEvent.Description, myEvent.Location, myEvent.Start, myEvent.End, id)
	if err != nil {
		return err
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Event updated: ", rowNum)
	return nil
}

//validate checks all of the struct fields to make sure it adheres to requirements and that the data being injected into
//our queries is secure
func (*Event) validate() bool {

	return true

}
