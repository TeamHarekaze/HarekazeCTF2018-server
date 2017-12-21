package controllers

import (
	"errors"
	"os"
	"regexp"

	"../models/QuestionModel"
	"./BaseController"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

// AdminQuestionAdd override BaseController
type AdminQuestionAdd struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/<APP_ADMIN_HASH>/question/add.
func (c *AdminQuestionAdd) Get() mvc.Result {
	return mvc.View{
		Name: "admin/questionAddForm.html",
		Data: context.Map{
			"Title":     "Question Add",
			"AdminHash": os.Getenv("APP_ADMIN_HASH"),
		},
	}
}

// Post handles POST: http://localhost:8080/<APP_ADMIN_HASH>/question/add.
func (c *AdminQuestionAdd) Post() mvc.Result {
	var (
		name     = c.Ctx.FormValue("name")
		flag     = c.Ctx.FormValue("flag")
		score    = c.Ctx.FormValue("score")
		sentence = c.Ctx.FormValue("sentence")
	)
	if !regexp.MustCompile(`^[0-9]+$`).MatchString(score) {
		err := errors.New("Score is number!!")
		return mvc.Response{Err: err, Code: 500}
	}
	questionModel := QuestionModel.New()
	err := questionModel.Save(name, flag, score, sentence)
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH"),
	}
}
