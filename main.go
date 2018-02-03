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
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers"
	"github.com/HayatoDoi/HarekazeCTF-Competition/lib/md2html"
	"github.com/comail/colog"
)

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	Env_load()
	// colog.SetFormatter(&colog.StdFormatter{Colors: true})
	colog.Register()
	app := iris.New()
	app.StaticWeb("/asset", "app/web/public") //debug only
	view := iris.HTML("app/web/views", ".html")
	view.Layout("layouts/layout.html")
	view.AddFunc("md2html", md2html.Md2Html)
	app.RegisterView(view)

	// 404 error page
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.View("error/404.html")
	})

	// make session DB
	sesstionDB := redis.New(service.Config{
		Network:     "tcp",
		Addr:        fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password:    os.Getenv("REDIS_PASSWORD"),
		Database:    "1",
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: service.DefaultRedisIdleTimeout,
	})
	iris.RegisterOnInterrupt(func() {
		sesstionDB.Close()
	})

	// make session manager
	sessionManager := sessions.New(sessions.Config{
		Cookie:  "HarekazeCTF-session",
		Expires: 30 * time.Minute,
	})

	// use redis db
	sessionManager.UseDatabase(sesstionDB)

	//route
	mvc.New(app.Party("/")).Register(sessionManager.Start).Handle(&controllers.HomeController{})
	mvc.New(app.Party("/user")).Register(sessionManager.Start).Handle(&controllers.UserController{})
	mvc.New(app.Party("/question")).Register(sessionManager.Start).Handle(&controllers.QuestionController{})
	mvc.New(app.Party("/answer")).Register(sessionManager.Start).Handle(&controllers.AnswerController{})
	mvc.New(app.Party("/ranking")).Register(sessionManager.Start).Handle(&controllers.RankingController{})
	//error page
	mvc.New(app.Party("/500")).Register(sessionManager.Start).Handle(&controllers.Error500Controller{})
	mvc.New(app.Party("/404")).Register(sessionManager.Start).Handle(&controllers.Error404Controller{})
	//admin
	if os.Getenv("RUN_ADMIN_CONSOLE") == "True" {
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH"))).Register(sessionManager.Start).Handle(&controllers.Admin{})
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionList{})
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question/add")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionAdd{})
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question/edit")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionEdit{})
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/question/delete")).Register(sessionManager.Start).Handle(&controllers.AdminQuestionDelete{})
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/team/enable")).Register(sessionManager.Start).Handle(&controllers.AdminTeamEnable{})
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/team/disable")).Register(sessionManager.Start).Handle(&controllers.AdminTeamDisable{})
		mvc.New(app.Party("/" + os.Getenv("APP_ADMIN_HASH") + "/feed")).Register(sessionManager.Start).Handle(&controllers.AdminFeed{})
		fmt.Printf("admin url is http://localhost:%s/%s\n", os.Getenv("APP_PORT"), os.Getenv("APP_ADMIN_HASH"))
	}
	// Run!!
	app.Run(iris.Addr(":" + os.Getenv("APP_PORT")))
}
