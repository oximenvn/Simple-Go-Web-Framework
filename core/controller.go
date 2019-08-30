package core

import (
	//"fmt"
	"net/http"
)

type Controller struct {
}

type Action func(w http.ResponseWriter, r *http.Request)
