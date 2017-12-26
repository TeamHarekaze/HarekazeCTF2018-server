package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"

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
	app.RegisterView(view)
	// make session manager
	sessionManager := sessions.New(sessions.Config{
		Cookie:  "HarekazeCTF-session",
		Expires: 1 * time.Minute,
	})

	//route
	app.Controller("/", new(controllers.HomeController), sessionManager)
	app.Controller("/user", new(controllers.UserController), sessionManager)
	//admin
	app.Controller("/"+os.Getenv("APP_ADMIN_HASH"), new(controllers.Admin), sessionManager)
	app.Controller("/"+os.Getenv("APP_ADMIN_HASH")+"/question", new(controllers.AdminQuestionList), sessionManager)
	app.Controller("/"+os.Getenv("APP_ADMIN_HASH")+"/question/add", new(controllers.AdminQuestionAdd), sessionManager)
	app.Controller("/"+os.Getenv("APP_ADMIN_HASH")+"/question/edit", new(controllers.AdminQuestionEdit), sessionManager)
	app.Controller("/"+os.Getenv("APP_ADMIN_HASH")+"/question/delete", new(controllers.AdminQuestionDelete), sessionManager)
	fmt.Printf("admin url is http://localhost:%s/%s\n", os.Getenv("APP_PORT"), os.Getenv("APP_ADMIN_HASH"))

	// Run!!
	app.Run(iris.Addr(":" + os.Getenv("APP_PORT")))
}
