package syd

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/core/exception"
	"github.com/elivoa/got/errorhandler"
	"github.com/elivoa/got/utils"
	"github.com/elivoa/gxl"
	"reflect"
)

// TODO: think out a better way to register this.
var Module = &core.Module{
	Name:            "syd server", // Don't use this. It's only used to display.
	Version:         "2.0",        // TODO: used to add to assets path to disable cache.
	VarName:         "Module",     // Variable name.
	BasePath:        utils.CurrentBasePath(),
	PackagePath:     "github.com/elivoa/syd", // package name used anywhere to locate important things.
	Description:     "SYD Platform Server side. --Secure api.",
	IsStartupModule: true, // Application only accept one startup module for now.
	Register: func() {
		// settings
		c := config.Config

		// config static resources
		// c.AddStaticResource("/pictures/", "/var/site/data/syd/pictures/")
		c.AddStaticResource("/static/", "static/") // TODO: test this, is this works now?

		c.Port = 8880
		c.SetDBInfo(3306, "sydplatform", "root", "eserver409$)(")

		gxl.Locale = gxl.CN

		// builtin functions
		// templates.RegisterFunc("HasAnyRole", HasAnyRole)

		// errorhandlers
		errorhandler.AddHandler("LoginError",
			reflect.TypeOf(LoginError{}),
			errorhandler.RedirectHandler("/account/login"),
		)
		errorhandler.AddHandler("TimeZoneNotFoundError",
			reflect.TypeOf(exception.TimeZoneNotFoundError{}),
			errorhandler.RedirectHandler("/account/login"),
		)

		// errorhandler.AddHandler("AccessDenied",
		// 	reflect.TypeOf(AccessDeniedError{}),
		// 	errorhandler.HandleAccessDeniedError,
		// )

		config.ReloadTemplate = true // disable reload template?
	},
}

// func HasAnyRole(w http.ResponseWriter, r *http.Request, roles ...string) bool {
// 	session := sessions.LongCookieSession(r)
// 	if userTokenRaw, ok := session.Values[config.USER_TOKEN_SESSION_KEY]; ok && userTokenRaw != nil {
// 		if userToken := userTokenRaw.(*model.UserToken); userToken != nil {
// 			// TODO: check if userToken is outdated.
// 			if outdated := false; !outdated {
// 				// TODO: update userToken.Tiemout
// 				// userToken := service.UserService.GetLogin(w, r)
// 				if userToken.Roles != nil {
// 					for _, requiredRole := range roles {
// 						requiredRole = strings.ToLower(requiredRole)
// 						for _, role := range userToken.Roles {
// 							if strings.ToLower(role) == requiredRole {
// 								return true
// 							}
// 						}
// 					}
// 				}

// 			}
// 		}
// 	}
// 	return false
// }
