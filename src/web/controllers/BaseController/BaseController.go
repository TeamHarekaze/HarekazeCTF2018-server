package BaseController

import (
	"math/rand"
	"time"

	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

// Base is super struct
type Base struct {
	mvc.C
	Manager *sessions.Sessions
	Session *sessions.Session
}

// BeginRequest will set the current session to the controller.
func (c *Base) BeginRequest(ctx context.Context) {
	c.C.BeginRequest(ctx)

	if c.Manager == nil {
		ctx.Application().Logger().Errorf(`UserController: sessions manager is nil, you should bind it`)
		ctx.StopExecution() // dont run the main method handler and any "done" handlers.
		return
	}

	c.Session = c.Manager.Start(ctx)
}

// LoggedUser is return username
func (c *Base) LoggedUser() string {
	return c.Session.GetString("username")
}

// LoginUser is login
func (c *Base) LoginUser(username string) {
	c.Session.Set("username", username)
}

// IsLoggedIn is Check login status
func (c *Base) IsLoggedIn() bool {
	return c.Session.Get("username") != nil
}

// GetLoggedUserName is Check login status
func (c *Base) GetLoggedUserName() string {
	return c.Session.GetString("username")
}

// Logout is logout
func (c *Base) Logout() {
	c.Session.Delete("username")
}

// MakeToken is generate taken and set taken in session.
func (c *Base) MakeToken() string {
	var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$%&")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 64)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	taken := string(b)
	c.Session.Set("_csrfToken", taken)
	return taken
}

// CheckTaken is check taken
func (c *Base) CheckTaken(taken string) bool {
	r := c.Session.Get("_csrfToken") == taken
	c.Session.Delete("_csrfToken")
	return r
}
