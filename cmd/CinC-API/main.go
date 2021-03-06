package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	cinc "github.com/jasonmoore30/CinC-API"
	"github.com/jasonmoore30/CinC-API/models"
)

func main() {
	cincObj, err := cinc.NewCinc()
	if err != nil {
		log.Println(err)
		return
	}

	//var dev bool
	//dev = true
	//err = godotenv.Load()
	_ = godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("Error loading .env file")

	//}

	//This would be the proper way to check, but we are going to use whatever works
	env := os.Getenv("ENVIRONMENT")
	port := os.Getenv("PORT")
	log.Println(port)
	var dsn string
	var dsn2 string
	if env == "dev" {
		config := &cinc.DBConfig{
			Connection: "tcp(156.143.17.176)",
			DBName:     "jcovington",
			DBUser:     "jcovington",
			DBPass:     "Z48tuaOs",
		}
		dsn = config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
	} else {
		dsn = "lg4zljacvp2tkm4x:clh3e6aww7a0600o@tcp(erxv1bzckceve5lh.cbetxkdyhwsb.us-east-1.rds.amazonaws.com)/ekrazwe0spgirfvb"
		//dsn2 = "mysql://lg4zljacvp2tkm4x:clh3e6aww7a0600o@erxv1bzckceve5lh.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306/ekrazwe0spgirfvb"
	}

	log.Println("DSN string being used: " + dsn)

	err = models.InitDB(dsn)
	if err != nil {
		err2 := models.InitDB(dsn2)
		if err2 != nil {
			log.Println("Another pinging error!")
		}
	}

	router := cincObj.Gin
	router.Run(":" + port)
	fmt.Printf("Now running on port: " + port)
}
