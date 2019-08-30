package main

import (
	//"fmt"
	//"net/http"

	"./controller"
	"./core"
	"./middleware"
)

func init() {
	core.Routes.AddRoute("/", "get", controller.Test.Action)
	core.Routes.AddRoute("/id", "get", controller.Test.Get123)
	core.Routes.AddRoute("/id/{id}", "get", controller.Test.Asd)
	core.Routes.AddRoute("/abc/{stt}/xyz/{ert}/dfg", "get", controller.Test.Get123)
	core.Routes.AddMiddleWare("/id", core.MiddleWare(middleware.TestBefore))
}
