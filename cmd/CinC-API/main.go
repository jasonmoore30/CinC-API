package main

import (
	"fmt"
	"log"

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

	var dev bool
	dev = true
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		dev = false
	}

	//This would be the proper way to check, but we are going to use whatever works
	//env := os.Getenv("ENVIRONMENT")
	var dsn string
	if dev {
		config := &cinc.DBConfig{
			Connection: "tcp(156.143.17.176)",
			DBName:     "jcovington",
			DBUser:     "jcovington",
			DBPass:     "Z48tuaOs",
		}
		dsn = config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
	} else {
		// config := &cinc.DBConfig{
		// 	Connection: "erxv1bzckceve5lh.cbetxkdyhwsb.us-east-1.rds.amazonaws.com",
		// 	DBName:     "ekrazwe0spgirfvb",
		// 	DBUser:     "lg4zljacvp2tkm4x",
		// 	DBPass:     "clh3e6aww7a0600o",
		// }
		//dsn = config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
		dsn = "mysql://lg4zljacvp2tkm4x:clh3e6aww7a0600o@erxv1bzckceve5lh.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306/ekrazwe0spgirfvb"
	}

	//dsn := config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
	log.Println("DSN string being used: " + dsn)
	models.InitDB(dsn)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting...")
		return
	}
	router := cincObj.Gin
	router.Run(":8000")
	fmt.Printf("Now running on port:8000")
}
