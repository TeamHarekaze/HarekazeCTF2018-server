package main

import (
    "os"
    "time"
    "log"
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
    app.Controller("/", new(controllers.HomeController))
    app.Controller("/user", new(controllers.UserController), sessionManager)

    // Run!!
    app.Run(iris.Addr(":" + os.Getenv("APP_PORT")))
}
