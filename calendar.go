package cinc

import (
	"database/sql"
	"encoding/json"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//rateLimit.... This is an example endpoint.
func getEvents(c *gin.Context) {

	// This opens the DB connection, using parameters from the cinc DBConfig object
	db, err := sql.Open("mysql", "astaxie:astaxie@/test?charset=utf8")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("SELECT * FROM ? WHERE id = ?")
	checkErr(err)

	res, err := stmt.Exec("ExampleTable", "4")
	checkErr(err)
	resJSON, err := json.Marshal(res)
	checkErr(err)
	c.Abort()
	c.JSON(200, resJSON)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
