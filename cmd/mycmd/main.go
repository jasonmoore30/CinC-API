package main

import (
	"fmt"

	cinc "github.com/jasonmoore30/CinC-API"
	"github.com/jasonmoore30/CinC-API/models"
)

func main() {
	cincObj, err := cinc.NewCinc()
	config := &cinc.DBConfig{
		Connection: "tcp(156.143.17.176)",
		DBName:     "jcovington",
		DBUser:     "jcovington",
		DBPass:     "Z48tuaOs",
	}
	dsn := config.DBUser + ":" + config.DBPass + "@" + config.Connection + "/" + config.DBName
	models.InitDB(dsn)
	if err != nil {
		fmt.Println(err)
	}
	router := cincObj.Gin
	router.Run(":8000")
	fmt.Printf("Now running on port:8000")
}
