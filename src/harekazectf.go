package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"

	"./md2html"
	"./web/controllers"
)

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	Env_load()
	app := iris.New()
	app.StaticWeb("/asset", "./web/public") //debug only
	view := iris.HTML("./web/views", ".html")
	view.Layout("layouts/layout.html")
	view.AddFunc("md2html", md2html.Md2Html)
	app.RegisterView(view)
	// make session manager
	sessionManager := sessions.New(sessions.Config{
		Cookie:  "HarekazeCTF-session",
		Expires: 30 * time.Minute,
	})

	//route
	mvc.New(app.Party("/")).Register(sessionManager.Start).Handle(&controllers.HomeController{})
	mvc.New(app.Party("/user")).Register(sessionManager.Start).Handle(&controllers.UserController{})
	mvc.New(app.Party("/question")).Register(sessionManager.Start).Handle(&controllers.QuestionController{})
	mvc.New(app.Party("/answer")).Register(sessionManager.Start).Handle(&controllers.AnswerController{})
	mvc.New(app.Party("/ranking")).Register(sessionManager.Start).Handle(&controllers.RankingController{})
	//admin
	mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH"))).Register(sessionManager.Start).Handle(&controllers.Admin{})
	mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionList{})
	mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question/add")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionAdd{})
	mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question/edit")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionEdit{})
	mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question/delete")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionDelete{})
	mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/team/enable")).Register(sessionManager.Start).Handle(&controllers.AdminTeamEnable{})
	mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/team/disable")).Register(sessionManager.Start).Handle(&controllers.AdminTeamDisable{})
	fmt.Printf("admin url is http://localhost:%s/%s\n", os.Getenv("APP_PORT"), os.Getenv("APP_ADMIN_HASH"))

	// Run!!
	app.Run(iris.Addr(":" + os.Getenv("APP_PORT")))
}
