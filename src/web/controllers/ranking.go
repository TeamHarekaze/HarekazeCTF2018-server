package controllers

import (
	"./BaseController"

	"../models/AnswerModel"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

type RankingController struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/ranking.
func (c *RankingController) Get() mvc.Result {

	answerModel := AnswerModel.New()
	ranks, err := answerModel.Ranking()
	if err != nil {
		return mvc.Response{Err: err}
	}

	var homeView = mvc.View{
		Name: "ranking/index.html",
		Data: context.Map{
			"Title": "Ranking",
			"Ranks": ranks,
		},
	}
	return homeView
}
