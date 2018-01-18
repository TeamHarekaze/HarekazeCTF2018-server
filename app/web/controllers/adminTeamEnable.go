package controllers

import (
	"fmt"
	"os"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/TeamModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/mvc"
)

// AdminTeamEnable override BaseController
type AdminTeamEnable struct {
	BaseController.Base
}

// GetBy handles GET: http://localhost:8080/<APP_ADMIN_HASH>/team/enable/<team id>.
func (c *AdminTeamEnable) GetBy(teamId int) mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/%s/team/enable/%d", os.Getenv("APP_ADMIN_HASH"), teamId))
		return mvc.Response{Path: "/user/login"}
	}

	teamModel := TeamModel.New()
	err := teamModel.Enable(teamId)
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH") + "/team",
	}
}
