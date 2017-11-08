package cinc

import (
	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //This is the driver for MySQL and is vital for DB Connectivity
)

//Cinc is a struct to pass the Gin router to the main package and other important
// info to the main package. I will add more props to this if we need them.
type Cinc struct {
	Environment string
	Gin         *gin.Engine
}

// DBConfig will but database connection configuration right here
//URL, username, password, etc for Joel's MySQL server
type DBConfig struct {
	Connection string
	DBName     string
	DBUser     string
	DBPass     string
}

// Cors ..
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// NewCinc creates a new Gin router the CinC website, already configured wiht the MySQL connection
func NewCinc() (*Cinc, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(Cors())
	//router.Use(rateLimit, gin.Recovery())

	// This is where all of our route endpoints will go.
	// The last parameter of each is a function that is defined in another file of the
	// cinc package. We can break it up into as many independent files as we want.
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Home",
		})
	})

	//Will add search/query string endpoints for each at a later time

	//Events CRUD: Ready for Test
	router.GET("/api/calendar/events", getEvents)
	router.GET("/api/calendar/events/:id", getEvent)
	router.POST("/api/calendar/events/new", newEvent)
	router.DELETE("/api/calendar/events/delete/:id", deleteEvent)
	router.PUT("/api/calendar/event/update/:id", updateEvent)

	//Courses CRUD

	//Experiences CRUD

	//Blog CRUD

	//Courses CRUD

	return &Cinc{
		Environment: "dev",
		Gin:         router,
	}, nil
}
