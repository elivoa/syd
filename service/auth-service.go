package service

import (
	"github.com/elivoa/got/db"
	"github.com/elivoa/got/logs"
	"github.com/elivoa/syd/dal/userdao"
	"log"
	"net/http"

	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"gopkg.in/session.v1"
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
	s.Srv.SetInternalErrorHandler(func(err error) {
		log.Println("+++++++++++++++++++++++++++++++----------------------\n[oauth2] error:", err.Error())
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
