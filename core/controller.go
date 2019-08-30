package core

import (
	"fmt"
	"net/http"
)

type Controller struct {
}

type ElderlyGent interface {
	SayHi()
	Sing(song string)
	SpendSalary(amount float32)
}

type Action func(w http.ResponseWriter, r *http.Request)

var Test Controller

func (c Controller) Action(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func (c Controller) Asd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "asd!")
}

func (c Controller) Get123(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "123!")
}
