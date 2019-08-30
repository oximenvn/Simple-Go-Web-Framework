package core

import (
	"net/http"
)

type MiddleWare struct {
	Next *MiddleWare
	Run  func(w http.ResponseWriter, r *http.Request)
	// DoNext func(w http.ResponseWriter, r *http.Request){
	// 	Next.Run(w,r)
	// }
}

// type Actor interface {
// 	Action(w http.ResponseWriter, r *http.Request)
// }

// func (m *MiddleWare) Run(w http.ResponseWriter, r *http.Request) {

// }

// type MiddleAction func(w http.ResponseWriter, r *http.Request, next *MiddleWare)

// type MiddleWareStack struct {
// }
