package controller

import (
	//"context"
	"encoding/json"
	"fmt"

	//"html/template"
	"log"
	"net/http"
	"time"

	"github.com/oximenvn/Simple-Go-Web-Framework/core"
	model "github.com/oximenvn/Simple-Go-Web-Framework/models"
)

type testController core.Controller

var Test testController

//Create a struct that holds information to be displayed in our HTML file
type Welcome struct {
	Name string
	Time string
}

func (Test testController) Action(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "Welcome!")

	//Instantiate a Welcome struct object and pass in some random information.
	//We shall get the name of the user as a query parameter from the URL
	welcome := Welcome{"Anonymous", time.Now().Format(time.Stamp)}

	//Takes the name from the URL query e.g ?name=Martin, will set welcome.Name = Martin.
	if name := r.FormValue("name"); name != "" {
		welcome.Name = name
	}

	core.ServeView(w, "views/welcome.html", welcome)
	abc := model.Persons{
		//Id:         4,
		Name:       "Thanh",
		Created_at: time.Now(),
		Created_by: "Thanh1",
		Updated_at: time.Now(),
		Updated_by: "Thanh2",
	}
	//fmt.Println(abc)
	core.Insert(abc)
	he, err := core.Finds(model.Persons{Name: "Thanh"})
	core.Check(err)
	they := he.([]model.Persons)
	fmt.Println(len(they))
	fmt.Println(they)
}

func (Test testController) Asd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "asd!")
	user := r.Context()
	log.Printf("%+v\n", user)
	//get value from Context
	log.Println(user.Value("id").(string))
	json.NewEncoder(w).Encode(user)
}

func (Test testController) Get123(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "123!")
}
