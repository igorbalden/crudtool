package ibsession

import (
	"github.com/astaxie/beego/session"
	"net/http"
)

//GlobalSessions holds all sessions
var GlobalSessions *session.Manager
var mConf *session.ManagerConfig

func init() {
	mConf = new(session.ManagerConfig)
	mConf.CookieName = "gosessionid"
	mConf.EnableSetCookie = true
	mConf.Gclifetime = 600
	GlobalSessions, _ = session.NewManager("memory", mConf)
	go GlobalSessions.GC()
}

//GetSession make session Id available
func GetSession(w http.ResponseWriter, r *http.Request) (session.Store, error) {
	Sess, err := GlobalSessions.SessionStart(w, r)
	return Sess, err
}
