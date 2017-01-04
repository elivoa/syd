package service

import (
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/coreservice/sessions"
	"github.com/elivoa/got/db"
	"github.com/elivoa/got/logs"
	"github.com/elivoa/syd"
	"github.com/elivoa/syd/dal/userdao"
	"github.com/elivoa/syd/model"
	"golang.org/x/oauth2"
	oa2errors "gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"gopkg.in/session.v1"
	"log"
	"net/http"
)

/*
改造过程：
1. 改掉session， 他用了一个session， 然而我框架中用了gorilla/session,是否要替换掉他用的session。
   ANSWER：暂时用他的也没关系吧。
2. 改掉Store，我需要使用数据库存储。所以需要实现一个数据库的Store。

*/
type AuthService struct {
	globalSessions *session.Manager // TODO replace this with gorilla/session.
	logs           *logs.Logger     // TODO: Inject request...
	Srv            *server.Server
	// session        *gsession.Session // Long Session, Stores token.
}

func NewAuthService() AuthService {
	newService := AuthService{
		logs: logs.Get("SERVICE:USER:LoginCheck"),
	}

	// init sessions.
	globalSessions, _ := session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go globalSessions.GC()

	{ // fake
		newService.globalSessions = globalSessions
	}

	// init in new() function.
	newService.init()

	return newService
}

// func (s *AuthService) LongSession(r *http.Request) *gsession.Session {
// 	if s.session == nil {
// 		s.session = sessions.LongCookieSession(r)
// 	}
// 	return s.session
// }

func (s *AuthService) init() {
	manager := manage.NewDefaultManager()
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:8880",
	})
	manager.MapClientStorage(clientStore)

	s.Srv = server.NewServer(server.NewConfig(), manager)
	s.Srv.SetUserAuthorizationHandler(s.userAuthorizeHandler)
	s.Srv.SetInternalErrorHandler(func(err error) (re *oa2errors.Response) {
		log.Println("++++++++++++++++--------\n[oauth2] Internal Error:", err.Error())
		panic(err.Error())
	})
}
func (s *AuthService) EntityManager() *db.Entity {
	return userdao.EntityManager()
}

//////////////////////////////////////

func (s *AuthService) GlobalSessions() *session.Manager {
	return s.globalSessions
}

func (s *AuthService) userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	us, err := s.globalSessions.SessionStart(w, r)
	uid := us.Get("UserID")
	if uid == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		us.Set("Form", r.Form)
		w.Header().Set("Location", "/auth/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	us.Delete("UserID")
	return

}

// --------------------------------------------------------------------------------
// Auth Session services.

// var OAUTH_TOKEN_SESSION_KEY = "oauth_token_session_key"
// TODO Inject w and r.
func (s *AuthService) OAuthToken(w http.ResponseWriter, r *http.Request) (*model.UserToken, error) {
	session := sessions.LongCookieSession(r)
	if userTokenRaw, ok := session.Values[config.USER_TOKEN_SESSION_KEY]; ok && userTokenRaw != nil {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ GET TOKEN:", userTokenRaw)
		if userToken := userTokenRaw.(*model.UserToken); userToken != nil {
			// TODO: check if outdated or null, redirect to login.
			return userToken, nil
		}
	}
	return nil, nil
}

func (s *AuthService) OAuthUpdateToken(w http.ResponseWriter, r *http.Request, token *oauth2.Token) {
	session := sessions.LongCookieSession(r)
	if userTokenRaw, ok := session.Values[config.USER_TOKEN_SESSION_KEY]; ok && userTokenRaw != nil {
		if userToken := userTokenRaw.(*model.UserToken); userToken != nil {
			// if outdated or null, redirect to login.
			userToken.Token = token
			session.Values[config.USER_TOKEN_SESSION_KEY] = userToken
			fmt.Println("SESSION: update token to ", token)
		}
	} else {
		userToken := &model.UserToken{
			Username: "TestUserName",
			Token:    token,
		}
		session.Values[config.USER_TOKEN_SESSION_KEY] = userToken
		fmt.Println("SESSION: Create and set token to ", token)
	}
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
}

func (s *AuthService) OAuthDeleteToken(w http.ResponseWriter, r *http.Request) error {
	// TODO remove token.
	session := sessions.LongCookieSession(r)
	delete(session.Values, config.USER_TOKEN_SESSION_KEY)
	return session.Save(r, w)
}

// ----------------------------------------------------------------------------------------------------
// Require Login and Requires

func (s *AuthService) Auth(w http.ResponseWriter, r *http.Request, requiredRoles model.Roles,
) (*model.UserToken, error) {

	userToken, err := s.OAuthToken(w, r)
	if err != nil {
		return nil, err
	}

	// 任何Role都要求必须登录.如果只登录不要求任何Role是ROLE_Login:1.
	if userToken == nil {
		if requiredRoles <= 0 {
			return nil, nil
		} else {
			return nil, &syd.LoginError{Message: "Not Login. UserToken is nil."}
		}
	}

	// 位运算判断角色是否符合
	if requiredRoles&^userToken.RolesInt() > 0 {
		return nil, &syd.AccessDeniedError{
			Message: fmt.Sprintf("Require Role. (Required:%b, Given:%b)",
				requiredRoles, userToken.RolesInt()),
		}
	}

	return userToken, nil
}
