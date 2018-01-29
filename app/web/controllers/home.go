package controllers

import (
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/mvc"
)

type HomeController struct {
	BaseController.Base
}

func (c *HomeController) Get() mvc.Result {
	var homeView = mvc.View{
		Name: "home/home.html",
		Data: map[string]interface{}{
			"Title":       "Hello Page",
			"MyMessage":   "Welcome to my awesome website",
			"CurrentPage": "home",
		},
	}
	return homeView
}
