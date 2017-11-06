package models

import "fmt"

//Event ..
type Event struct {
	ID          int64  `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Date        string `db:"date" json:"date"`
	Location    string `db:"location" json:"location"`
	Start       string `db:"start_time" json:"start_time"`
	End         string `db:"end_time" json:"end_time"`
}

//GetEvents returns an array of all the event objects in the database
func GetEvents() ([]*Event, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("SELECT * FROM Events")
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
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Location, &event.Start, &event.End)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	fmt.Println("Returning events!")
	return events, nil
}

//GetEvent returns the event with the given unique identifier
func GetEvent(id string) (*Event, error) {
	row, err := db.Query("SELECT * FROM Events WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	var event = new(Event)
	for row.Next() {
		err = row.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Location, &event.Start, &event.End)
		if err != nil {
			return nil, err
		}
	}
	return event, nil
}
