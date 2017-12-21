package controllers

import (
	"os"

	"../models/QuestionModel"
	"./BaseController"
	"github.com/kataras/iris/mvc"
)

// AdminQuestionEdit override BaseController
type AdminQuestionDelete struct {
	BaseController.Base
}

// Any handles GET: http://localhost:8080/<APP_ADMIN_HASH>/question/delete/<question id>.
func (c *AdminQuestionDelete) AnyBy(questionId int) mvc.Result {
	questionModel := QuestionModel.New()
	err := questionModel.Delete(questionId)
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH"),
	}
}
