package controllers

import (
	"./BaseController"

	"../models/QuestionModel"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

type QuestionController struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/question.
func (c *QuestionController) Get() mvc.Result {
	// c.GetLoggedUserName()

	questionModel := QuestionModel.New()
	questions, _ := questionModel.FindAllEnable()

	var homeView = mvc.View{
		Name: "question/questionList.html",
		Data: context.Map{
			"Title":     "Question",
			"Questions": questions,
		},
	}
	return homeView
}
