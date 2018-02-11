package controllers

import (
	"regexp"

	"github.com/TeamHarekaze/HarekazeCTF2018-server/datamodels/TeamModel"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/datamodels/UserModel"
	"github.com/TeamHarekaze/HarekazeCTF2018-server/web/controllers/BaseController"
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

var (
	username                        string
	email                           string
	password                        string
	password_confirmation           string
	makejointeam                    string
	make_team_name                  string
	make_team_password              string
	make_team_password_confirmation string
	join_team_name                  string
	join_team_password              string
	token                           string
)

func (c *UserController) registerTemplete(msgs ...string) mvc.Result {
	if len(msgs) == 0 {
		return mvc.View{
			Name: "user/register.html",
			Data: context.Map{
				"Title":       "User Registration",
				"Token":       c.MakeToken("/user/register"),
				"IsLoggedIn":  c.IsLoggedIn(),
				"CurrentPage": "register",
			},
		}
	} else {
		return mvc.View{
			Name: "user/register.html",
			Data: context.Map{
				"Err":         msgs[0],
				"Username":    username,
				"Email":       email,
				"Title":       "User Registration",
				"Token":       c.MakeToken("/user/register"),
				"IsLoggedIn":  c.IsLoggedIn(),
				"CurrentPage": "register",
			},
		}
	}
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.IsLoggedIn() {
		c.Logout()
	}
	return c.registerTemplete()
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	// get firstname, username and password from the form.
	username = c.Ctx.FormValue("username")
	email = c.Ctx.FormValue("email")
	password = c.Ctx.FormValue("password")
	password_confirmation = c.Ctx.FormValue("password_confirmation")
	makejointeam = c.Ctx.FormValue("makejointeam")
	make_team_name = c.Ctx.FormValue("make_team_name")
	make_team_password = c.Ctx.FormValue("make_team_password")
	make_team_password_confirmation = c.Ctx.FormValue("make_team_password_confirmation")
	join_team_name = c.Ctx.FormValue("join_team_name")
	join_team_password = c.Ctx.FormValue("join_team_password")
	token = c.Ctx.FormValue("csrf_token")

	if !c.CheckTaken(token, "/user/register") {
		return c.registerTemplete("token error.")
	}
	// validation check
	if username == "" || email == "" || password == "" {
		return c.registerTemplete("User name or user email or user password is null.")
	} else if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		return c.registerTemplete("User name is 'a-z A-Z 0-9' only.")
	} else if !regexp.MustCompile(`^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$`).MatchString(email) {
		return c.registerTemplete("The format of the email is incorrect.")
	} else if password != password_confirmation {
		return c.registerTemplete("User passwords do not match")
	} else if makejointeam != "join_team" && makejointeam != "make_team" {
		return c.registerTemplete("Unkown radio status.")
	} else if makejointeam == "make_team" {
		if make_team_name == "" || make_team_password == "" {
			return c.registerTemplete("Team name or team password is null.")
		} else if make_team_password != make_team_password_confirmation {
			return c.registerTemplete("Team passwords do not match")
		}
	}

	teamModel := TeamModel.New()
	userModel := UserModel.New()

	// team password check
	if makejointeam == "join_team" {
		status, err := teamModel.PasswordCheck(join_team_name, join_team_password)
		if err != nil {
			return c.Error(err)
		}
		if !status {
			return c.registerTemplete("Team password or Team name do not match.")
		}
	}

	//used check
	usernameUsed, err := userModel.UsedChack(username, email)
	if usernameUsed {
		return c.registerTemplete("Username or email is already used.")
	}
	if makejointeam == "make_team" {
		teamnameUsed, err := teamModel.UsedChack(make_team_name)
		if err != nil {
			return c.Error(err)
		}
		if teamnameUsed {
			return c.registerTemplete("Teamname is already used.")
		}
	}

	//add team
	if makejointeam == "make_team" {
		err = teamModel.Add(make_team_name, make_team_password)
		if err != nil {
			return c.Error(err)
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
		return c.Error(err)
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

func (c *UserController) loginTemplete(msgs ...string) mvc.Result {
	if len(msgs) != 1 {
		return mvc.View{
			Name: "user/login.html",
			Data: context.Map{
				"Title":       "User Login",
				"Token":       c.MakeToken("/user/login"),
				"IsLoggedIn":  c.IsLoggedIn(),
				"CurrentPage": "login",
			},
		}
	} else {
		return mvc.View{
			Name: "user/login.html",
			Data: context.Map{
				"Err":         msgs[0],
				"Email":       email,
				"Title":       "User Login",
				"Token":       c.MakeToken("/user/login"),
				"IsLoggedIn":  c.IsLoggedIn(),
				"CurrentPage": "login",
			},
		}
	}
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	if c.IsLoggedIn() {
		// if it's already logged in then destroy the previous session.
		c.Logout()
	}
	return c.loginTemplete()
}

// PostLogin handles POST: http://localhost:8080/user/login.
func (c *UserController) PostLogin() mvc.Result {
	email = c.Ctx.FormValue("email")
	password = c.Ctx.FormValue("password")
	token = c.Ctx.FormValue("csrf_token")

	if !c.CheckTaken(token, "/user/login") {
		return c.loginTemplete("Token error.")
	} else if !regexp.MustCompile(`^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$`).MatchString(email) {
		return c.loginTemplete("The format of the email is incorrect.")
	}

	userModel := UserModel.New()
	existence, err := userModel.PasswordCheck(email, password)
	if err != nil {
		return c.Error(err)
	}
	if existence == false {
		return c.loginTemplete("User not found.")
	}
	username, err := userModel.GetNameFromEmail(email)

	if err != nil {
		return c.Error(err)
	}
	if c.LoginUser(username) != nil {
		return c.Error("Session error.")
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

	teamModel := TeamModel.New()
	member, err := teamModel.GetMember(c.GetLoggedTeamID())

	if err != nil {
		return c.Error(err)
	}

	return mvc.View{
		Name: "user/me.html",
		Data: context.Map{
			"UserName":    c.GetLoggedUserName(),
			"TeamName":    c.GetLoggedTeamName(),
			"Member":      member,
			"IsLoggedIn":  c.IsLoggedIn(),
			"CurrentPage": "me",
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
