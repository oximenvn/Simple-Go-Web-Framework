package main

import (
	"time"
)

type Tables struct {
	Persons struct {
		Id         int `db:"NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Name       string
		Created_at time.Time `db:"DEFAULT CURRENT_TIMESTAMP"`
		Creatd_by  string    `db:"DEFAULT CURRENT_USER()"`
		Updated_at time.Time `db:"DEFAULT now() ON UPDATE now()"`
		Updated_by string    `db:"DEFAULT CURRENT_USER() ON UPDATE CURRENT_USER()"`
	}

	Cars struct {
		Id  int `db:"NOT NULL AUTO_INCREMENT PRIMARY KEY"`
		Own int `db:"FOREIGN KEY REFERENCES Persons(Id)"`
	} `db:"INDEX cars_idx (Id)"`
}
