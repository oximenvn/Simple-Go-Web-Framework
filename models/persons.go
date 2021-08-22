package model

import (
	"time"
)

type Persons struct {
	Id         int `db:"NOT NULL AUTO_INCREMENT PRIMARY KEY"`
	Name       string
	Created_at time.Time `db:"DEFAULT CURRENT_TIMESTAMP"`
	Created_by string
	Updated_at time.Time `db:"DEFAULT now() ON UPDATE now()"`
	Updated_by string
}
