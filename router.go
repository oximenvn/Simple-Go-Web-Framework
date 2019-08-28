package main

import (
	"fmt"
	"net/http"
)

func routing() {
	http.HandleFunc("/", welcome)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}
