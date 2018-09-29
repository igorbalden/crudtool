package ibsession

import (
	"github.com/astaxie/beego/session"
	"github.com/igorbalden/crudtool/config"
	"net/http"
)

//GlobalSessions holds all sessions
var GlobalSessions *session.Manager
var mConf *session.ManagerConfig

func init() {
	defValues := config.Config
	mConf = new(session.ManagerConfig)
	mConf.CookieName = "gosessionid"
	mConf.EnableSetCookie = true
	mConf.Gclifetime = defValues.SessGclifetime
	GlobalSessions, _ = session.NewManager("memory", mConf)
	go GlobalSessions.GC()
}

//GetSession make session Id available
func GetSession(w http.ResponseWriter, r *http.Request) (session.Store, error) {
	Sess, err := GlobalSessions.SessionStart(w, r)
	return Sess, err
}
