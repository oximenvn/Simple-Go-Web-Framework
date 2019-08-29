package main

import (
	"fmt"
	"net/http"

	"./core"
)

func init() {
	core.Routes.AddRoute("/", "get", core.Test.Action)
	core.Routes.AddRoute("/123", "get", core.Test.Get123)
	core.Routes.AddRoute("/123/asd", "get", core.Test.Asd)
}

func routing(w http.ResponseWriter, r *http.Request) {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/asdf", asdf)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func asdf(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "asdf!")
}
