package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/redismodels/RankingCache"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/mvc"
)

// AdminFeed override BaseController
type AdminFeed struct {
	BaseController.Base
}

// JSONFeed is struct for JSON feed
type JSONFeed struct {
	Standings []RankingCache.Rank `json:"standings"`
}

// Get handles GET: http://localhost:8080/<APP_ADMIN_HASH>/feed.
// Display JSON feed for CTFtime.org
func (c *AdminFeed) Get() mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/%s/feed", os.Getenv("APP_ADMIN_HASH")))
		return mvc.Response{Path: "/user/login"}
	}

	rankingCache := RankingCache.New()
	defer rankingCache.Close()
	ranks, err := rankingCache.Rank()
	if err != nil {
		return c.Error(err)
	}

	var feed = JSONFeed{}
	feed.Standings = ranks

	output, err := json.Marshal(feed)
	if err != nil {
		return c.Error(err)
	}

	return mvc.Response{
		ContentType: "application/json",
		Text:        string(output),
	}
}
