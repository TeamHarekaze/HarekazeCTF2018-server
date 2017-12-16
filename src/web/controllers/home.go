package controllers

import (
	"github.com/kataras/iris/mvc"
)

type HomeController struct {
	mvc.C
}

func (c *HomeController) Get() mvc.Result {
	var homeView = mvc.View{
		Name: "home/home.html",
		Data: map[string]interface{}{
			"Title":     "Hello Page",
			"MyMessage": "Welcome to my awesome website",
		},
	}
	return homeView
}
