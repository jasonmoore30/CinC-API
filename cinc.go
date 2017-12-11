package cinc

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //This is the driver for MySQL and is vital for DB Connectivity
	"github.com/joho/godotenv"
)

//Cinc is a struct to pass the Gin router to the main package and other important
// info to the main package. I will add more props to this if we need them.
type Cinc struct {
	Environment string
	Gin         *gin.Engine
	PrivateKey  *rsa.PrivateKey
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
	// router.DELETE("/api/blog/posts/delete/:id", deletePost)
	// router.PUT("/api/blog/posts/update/:id", updatePost)
	auth.DELETE("/api/blog/posts/delete/:id", deletePost)
	auth.PUT("/api/blog/posts/update/:id", updatePost)

	//Authentication

	router.GET("/api/public_key", publicKey)
	router.POST("/api/authenticate", authUser)
	router.POST("/api/register", registerUser)

	//Here is where we attach the private key to our object
	//we will read  the private key from S3 storage or generate one and write it to there

	//Connect to S3 Bucket using env credentials

	//var myKey *rsa.PrivateKey
	myKey, err := generateRSA()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error with RSA read/gen")
		return nil, err
	}
	if myKey == nil {
		fmt.Println("THE DAMN KEY IS EMPTY")
	}

	//Other places that need to deal with this file:
	//Func GenerateUserToken in models/token.go
	//Func GetPublicPem	in models/token.go

	cinc := &Cinc{
		Environment: "dev",
		Gin:         router,
		PrivateKey:  myKey,
	}

	return cinc, nil
}

func generateRSA() (*rsa.PrivateKey, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Could not load env file")
		return nil, err
	}

	mySession, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Println("Couldn't connect to the AWS platform")
		return nil, err
	}
	svc := s3.New(mySession)

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String("id_rsa"),
	})
	defer result.Body.Close()

	if err != nil {
		//If not exists, GENERATE and STORE
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}

		marshalKey := x509.MarshalPKCS1PrivateKey(key)
		privPem := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: marshalKey,
		}

		blockToWrite := pem.EncodeToMemory(privPem)

		r := bytes.NewReader(blockToWrite)
		svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String("id_rsa"),
			Body:   r,
		})

		return key, nil
	}
	//fmt.Println("making it here 1!")

	var privKey = make([]byte, 1675)

	_, err = result.Body.Read(privKey)
	//fmt.Println("Raw Bytes: ", privKey)
	if err != nil {
		//log.Println("Error reading the bytes from S3 in generateRSA. Read ", n, " bytes before failure")
	}

	pemBlock, rest := pem.Decode(privKey)
	// if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
	if pemBlock == nil {
		log.Fatal("failed to decode PEM block containing private key")
	}
	//fmt.Println("Making it here 2!")

	myKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if myKey == nil {
		myKey, err := x509.ParsePKCS1PrivateKey(rest)
		if err != nil {
			return nil, err
		}
		fmt.Println("Inner key return.")
		return myKey, nil
	}

	//log.Println("Here is my key: ", myKey)
	return myKey, nil
}
