package controllers

import (
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/redismodels/RankingCache"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

type RankingController struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/ranking.
func (c *RankingController) Get() mvc.Result {

	rankingCache := RankingCache.New()
	defer rankingCache.Close()
	ranks, err := rankingCache.Rank()
	if err != nil {
		return mvc.Response{Err: err}
	}

	var homeView = mvc.View{
		Name: "ranking/index.html",
		Data: context.Map{
			"Title":      "Ranking",
			"Ranks":      ranks,
			"IsLoggedIn": c.IsLoggedIn(),
		},
	}
	return homeView
}
