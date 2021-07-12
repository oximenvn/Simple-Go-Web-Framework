package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DataBase struct {
	Godb *sql.DB
	Name string
}

const shortDuration = 1 * time.Second

func (db DataBase) make() {
	fmt.Println("run make")

	// create db
	query_create := "CREATE DATABASE IF NOT EXISTS " + db.Name
	query_use := "USE " + db.Name
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	tx, err := db.Godb.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := tx.Exec(query_create)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)
	rows, err = tx.Exec(query_use)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// query := "SELECT *	FROM `help_topic`	LIMIT 50"
	// rows, err = db.Godb.Query(query)
	// if err != nil {
	// 	log.Fatal(err)
	// 	fmt.Println(err)
	// }
	// fmt.Println(rows)

	// ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
	// defer cancel()

	query_table := "CREATE TABLE Persons (		PersonID int,		LastName varchar(255),		FirstName varchar(255),		Address varchar(255),		City varchar(255)	);"
	if ctx, err := db.Godb.QueryContext(ctx, query_table); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		fmt.Println(ctx)
	}
}

func (db DataBase) createTable(model interface{}) {
	fmt.Println("db create ....")
	fmt.Println(model)
	v := reflect.ValueOf(model)
	switch v.Kind() {
	case reflect.String:
		fmt.Println(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Println(v.Int())
	default:
		fmt.Printf("unhandled kind %s", v.Kind())
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}

	s := v.Type()
	fmt.Println("parse ....")
	fmt.Println(v)
	fmt.Println(s)
	nameOfT := s.Name()
	fmt.Println(nameOfT)

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Println(f.Name)
		fmt.Println(f.Type)
		// fmt.Printf("%d: %s %s = %v\n", i,
		// 	s.Field(i).Name, f.Type(), f.Interface())
	}

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
	make()
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

type Persons struct {
	BaseModel
	PersonID  int
	LastName  string
	FirstName string
	Address   string
	City      string
}

func (model Persons) find() {
	s := reflect.ValueOf(&model).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func main() {
	aa := BaseModel{1, "abc"}
	aa.find()
	bb := Persons{}
	//bb.find()

	dsn := "root:123456@tcp(godockerDB:3306)/?charset=utf8"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	godb := DataBase{db, "go_db"}
	//godb.make()

	godb.createTable(bb)
}
