package main

import (
	//"fmt"
	"net/http"

	controller "github.com/oximenvn/Simple-Go-Web-Framework/controllers"
	core "github.com/oximenvn/Simple-Go-Web-Framework/core"
	middleware "github.com/oximenvn/Simple-Go-Web-Framework/middlewares"
)

func init() {
	core.Routes.AddRoute("/", "get", http.HandlerFunc(controller.Test.Action))

	commonHandlers := core.Middleware(middleware.LoggingHandler)

	core.Routes.AddRoute("/id", "get", http.HandlerFunc(controller.Test.Get123))
	core.Routes.AddRoute("/abc/{stt}/xyz/{ert}/dfg", "get", http.HandlerFunc(controller.Test.Get123))

	authenHandlers := commonHandlers.Append(middleware.AuthHandler)
	core.Routes.AddRoute("/id/{id}", "get", authenHandlers.ThenFunc(controller.Test.Asd))

	// core.Routes.AddMiddleWare("/id", core.MiddleWare(middleware.TestBefore))
}
