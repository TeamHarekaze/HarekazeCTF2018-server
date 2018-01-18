package controllers

import (
	"fmt"
	"os"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/TeamModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/mvc"
)

// AdminTeamDisable override BaseController
type AdminTeamDisable struct {
	BaseController.Base
}

// AnyBy handles GET: http://localhost:8080/<APP_ADMIN_HASH>/team/disable/<team id>.
func (c *AdminTeamDisable) AnyBy(teamId int) mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/%s/team/disable/%d", os.Getenv("APP_ADMIN_HASH"), teamId))
		return mvc.Response{Path: "/user/login"}
	}

	teamModel := TeamModel.New()
	err := teamModel.Disable(teamId)
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH") + "/team",
	}
}
