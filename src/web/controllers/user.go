package controllers

import (
	"errors"
	"fmt"
	"regexp"

	"../models"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

// UserController is our /user controller.
// UserController is responsible to handle the following requests:
// GET              /user/register
// POST             /user/register
// GET              /user/login
// POST             /user/login
// GET              /user/me
// All HTTP Methods /user/logout
type UserController struct {
	// mvc.C is just a lightweight lightweight alternative
	// to the "mvc.Controller" controller type,
	// use it when you don't need mvc.Controller's fields
	// (you don't need those fields when you return values from the method functions).
	mvc.C

	// Our UserService, it's an interface which
	// is binded from the main application.
	// Service services.UserService

	// Session-relative things.
	Manager *sessions.Sessions
	Session *sessions.Session
}

// BeginRequest will set the current session to the controller.
//
// Remember: iris.Context and context.Context is exactly the same thing,
// iris.Context is just a type alias for go 1.9 users.
// We use context.Context here because we don't need all iris' root functions,
// when we see the import paths, we make it visible to ourselves that this file is using only the context.
func (c *UserController) BeginRequest(ctx context.Context) {
	c.C.BeginRequest(ctx)

	if c.Manager == nil {
		ctx.Application().Logger().Errorf(`UserController: sessions manager is nil, you should bind it`)
		ctx.StopExecution() // dont run the main method handler and any "done" handlers.
		return
	}

	c.Session = c.Manager.Start(ctx)
	if c.Session.Get("username") != nil {
		txt := fmt.Sprintf("login user is %s", c.Session.Get("username"))
		fmt.Println(txt)
	}
}

func (c *UserController) isLoggedIn() bool {
	return c.Session.Get("username") != nil
}

func (c *UserController) logout() {
	c.Session.Clear()
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}

	return mvc.View{
		Name: "user/register.html",
		Data: context.Map{"Title": "User Registration"},
	}
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	// get firstname, username and password from the form.
	var (
		username              = c.Ctx.FormValue("username")
		email                 = c.Ctx.FormValue("email")
		password              = c.Ctx.FormValue("password")
		password_confirmation = c.Ctx.FormValue("password_confirmation")
	)
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		err := errors.New("User name is 'a-z A-Z 0-9' only")
		return mvc.Response{Err: err}
	} else if password != password_confirmation {
		err := errors.New("Passwords do not match")
		return mvc.Response{Err: err}
	}

	userModel := UserModel.New()
	used, err := userModel.UsedChack(username, email)
	if used {
		err := errors.New("username or email is already used")
		return mvc.Response{Err: err}
	}
	err = userModel.Add(username, email, password)
	if err != nil {
		fmt.Println(err)
		err := errors.New("db insert error")
		return mvc.Response{Err: err}
	}
	// create the new user, the password will be hashed by the service.
	// u, err := c.Service.Create(password, datamodels.User{
	//     Name:  username,
	// })

	// set the user's id to this session even if err != nil,
	// the zero id doesn't matters because .getCurrentUserID() checks for that.
	// If err != nil then it will be shown, see below on mvc.Response.Err: err.
	c.Session.Set("username", username)

	return mvc.Response{
		// if not nil then this error will be shown instead.
		Err: err,
		// redirect to /user/me.
		Path: "/user/me",
		// When redirecting from POST to GET request you -should- use this HTTP status code,
		// however there're some (complicated) alternatives if you
		// search online or even the HTTP RFC.
		// Status "See Other" RFC 7231, however iris can automatically fix that
		// but it's good to know you can set a custom code;
		// Code: 303,
	}

}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	if c.isLoggedIn() {
		// if it's already logged in then destroy the previous session.
		c.logout()
	}

	return mvc.View{
		Name: "user/login.html",
		Data: context.Map{
			"Title": "User Login",
		},
	}
}

// PostLogin handles POST: http://localhost:8080/user/login.
func (c *UserController) PostLogin() mvc.Result {
	var (
		email    = c.Ctx.FormValue("email")
		password = c.Ctx.FormValue("password")
	)

	userModel := UserModel.New()
	username, err := userModel.PasswordCheck(email, password)
	if err != nil {
		return mvc.Response{Err: errors.New("email or username is incorrect")}
	}

	c.Session.Set("username", username)

	return mvc.Response{
		Path: "/user/me",
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	username := c.Session.Get("username")
	fmt.Println(username)
	if !c.isLoggedIn() {
		// if it's not logged in then redirect user to the login page.
		return mvc.Response{Path: "/user/login"}
	}
	// username := c.Session.Get("username")
	fmt.Println(username)

	return mvc.View{
		Name: "user/me.html",
		Data: context.Map{
			"User": username,
		},
	}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}
	c.Ctx.Redirect("/")
}
