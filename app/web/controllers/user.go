package controllers

import (
	"errors"
	"regexp"

	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/TeamModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/datamodels/UserModel"
	"github.com/HayatoDoi/HarekazeCTF-Competition/app/web/controllers/BaseController"
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
			"Title":       "User Registration",
			"Token":       c.MakeToken("/user/register"),
			"currentPage": "register",
		},
	}
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	// get firstname, username and password from the form.
	var (
		username                        = c.Ctx.FormValue("username")
		email                           = c.Ctx.FormValue("email")
		password                        = c.Ctx.FormValue("password")
		password_confirmation           = c.Ctx.FormValue("password_confirmation")
		makejointeam                    = c.Ctx.FormValue("makejointeam")
		make_team_name                  = c.Ctx.FormValue("make_team_name")
		make_team_password              = c.Ctx.FormValue("make_team_password")
		make_team_password_confirmation = c.Ctx.FormValue("make_team_password_confirmation")
		join_team_name                  = c.Ctx.FormValue("join_team_name")
		join_team_password              = c.Ctx.FormValue("join_team_password")
		token                           = c.Ctx.FormValue("csrf_token")
	)

	if !c.CheckTaken(token, "/user/register") {
		err := errors.New("token error!!")
		return mvc.Response{Err: err, Code: 400}
	}
	// validation check
	if username == "" || email == "" || password == "" {
		return mvc.Response{Err: errors.New("user name or user email or user password is null")}
	} else if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		err := errors.New("User name is 'a-z A-Z 0-9' only")
		return mvc.Response{Err: err}
	} else if !regexp.MustCompile(`^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$`).MatchString(email) {
		err := errors.New("The format of the email is incorrect.")
		return mvc.Response{Err: err}
	} else if password != password_confirmation {
		err := errors.New("User passwords do not match")
		return mvc.Response{Err: err}
	} else if makejointeam != "join_team" && makejointeam != "make_team" {
		return mvc.Response{Err: errors.New("unkown radio status")}
	} else if makejointeam == "make_team" {
		if make_team_name == "" || make_team_password == "" {
			return mvc.Response{Err: errors.New("team name or team password is null")}
		} else if make_team_password != make_team_password_confirmation {
			return mvc.Response{Err: errors.New("Team passwords do not match")}
		}
	}

	teamModel := TeamModel.New()
	userModel := UserModel.New()

	// team password check
	if makejointeam == "join_team" {
		status, err := teamModel.PasswordCheck(join_team_name, join_team_password)
		if err != nil {
			return mvc.Response{Err: err}
		}
		if !status {
			return mvc.Response{Err: errors.New("Team password or Team name do not match")}
		}
	}

	//used check
	usernameUsed, err := userModel.UsedChack(username, email)
	if usernameUsed {
		err := errors.New("username or email is already used")
		return mvc.Response{Err: err}
	}
	if makejointeam == "make_team" {
		teamnameUsed, err := teamModel.UsedChack(make_team_name)
		if err != nil {
			return mvc.Response{Err: err}
		}
		if teamnameUsed {
			return mvc.Response{Err: errors.New("teamname is already used")}
		}
	}

	//add team
	if makejointeam == "make_team" {
		err = teamModel.Add(make_team_name, make_team_password)
		if err != nil {
			return mvc.Response{Err: err}
		}
	}
	//get team name
	teamName := join_team_name
	if makejointeam == "make_team" {
		teamName = make_team_name
	}
	//add user
	err = userModel.Add(username, email, password, teamName)
	if err != nil {
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
		Path: c.GetRedirectPath(),
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
			"Title":       "User Login",
			"Token": c.MakeToken("/user/login"),
			"currentPage": "login",
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

	if !c.CheckTaken(token, "/user/login") {
		err := errors.New("token error!!")
		return mvc.Response{Err: err, Code: 400}
	} else if !regexp.MustCompile(`^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$`).MatchString(email) {
		err := errors.New("The format of the email is incorrect.")
		return mvc.Response{Err: err}
	}

	userModel := UserModel.New()
	existence, err := userModel.PasswordCheck(email, password)
	if err != nil {
		return mvc.Response{Err: err}
	}
	if existence == false {
		return mvc.Response{Err: errors.New("user not found")}
	}
	username, err := userModel.GetNameFromEmail(email)

	if err != nil {
		return mvc.Response{Err: err}
	}
	if c.LoginUser(username) != nil {
		return mvc.Response{Err: errors.New("session error")}
	}

	return mvc.Response{
		Path: c.GetRedirectPath(),
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.IsLoggedIn() {
		// if it's not logged in then redirect user to the login page.
		return mvc.Response{Path: "/user/login"}
	}

	return mvc.View{
		Name: "user/me.html",
		Data: context.Map{
			"UserName": c.GetLoggedUserName(),
			"UserID":   c.GetLoggedUserID(),
			"TeamName": c.GetLoggedTeamName(),
			"TeamID":   c.GetLoggedTeamID(),
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
