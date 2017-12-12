package cinc

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/jasonmoore30/CinC-API/models"
)

func getExperiences(c *gin.Context) {
	experiences, err := models.GetExperiences(false)
	checkErr(err)
	c.JSON(200, experiences)
}
func getExperiencesAdmin(c *gin.Context) {
	experiences, err := models.GetExperiences(true)
	checkErr(err)
	c.JSON(200, experiences)
}

func getExperience(c *gin.Context) {
	id := c.Params.ByName("id")
	experience, err := models.GetExperience(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}
	c.JSON(200, experience)
}

func newExperience(c *gin.Context) {
	var myExperience = new(models.Experience)
	err := c.Bind(myExperience)
	//validate
	checkErr(err)

	err = models.AddExperience(myExperience)
	checkErr(err)
	c.AbortWithStatus(200)
}

func deleteExperience(c *gin.Context) {
	id := c.Params.ByName("id")
	err := models.DeleteExperience(id)
	checkErr(err)
	c.AbortWithStatus(http.StatusOK)
}

func updateExperience(c *gin.Context) {
	id := c.Params.ByName("id")
	var myExperience = new(models.Experience)

	err := c.Bind(myExperience)
	checkErr(err)

	err = models.UpdateExperience(myExperience, id)
	checkErr(err)

	c.AbortWithStatus(http.StatusOK)
}

func approveExperience(c *gin.Context) {
	id := c.Params.ByName("id")
	err := models.ApproveExperience(id)
	checkErr(err)

	c.AbortWithStatus(http.StatusAccepted)
}
