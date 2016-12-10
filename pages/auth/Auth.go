package auth

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/syd/model"
	"github.com/elivoa/syd/service"
	"gopkg.in/session.v1"
	"net/url"
)

type Auth struct {
	core.Page
	Title string

	LoginUser   *model.User
	FormMessage string // `scope:"flash"` //
	FormError   string `query:"errmsg"` // use query to immulate Flash message.
	us          session.Store
}

func (p *Auth) SetupRender() *exit.Exit {
	fmt.Println("!!!!!!!!!!!!!setup render !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	var err error
	p.us, err = service.Auth.GlobalSessions().SessionStart(p.W, p.R)
	if err != nil {
		fmt.Println("!!!!!!!!!!!!! Error !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		return exit.Error(err)
	}
	if p.us.Get("LoggedInUserID") == nil {
		fmt.Println("!!!!!!!!!!!!! Redirecty to Login !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		return exit.Redirect("/auth/login")
	}
	// default return to auth.html
	// outputHTML(w, r, "static/auth.html")
	return nil
}

func (p *Auth) OnSuccessFromAuthForm() *exit.Exit {
	fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&8888888888888888888888888")
	{ // special process login

		// first call setuprender
		setuprenderResult := p.SetupRender()
		fmt.Println("result is : ", setuprenderResult)
		if nil != setuprenderResult {
			return setuprenderResult
		}

		form := p.us.Get("Form").(url.Values)
		u := new(url.URL)
		u.Path = "/auth/authorize"
		u.RawQuery = form.Encode()

		// w.Header().Set("Location", u.String())
		// w.WriteHeader(http.StatusFound)

		p.us.Delete("Form")
		p.us.Set("UserID", p.us.Get("LoggedInUserID"))
		fmt.Println("-------------------------------------")
		fmt.Println(">>>", u.String())
		return exit.Redirect(u.String())
	}
}
