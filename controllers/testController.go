package controller

import (
	"fmt"
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
}

func (Test testController) Get123(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "123!")
}
