package model

type Cars struct {
	Id        int `db:"NOT NULL AUTO_INCREMENT PRIMARY KEY, INDEX cars_idx (Id)"`
	Own_id    int `db:", FOREIGN KEY (Own_id) REFERENCES Persons(Id) ON DELETE CASCADE ON UPDATE CASCADE"`
	UseDiesel bool
}
