package cinc

import (
	"crypto/rsa"

	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //This is the driver for MySQL and is vital for DB Connectivity
)

//Cinc is a struct to pass the Gin router to the main package and other important
// info to the main package. I will add more props to this if we need them.
type Cinc struct {
	Environment string
	Gin         *gin.Engine
	privateKey  *rsa.PrivateKey
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

	auth := router.Group("/", AuthToken)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
	//Events CRUD: Ready for Test
	router.GET("/api/calendar/events", getEvents)
	router.GET("/api/calendar/events/:id", getEvent)
	router.POST("/api/calendar/events/new", newEvent)
	// router.DELETE("/api/calendar/events/delete/:id", deleteEvent)
	// router.PUT("/api/calendar/events/update/:id", updateEvent)
	auth.DELETE("/api/calendar/events/delete/:id", deleteEvent)
	auth.PUT("/api/calendar/events/update/:id", updateEvent)

	//Courses CRUD
	router.GET("/api/courses", getCourses)
	router.GET("/api/courses/:id", getCourse)
	router.POST("/api/courses/new", newCourse)
	// router.DELETE("/api/courses/delete/:id", deleteCourse)
	// router.PUT("/api/courses/update/:id", updateCourse)
	auth.DELETE("/api/courses/delete/:id", deleteCourse)
	auth.PUT("/api/courses/update/:id", updateCourse)

	//Experiences CRUD
	router.GET("/api/experiences", getExperiences)
	router.GET("/api/experiences/:id", getExperience)
	router.POST("/api/experiences/new", newExperience)
	// router.DELETE("/api/experiences/delete/:id", deleteExperience)
	// router.PUT("/api/experiences/update/:id", updateExperience)
	auth.DELETE("/api/experiences/delete/:id", deleteExperience)
	auth.PUT("/api/experiences/update/:id", updateExperience)

	//Blog CRUD
	router.GET("/api/blog/posts", getPosts)
	router.GET("/api/blog/posts/:id", getPost)
	router.POST("/api/blog/posts/new", newPost)
	router.DELETE("/api/blog/posts/delete/:id", deletePost)
	router.PUT("/api/blog/posts/update/:id", updatePost)
	auth.DELETE("/api/blog/posts/delete/:id", deletePost)
	auth.PUT("/api/blog/posts/update/:id", updatePost)

	//Authentication

	router.GET("/api/public_key", publicKey)
	router.POST("/api/authenticate", authUser)

	//only a current admin can add another administrator
	auth.POST("/api/register", registerUser)

	//Here is where we attach the private key to our object
	//we will read  the private key from S3 storage or generate one and write it to there

	//Other places that need to deal with this file:
	//Func GenerateUserToken
	//Func GetPublicPem

	//IF we dont have an RSA, make one

	return &Cinc{
		Environment: "dev",
		Gin:         router,
		//PrivateKey:  key,
	}, nil
}
