package cinc

import (
	gin "github.com/gin-gonic/gin"
	// import "database/sql"
	// import _ "github.com/go-sql-driver/mysql"
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
	DBPass     string
}

// NewCinc creates a new Gin router the CinC website, already configured wiht the MySQL connection
func NewCinc(config *DBConfig) (*Cinc, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	//router.Use(rateLimit, gin.Recovery())

	// This is where all of our route endpoints will go.
	// The last parameter of each is a function that is defined in another file of the
	// cinc package. We can break it up into as many independent files as we want.

	//router.GET("/api/calendar/events/", getEvents)

	return &Cinc{
		Environment: "dev",
		Gin:         router,
	}, nil

}
