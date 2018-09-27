package auth

import (
	"net/http"

	"github.com/igorbalden/crudtool/services/ibsession"
)

//RequiresLogin is a middleware which will be used for each httpHandler to check if there is any active session
func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := ibsession.GetSession(w, r)
		if sess.Get("loggedin") != "true" {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler(w, r)
	}
}
