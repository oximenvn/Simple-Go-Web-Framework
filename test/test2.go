package main

import (
	"database/sql"
	"log"

	"../models"
	_ "github.com/go-sql-driver/mysql"
)

func (i models.Tables) Say() {

}
func main() {

	dsn := "root:123456@tcp(godockerDB:3306)/?charset=utf8"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//godb.make()

	godb.createTable(bb)
}
