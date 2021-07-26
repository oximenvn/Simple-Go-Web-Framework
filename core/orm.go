package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type stype struct {
	Type  string
	Value interface{}
}

const DEADLINE = 1 * time.Second

const INT = "int"
const STRING = "varchar(255)"
const BOOLEAN = "Boolean"
const TIME = "time.Time"
const TIMESTAMP = "timestamp"

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

	return err
}

/*Get sql data type for each golang type*/
func getType(a_type reflect.Type) string {
	switch a_type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return INT
	case reflect.String:
		return STRING
	case reflect.Bool:
		return BOOLEAN
	case reflect.Struct:
		switch a_type.String() {
		case TIME:
			return TIMESTAMP
		default:
			return ""
		}
	default:
		return ""
	}
}

/*get db tag on Tag struct*/
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
	create_str := "CREATE TABLE IF NOT EXISTS " + a_struct.Name + " ("
	fields_str := ""
	for i := 0; i < a_struct.Type.NumField(); i++ {
		f := a_struct.Type.Field(i)
		fmt.Println(f.Name)
		fmt.Println(f.Type)
		fmt.Println(f.Tag)
		fields_str = fields_str + fmt.Sprintf("%s %s %s,", f.Name, getType(f.Type), getTagDb(f.Tag))
	}
	fields_str = fields_str[:len(fields_str)-1]
	last_str := ");"
	if len(getTagDb(a_struct.Tag)) > 0 {
		last_str = ", " + getTagDb(a_struct.Tag) + ");"
	}
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

	config := GetConfig()
	db_name := config.Database.Name
	query_use := "USE " + db_name
	deadline := time.Now().Add(DEADLINE)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// use this db
	result, err := tx.Exec(query_use)
	if err != nil {
		log.Fatal(err)
		log.Fatal(result)
	}
	// create table
	result, err = tx.Exec(create_table)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
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
	if v.Kind() != reflect.Struct {
		log.Fatal("unsupport type")
		return errors.New("tables: it is not the tables")
	}

	s := v.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Println(f.Name)
		fmt.Println(f.Type)
		fmt.Println(f.Tag)
		createATable(db, f)
	}

	defer db.Close()
	return nil
}

func Insert(record interface{}) error {
	fmt.Println("save record ....")
	db, err := connect()
	if err != nil {
		fmt.Println(db)
	}
	defer db.Close()

	name, fields, err := parse(record)
	//fmt.Println(name, fields, err)

	query := "INSERT INTO %s ( %s ) VALUES ( %s);"

	list_name := make([]string, 0, len(fields))
	list_values := make([]string, 0, len(fields))
	for k, v := range fields {
		list_name = append(list_name, k)
		list_values = append(list_values, toString(v))
	}

	query = fmt.Sprintf(query, name, strings.Join(list_name, ","), strings.Join(list_values, ","))
	fmt.Println(query)
	// execute query
	config := GetConfig()
	db_name := config.Database.Name
	query_use := "USE " + db_name
	deadline := time.Now().Add(DEADLINE)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// use this db
	result, err := tx.Exec(query_use)
	if err != nil {
		log.Fatal(err)
	}
	// insert record
	result, err = tx.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func toString(s stype) string {
	switch s.Type {
	case INT:
		return strconv.FormatInt(s.Value.(int64), 10)
	case STRING:
		return fmt.Sprintf("'%s'", s.Value.(string))
	case BOOLEAN:
		return strconv.FormatBool(s.Value.(bool))
	case TIMESTAMP:
		return fmt.Sprintf("'%s'", s.Value.(time.Time).Format("2006-01-02 15:04:05"))
	default:
		return ""
	}
}

func getValue(f reflect.Value) interface{} {
	f_value := f.Type()
	switch f_value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int()
	case reflect.String:
		return f.String()
	case reflect.Bool:
		return f.Bool()
	case reflect.Struct:
		switch f_value.String() {
		case "time.Time":
			a_time := f.Interface().(time.Time)
			return a_time // .Unix() //time.Unix(f.Int(), 0)
		default:
			return ""
		}
	default:
		return ""
	}
}

func parse(obj interface{}) (string, map[string]stype, error) {

	result := make(map[string]stype)
	// verify
	struct_value := reflect.ValueOf(obj)
	if struct_value.Kind() != reflect.Struct {
		log.Fatal("unsupport type")
		return "", result, errors.New("parse: unsupport type")
	}

	struct_type := struct_value.Type()
	struct_name := struct_type.Name()
	//fmt.Println(struct_name)
	for i := 0; i < struct_value.NumField(); i++ {
		f := struct_value.Field(i)
		v := struct_type.Field(i)
		// fmt.Println(v.Name)
		// fmt.Println(f.Type().Name())
		// fmt.Println(f.Type())
		// fmt.Println(f)
		result[v.Name] = stype{getType(f.Type()), getValue(f)}

	}

	return struct_name, result, nil
}

func Finds(record interface{}) error {
	fmt.Println("find records ....")
	db, err := connect()
	if err != nil {
		fmt.Println(db)
	}
	defer db.Close()

	name, fields, err := parse(record)
	//fmt.Println(name, fields, err)

	query := "SELECT * FROM %s %s;"
	list_values := make([]string, 0, len(fields))
	for k, v := range fields {
		// skip if is default value
		f_type := reflect.TypeOf(v.Value)
		v_default := reflect.New(f_type).Elem().Interface()
		if v_default == v.Value {
			continue
		}
		// list_name = append(list_name, k)
		condition := k + " = " + toString(v)
		list_values = append(list_values, condition)
	}
	where_str := ""
	if len(list_values) > 0 {
		where_str = " WHERE " + strings.Join(list_values, ",")
	}
	query = fmt.Sprintf(query, name, where_str)
	fmt.Println(query)
	// execute query
	config := GetConfig()
	db_name := config.Database.Name
	query_use := "USE " + db_name

	_, err = db.Exec(query_use)
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	// r_type := reflect.TypeOf(record)
	// emp := reflect.New(r_type).Elem().Interface()
	// res := reflect.MakeSlice(reflect.SliceOf(r_type), 0, 0)

	list_fields := make([]interface{}, len(fields))
	list_fields2 := make([]interface{}, len(fields))
	aaa := make([]string, len(fields))
	i := 0
	for _, v := range fields {
		list_fields[i] = &aaa[i]
		list_fields2[i] = (v.Value)
		i += 1
	}

	for result.Next() {
		// var Id string
		// var Name, Created_by, Updated_by string
		// var Created_at, Update_at string
		// var aaa [5]string
		err = result.Scan(list_fields...)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(list_fields)
	}
	return nil
}
