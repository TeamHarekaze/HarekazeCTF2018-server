package controllers

import (
	"errors"
	"fmt"
	"regexp"

	"../models/UserModel"
	"./BaseController"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
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
	BaseController.Base
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.IsLoggedIn() {
		c.Logout()
	}

	return mvc.View{
		Name: "user/register.html",
		Data: context.Map{
			"Title": "User Registration",
			"Token": c.MakeToken(),
		},
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
		token                 = c.Ctx.FormValue("csrf_token")
	)
	if !c.CheckTaken(token) {
		err := errors.New("token error!!")
		return mvc.Response{Err: err, Code: 400}
	}
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
	c.LoginUser(username)

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
	if c.IsLoggedIn() {
		// if it's already logged in then destroy the previous session.
		c.Logout()
	}

	return mvc.View{
		Name: "user/login.html",
		Data: context.Map{
			"Title": "User Login",
			"Token": c.MakeToken(),
		},
	}
}

// PostLogin handles POST: http://localhost:8080/user/login.
func (c *UserController) PostLogin() mvc.Result {
	var (
		email    = c.Ctx.FormValue("email")
		password = c.Ctx.FormValue("password")
		token    = c.Ctx.FormValue("csrf_token")
	)

	if !c.CheckTaken(token) {
		err := errors.New("token error!!")
		return mvc.Response{Err: err, Code: 400}
	}

	userModel := UserModel.New()
	username, err := userModel.PasswordCheck(email, password)
	if err != nil {
		return mvc.Response{Err: errors.New("email or username is incorrect")}
	}

	c.LoginUser(username)

	return mvc.Response{
		Path: "/user/me",
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	username := c.LoggedUser()
	fmt.Println(username)
	if !c.IsLoggedIn() {
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
	if c.IsLoggedIn() {
		c.Logout()
	}
	c.Ctx.Redirect("/")
}
