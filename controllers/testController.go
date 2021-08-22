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
	//Instantiate a Welcome struct object and pass in some random information.
	//We shall get the name of the user as a query parameter from the URL
	welcome := Welcome{"Anonymous", "Now"}
	//Takes the name from the URL query e.g ?name=Martin, will set welcome.Name = Martin.
	name := r.FormValue("name")
	if name == "" {
		core.ServeView(w, "views/welcome.html", welcome)
		return
	}

	temp := model.Persons{
		//Id:         4,
		Name:       name,
		Created_at: time.Now(),
		Created_by: "test",
		Updated_at: time.Now(),
		Updated_by: "test",
	}

	he, err := core.Finds(model.Persons{Name: name})
	core.Check(err)
	they := he.([]model.Persons)
	fmt.Println(len(they))
	if len(they) == 0 {
		core.Insert(temp)
		welcome.Name = temp.Name
		welcome.Time = "Now"
	} else {
		welcome.Name = they[0].Name
		welcome.Time = they[0].Created_at.Format(time.Stamp)
	}

	core.ServeView(w, "views/welcome.html", welcome)
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
