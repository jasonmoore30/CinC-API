package cinc

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/jasonmoore30/CinC-API/models"
)

func getPosts(c *gin.Context) {
	posts, err := models.GetPosts()
	checkErr(err)
	c.JSON(200, posts)
}

func getPost(c *gin.Context) {
	id := c.Params.ByName("id")
	post, err := models.GetPost(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}
	c.JSON(200, post)
}

func newPost(c *gin.Context) {
	var myPost = new(models.Post)
	err := c.Bind(myPost)
	checkErr(err)

	err = models.AddPost(myPost)
	checkErr(err)
	c.AbortWithStatus(200)
}

func deletePost(c *gin.Context) {
	id := c.Params.ByName("id")
	err := models.DeletePost(id)
	checkErr(err)
	c.AbortWithStatus(http.StatusOK)

}

func updatePost(c *gin.Context) {
	id := c.Params.ByName("id")
	var myPost = new(models.Post)

	err := c.Bind(myPost)
	checkErr(err)

	err = models.UpdatePost(myPost, id)
	checkErr(err)

	c.AbortWithStatus(http.StatusOK)
}
