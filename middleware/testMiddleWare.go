package middleware

import (
	"log"
	"net/http"

	"../core"
)

type TestBefore core.MiddleWare

var ttt TestBefore = TestBefore{Run: func(w http.ResponseWriter, r *http.Request) {
	log.Println("testBefore")
	// ttt.Next.Run(w, r)
}}
