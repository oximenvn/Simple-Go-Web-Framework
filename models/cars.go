package model

type Cars struct {
	Id  int `db:"NOT NULL AUTO_INCREMENT PRIMARY KEY, INDEX cars_idx (Id)"`
	Own int `db:"REFERENCES Persons(Id)"`
}
