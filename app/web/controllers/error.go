package controllers

import (
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/mvc"
)

type Error500Controller struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/error/500.
func (c *Error500Controller) Get() mvc.Result {
	return mvc.View{
		Name: "error/500.html",
	}
}

type Error404Controller struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/404.
func (c *Error404Controller) Get() mvc.Result {
	return mvc.View{
		Name: "error/404.html",
	}
}
