package controller

import (
	//"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"../core"
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

	//We tell Go exactly where we can find our html file. We ask Go to parse the html file (Notice
	// the relative path). We wrap it in a call to template.Must() which handles any errors
	// and halts if there are fatal errors

	templates := template.Must(template.ParseFiles("views/welcome.html"))

	//Our HTML comes with CSS that go needs to provide when we run the app. Here we tell go to create
	// a handle that looks in the static directory, go then uses the "/static/" as a url that our
	//html can refer to when looking for our css and other files.

	//Takes the name from the URL query e.g ?name=Martin, will set welcome.Name = Martin.
	if name := r.FormValue("name"); name != "" {
		welcome.Name = name
	}
	//If errors show an internal server error message
	//I also pass the welcome struct to the welcome-template.html file.
	if err := templates.ExecuteTemplate(w, "welcome.html", welcome); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
