package cinc

import (
	"fmt"

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
	checkErr(err)
	c.JSON(200, event)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
