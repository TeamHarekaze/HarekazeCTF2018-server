package controllers

import (
	"../../redisClient/RankingCache"
	"./BaseController"

	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

type RankingController struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/ranking.
func (c *RankingController) Get() mvc.Result {

	rankingCache := RankingCache.New()
	ranks, err := rankingCache.Rank()
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
