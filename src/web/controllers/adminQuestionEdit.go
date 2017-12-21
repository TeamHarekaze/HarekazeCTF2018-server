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

// AdminQuestionEdit override BaseController
type AdminQuestionEdit struct {
	BaseController.Base
}

// GetBy handles GET: http://localhost:8080/<APP_ADMIN_HASH>/question/edit/<question id>.
func (c *AdminQuestionEdit) GetBy(questionId int) mvc.Result {
	questionModel := QuestionModel.New()
	question, err := questionModel.FindId(questionId)
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.View{
		Name: "admin/questionEditForm.html",
		Data: context.Map{
			"Title":     "Question Edit",
			"AdminHash": os.Getenv("APP_ADMIN_HASH"),
			"Question":  question,
		},
	}
}

// PostBy handles GET: http://localhost:8080/<APP_ADMIN_HASH>/question/edit/<question id>.
func (c *AdminQuestionEdit) PostBy(questionId int) mvc.Result {
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
	err := questionModel.Update(questionId, name, flag, score, sentence)
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH"),
	}
}
