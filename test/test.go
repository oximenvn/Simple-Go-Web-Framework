package main

import (
	"fmt"
	"reflect"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"log"
)

type DataBase struct {
	Godb *sql.DB
	Name string
}

func (db DataBase) make(){
	query_create := "CREATE DATABASE IF NOT EXISTS " + db.Name
	//query_use = "USE " + db.Name
	rows, err := db.Godb.Query(query_create)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

}

type BaseModel struct {
	Create int
	Update string
}

type Curder interface {
	find() BaseModel
	insert() BaseModel
	update() BaseModel
	delete() BaseModel
}

func (model BaseModel) find() BaseModel {
	s := reflect.ValueOf(&model).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	return BaseModel{}
}

func main() {
	aa := BaseModel{1, "abc"}
	aa.find()

	dsn := "root:123456@tcp(localhost:3306)/"

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    // defer sql.Close()

    // rows, err := db.Query("select id, first_name from user limit 10")
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer rows.Close()

	godb := DataBase{db, "go_db"}
	godb.make()
}
