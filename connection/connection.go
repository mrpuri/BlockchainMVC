package connection

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" //sql driver
)

//Hello for test and create Database connection

var Db, err = sql.Open("mysql",
	"root:@tcp(127.0.0.1:3306)/Blockchain") //database name is Blockchain, user: root and password: ""

func r() {
	if err != nil {
		log.Fatal(err)
	}
}
