package controllers

import (
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/redismodels/SolveCache"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
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
				"Title":       "Question",
				"IsLoggedIn":  c.IsLoggedIn(),
				"CurrentPage": "questions",
			},
		}
	}

	solveCache := SolveCache.New()
	defer solveCache.Close()
	questions, err := solveCache.List(c.GetLoggedTeamName())

	// questionModel := QuestionModel.New()
	// questions, err := questionModel.List(c.GetLoggedTeamName())
	if err != nil {
		return c.Error(err)
	}

	var homeView = mvc.View{
		Name: "question/questionList.html",
		Data: context.Map{
			"Title":       "Question",
			"Questions":   questions,
			"IsLoggedIn":  c.IsLoggedIn(),
			"CurrentPage": "questions",
		},
	}
	return homeView
}
