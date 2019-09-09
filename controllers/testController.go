package controller

import (
	//"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../core"
)

type testController core.Controller

var Test testController

func (Test testController) Action(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
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
