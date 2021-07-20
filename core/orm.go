package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const DEADLINE = 1 * time.Second

/* Connect to sql server. Mysql support only */
func connect() (*sql.DB, error) {
	config := GetConfig()

	//dsn := "root:123456@tcp(godockerDB:3306)/?charset=utf8"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.Database.User, config.Database.Pass, config.Database.Host, config.Database.Post)

	db, err := sql.Open(config.Database.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect OK.")
	return db, err
}

/*Create database if not exist. If it is exist, use it.*/
func createDataBase(db *sql.DB) error {
	fmt.Println(" create database...")
	config := GetConfig()
	db_name := config.Database.Name
	var err error
	query_create := "CREATE DATABASE IF NOT EXISTS " + db_name
	query_use := "USE " + db_name
	deadline := time.Now().Add(DEADLINE)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// create db if not exists
	result, err := tx.Exec(query_create)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// use this db
	result, err = tx.Exec(query_use)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// query_table := "CREATE TABLE Persons (		PersonID int,		LastName varchar(255),		FirstName varchar(255),		Address varchar(255),		City varchar(255)	);"
	// if ctx, err := db.Godb.QueryContext(ctx, query_table); err != nil {
	// 	log.Fatalf("unable to connect to database: %v", err)
	// 	fmt.Println(ctx)
	// }

	return err
}

func getType(a_type reflect.Type) string {
	switch a_type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.String:
		return "varchar(255)"
	case reflect.Bool:
		return "Boolean"
	case reflect.Struct:
		switch a_type.String() {
		case "time.Time":
			return "timestamp"
		default:
			return ""
		}
	default:
		return ""
	}
}

func getTagDb(a_tag reflect.StructTag) string {
	return a_tag.Get("db")
}

/*Parse s struct to a CREATE command. */
func parseStructField(a_struct reflect.StructField) (string, error) {
	// verify
	struct_value := reflect.ValueOf(a_struct)
	if struct_value.Kind() != reflect.Struct {
		log.Fatal("unsupport type")
		return "", errors.New("parse: unsupport type")
	}

	//struct_type := a_struct.Type
	struct_name := a_struct.Name
	struct_tag := a_struct.Tag
	fmt.Println(struct_name + string(struct_tag))
	create_str := "CREATE TABLE " + a_struct.Name + " ("
	fields_str := ""
	for i := 0; i < a_struct.Type.NumField(); i++ {
		f := a_struct.Type.Field(i)
		fmt.Println(f.Name)
		fmt.Println(f.Type)
		fmt.Println(f.Tag)
		fields_str = fields_str + f.Name + " " + getType(f.Type) + " " + getTagDb(f.Tag) + ","
		//createATable(db, f)
		// fmt.Printf("%d: %s %s = %v\n", i,
		// 	s.Field(i).Name, f.Type(), f.Interface())
	}
	fields_str = fields_str[:len(fields_str)-1]
	last_str := " " + getTagDb(a_struct.Tag) + ");"
	return create_str + fields_str + last_str, nil
}

/* Create a table if it is not exist.
Return a error if it is exist.*/
func createATable(db *sql.DB, table reflect.StructField) error {
	create_table, err := parseStructField(table)
	fmt.Println(create_table)
	if err != nil {
		return err
	}

	deadline := time.Now().Add(DEADLINE)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	if ctx, err := db.QueryContext(ctx, create_table); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		fmt.Println(ctx)
	}
	return nil
}

/* Migrate all object on Tables struct */
func Migrate(tables interface{}) error {
	fmt.Println("db create ....")
	db, err := connect()
	if err != nil {
		fmt.Println(db)
	}

	//create db
	err = createDataBase(db)
	Check(err)

	// create tables
	fmt.Println(tables)
	v := reflect.ValueOf(tables)
	// switch v.Kind() {
	// case reflect.Bool:
	// 	fmt.Println(v.Bool())
	// case reflect.String:
	// 	fmt.Println(v.String())
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 	fmt.Println(v.Int())
	// case reflect.Struct:
	// 	fmt.Println(v.Type().Name())
	// default:
	// 	fmt.Printf("unhandled kind %s", v.Kind())
	// }
	if v.Kind() != reflect.Struct {
		log.Fatal("unsupport type")
		return errors.New("tables: it is not the tables")
	}
	// if v.Kind() == reflect.Ptr {
	// 	v = v.Elem()
	// }

	// if v.Kind() != reflect.Struct {
	// 	log.Fatal("unexpected type")
	// }
	//v = v.Elem()
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
		fmt.Println(f.Tag)
		createATable(db, f)
		// fmt.Printf("%d: %s %s = %v\n", i,
		// 	s.Field(i).Name, f.Type(), f.Interface())
	}

	return nil
}
