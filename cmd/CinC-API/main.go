package main

import (
	"fmt"
	"log"

	"os"

	cinc "github.com/jasonmoore30/CinC-API"
	"github.com/jasonmoore30/CinC-API/models"
	"github.com/joho/godotenv"
)

func main() {
	cincObj, err := cinc.NewCinc()
	if err != nil {
		log.Println(err)
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENVIRONMENT")
	var dsn string
	if env == "dev" {
		config := &cinc.DBConfig{
			Connection: "tcp(156.143.17.176)",
			DBName:     "jcovington",
			DBUser:     "jcovington",
			DBPass:     "Z48tuaOs",
		}
		dsn = config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
	} else {
		config := &cinc.DBConfig{
			Connection: "erxv1bzckceve5lh.cbetxkdyhwsb.us-east-1.rds.amazonaws.com",
			DBName:     "ekrazwe0spgirfvb",
			DBUser:     "lg4zljacvp2tkm4x",
			DBPass:     "clh3e6aww7a0600o",
		}
		dsn = config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
	}

	//dsn := config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
	log.Println("DSN string being used: " + dsn)
	models.InitDB(dsn)
	if err != nil {
		fmt.Println(err)
	}
	router := cincObj.Gin
	router.Run(":8000")
	fmt.Printf("Now running on port:8000")
}
