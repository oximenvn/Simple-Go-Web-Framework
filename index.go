package main

import (
	"fmt"
	"net/http"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func main() {
	http.HandleFunc("/", welcome)
	http.ListenAndServe(":8080", nil)
}
