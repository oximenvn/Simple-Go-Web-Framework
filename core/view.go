package core

import (
	"html/template"
	"net/http"
	"strings"
)

func ServeView(w http.ResponseWriter, path string, value interface{}) {

	//We tell Go exactly where we can find our html file. We ask Go to parse the html file (Notice
	// the relative path). We wrap it in a call to template.Must() which handles any errors
	// and halts if there are fatal errors
	templates := template.Must(template.ParseFiles(path))
	//If errors show an internal server error message
	//I also pass the value struct to the file.
	paths := strings.Split(path, "/")
	fileName := paths[len(paths)-1]
	if err := templates.ExecuteTemplate(w, fileName, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
