package main

import (
	//"fmt"
	"net/http"

	"./controllers"
	"./core"
	"./middleware"
)

func init() {
	core.Routes.AddRoute("/", "get", http.HandlerFunc(controller.Test.Action))
	core.Routes.AddRoute("/id", "get", http.HandlerFunc(controller.Test.Get123))
	core.Routes.AddRoute("/abc/{stt}/xyz/{ert}/dfg", "get", http.HandlerFunc(controller.Test.Get123))

	commonHandlers := core.Middleware(middleware.LoggingHandler)
	core.Routes.AddRoute("/id/{id}", "get", commonHandlers.ThenFunc(controller.Test.Asd))

	// core.Routes.AddMiddleWare("/id", core.MiddleWare(middleware.TestBefore))
}
