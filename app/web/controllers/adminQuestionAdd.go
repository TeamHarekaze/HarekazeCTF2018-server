package controllers

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/QuestionModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
)

// AdminQuestionAdd override BaseController
type AdminQuestionAdd struct {
	BaseController.Base
}

// Get handles GET: http://localhost:8080/<APP_ADMIN_HASH>/question/add.
func (c *AdminQuestionAdd) Get() mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/%s/question/add", os.Getenv("APP_ADMIN_HASH")))
		return mvc.Response{Path: "/user/login"}
	}

	return mvc.View{
		Name: "admin/questionAddForm.html",
		Data: context.Map{
			"Title":      "Question Add",
			"AdminHash":  os.Getenv("APP_ADMIN_HASH"),
			"Token":      c.MakeToken(fmt.Sprintf("/%s/question/add", os.Getenv("APP_ADMIN_HASH"))),
			"IsLoggedIn": c.IsLoggedIn(),
		},
	}
}

// Post handles POST: http://localhost:8080/<APP_ADMIN_HASH>/question/add.
func (c *AdminQuestionAdd) Post() mvc.Result {
	if !c.IsLoggedIn() {
		c.SetRedirectPath(fmt.Sprintf("/%s/question/add", os.Getenv("APP_ADMIN_HASH")))
		return mvc.Response{Path: "/user/login"}
	}

	var (
		name               = c.Ctx.FormValue("name")
		flag               = c.Ctx.FormValue("flag")
		score              = c.Ctx.FormValue("score")
		genre              = c.Ctx.FormValue("genre")
		publish_start_time = c.Ctx.FormValue("publish_start_time")
		publish_now        = c.Ctx.FormValue("publish_now")
		sentence           = c.Ctx.FormValue("sentence")
		token              = c.Ctx.FormValue("csrf_token")
	)
	if !c.CheckTaken(token, fmt.Sprintf("/%s/question/add", os.Getenv("APP_ADMIN_HASH"))) {
		err := errors.New("token error!!")
		return mvc.Response{Err: err, Code: 400}
	}
	if publish_now == "off" && !regexp.MustCompile(`/^(\d{4})-(\d{2})-(\d{2})\s(\d{2}):(\d{2}):(\d{2})$/`).MatchString(publish_start_time) {
		err := errors.New("publish_start_time is yyyy-MM-dd HH:mm:ss!!")
		return mvc.Response{Err: err, Code: 500}
	} else if !regexp.MustCompile(`^[0-9]+$`).MatchString(score) {
		err := errors.New("Score is number!!")
		return mvc.Response{Err: err, Code: 500}
	}
	questionModel := QuestionModel.New()
	err := questionModel.Save(map[string]string{
		"name":               name,
		"flag":               flag,
		"score":              score,
		"genre":              genre,
		"publish_now":        publish_now,
		"publish_start_time": publish_start_time,
		"auther_id":          c.GetLoggedUserID(),
		"sentence":           sentence,
	})
	if err != nil {
		return mvc.Response{Err: err, Code: 500}
	}
	return mvc.Response{
		Path: "/" + os.Getenv("APP_ADMIN_HASH"),
	}
}
