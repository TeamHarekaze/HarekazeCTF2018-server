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
	if c.IsBeforeCompetition() {
		return mvc.View{
			Name: "question/questionNoList.html",
			Data: context.Map{
				"Title": "Question",
			},
		}
	}

	questionModel := QuestionModel.New()
	questions, err := questionModel.List(c.GetLoggedTeamName())
	if err != nil {
		return mvc.Response{Err: err}
	}

	var homeView = mvc.View{
		Name: "question/questionList.html",
		Data: context.Map{
			"Title":     "Question",
			"Questions": questions,
		},
	}
	return homeView
}
