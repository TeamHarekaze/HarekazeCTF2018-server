package controllers

import (
	"fmt"
	"os"

	"github.com/TeamHarekaze/HarekazeCTF2018-server/web/controllers/BaseController"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

// AdminController override BaseController
type Admin struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/<APP_ADMIN_HASH>.
func (c *Admin) Get() mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/%s", os.Getenv("APP_ADMIN_HASH")))
		return mvc.Response{Path: "/user/login"}
	}
	return mvc.View{
		Name: "admin/index.html",
		Data: context.Map{
			"Title":      "Admin Page",
			"AdminHash":  os.Getenv("APP_ADMIN_HASH"),
			"IsLoggedIn": c.IsLoggedIn(),
		},
	}
}
