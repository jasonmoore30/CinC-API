package cinc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonmoore30/CinC-API/models"
)

// Event ..
// type Event struct {
// 	ID          int64  `db:"id" json:"id"`
// 	Title       string `db:"title" json:"title"`
// 	Description string `db:"description" json:"description"`
// 	Date        string `db:"date" json:"date"`
// 	Location    string `db:"location" json:"location"`
// 	Start       string `db:"start_time" json:"start_time"`
// 	End         string `db:"end_time" json:"end_time"`
// }

// getEvents is our handler func to write a nice, accurate response or error message
func getEvents(c *gin.Context) {
	events, err := models.GetEvents()
	checkErr(err)
	c.JSON(200, events)
}

func getEvent(c *gin.Context) {
	id := c.Params.ByName("id")
	event, err := models.GetEvent(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}
	c.JSON(200, event)
}

//for now, events go directly into the events table
// TODO: direct post methods to "holding tables" where they can await approval before being sent
// into the actual main data tables
func newEvent(c *gin.Context) {
	var myEvent = new(models.Event)

	//if I understand the documentation correctly, this is all you have to do to put json data into struct form.
	err := c.Bind(myEvent)
	checkErr(err)

	err = models.AddEvent(myEvent)
	checkErr(err)
}

func deleteEvent(c *gin.Context) {
	id := c.Params.ByName("id")
	err := models.DeleteEvent(id)
	checkErr(err)
	c.AbortWithStatus(http.StatusOK)
}

func updateEvent(c *gin.Context) {
	id := c.Params.ByName("id")
	var myEvent = new(models.Event)

	err := c.Bind(myEvent)
	checkErr(err)

	err = models.UpdateEvent(myEvent, id)
	checkErr(err)

	c.AbortWithStatus(http.StatusOK)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
