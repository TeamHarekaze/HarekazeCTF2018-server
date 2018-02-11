package controllers

import (
	"github.com/TeamHarekaze/HarekazeCTF2018-server/cachemodels/RankingCache"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/web/controllers/BaseController"
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
		return c.Error(err)
	}

	var homeView = mvc.View{
		Name: "ranking/index.html",
		Data: context.Map{
			"Title":       "Ranking",
			"Ranks":       ranks,
			"IsLoggedIn":  c.IsLoggedIn(),
			"CurrentPage": "ranking",
		},
	}
	return homeView
}
