package core

type DataBase struct {
}

type BaseModel struct {
}

type Curder interface {
	find() BaseModel
	insert() BaseModel
	update() BaseModel
	delete() BaseModel
}

// func(model BaseModel) find() BaseModel{

// }
