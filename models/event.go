package models

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
func GetEvents() ([]Event, error) {
	var events []Event
	rows, err := db.Query("SELECT * FROM Events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event = new(Event)

		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.Location, &event.Start, &event.End)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}
	return events, nil
}
