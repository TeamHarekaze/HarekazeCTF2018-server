package controllers

import (
	"fmt"
	"os"

	"github.com/TeamHarekaze/HarekazeCTF2018-server/app/datamodels/QuestionModel"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/app/web/controllers/BaseController"
	"github.com/kataras/iris/mvc"
)

// AdminQuestionEdit override BaseController
type AdminQuestionDelete struct {
	BaseController.Base
}

// Any handles GET: http://localhost:8080/<APP_ADMIN_HASH>/question/delete/<question id>.
func (c *AdminQuestionDelete) AnyBy(questionId int) mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/%s/question/delete/%d", os.Getenv("APP_ADMIN_HASH"), questionId))
		return mvc.Response{Path: "/user/login"}
	}

	questionModel := QuestionModel.New()
	err := questionModel.Delete(questionId)
	if err != nil {
		return c.Error(err)
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH"),
	}
}
