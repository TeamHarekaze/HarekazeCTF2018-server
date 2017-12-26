package controllers

import (
	"os"

	"../models/TeamModel"
	"./BaseController"
	"github.com/kataras/iris/mvc"
)

// AdminTeamEnable override BaseController
type AdminTeamEnable struct {
	BaseController.Base
}

// GetBy handles GET: http://localhost:8080/<APP_ADMIN_HASH>/team/enable/<team id>.
func (c *AdminTeamEnable) GetBy(teamId int) mvc.Result {
	teamModel := TeamModel.New()
	err := teamModel.Enable(teamId)
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH") + "/team",
	}
}
