package auth

import (
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service"
)

type Login struct {
	core.Page

	Referer string `query:"referer"` // return here if non-empty

	// old things.
	//
	Title       string
	LoginUser   *model.User
	FormMessage string // `scope:"flash"` //
	FormError   string `query:"errmsg"` // use query to immulate Flash message.
}

func (p *Login) SetupRender() {

}

func (p *Login) OnSuccessFromLoginForm() *exit.Exit {
	// fmt.Printf("-------------- login form success -----------------\n")
	// fmt.Println("Username ", p.LoginUser)

	{ // special process login
		us, err := service.Auth.GlobalSessions().SessionStart(p.W, p.R)
		if err != nil {
			return exit.Error(err)
		}
		us.Set("LoggedInUserID", "000000")
		return exit.Redirect("/auth/auth")
	}

	// temporally disabled.
	_, err := service.User.Login(p.LoginUser.Username, p.LoginUser.Password, p.W, p.R)
	if err != nil {
		// error can't login, How to redirect to the current page and show errors.
		p.FormError = "Error: Login failed!"

		// TODO: immulate flash message. automatically return empty page with parameter.
		url := services.Link.GeneratePageUrlWithContextAndQueryParameters("account/login",
			map[string]interface{}{"errmsg": "Login failed! " + err.Error()},
		)
		return exit.Redirect(url) // return nil // <-- should return nil
	} else {
		// service already set userToken to session and cookie. redirect if needed.

		p.FormMessage = "Login Success!" // nouse! No one can see this.
		return exit.Redirect("/")        // Return to homepage; TODO: return to where I comes from!
	}
}

// TODO: Is this right?
// Usages: /auth/login:logout?refer=http.....
func (p *Login) Onlogout() *exit.Exit {
	fmt.Println("Service logout")
	err := service.Auth.OAuthDeleteToken(p.W, p.R)
	if nil != err {
		return exit.Error(err)
	}
	url := route.GetRefererFromURL(p.R)
	return exit.RedirectFirstValid(url, "/")
}
