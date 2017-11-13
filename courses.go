package cinc

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/jasonmoore30/CinC-API/models"
)

func getCourses(c *gin.Context) {
	courses, err := models.GetCourses()
	checkErr(err)
	c.JSON(200, courses)
}

func getCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	course, err := models.GetCourse(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}
	c.JSON(200, course)
}

func newCourse(c *gin.Context) {
	var myCourse = new(models.Course)

	//if I understand the documentation correctly, this is all you have to do to put json data into struct form.
	err := c.Bind(myCourse)
	checkErr(err)

	err = models.AddCourse(myCourse)
	checkErr(err)
	c.AbortWithStatus(200)
}

func deleteCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	err := models.DeleteCourse(id)
	checkErr(err)
	c.AbortWithStatus(http.StatusOK)
}

func updateCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	var myCourse = new(models.Course)

	err := c.Bind(myCourse)
	checkErr(err)

	err = models.UpdateCourse(myCourse, id)
	checkErr(err)

	c.AbortWithStatus(http.StatusOK)
}
