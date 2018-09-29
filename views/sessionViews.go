package views

import (
	//"errors"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	//"html/template"
	"log"
	"net/http"

	"github.com/igorbalden/crudtool/services/ibsession"
	"github.com/igorbalden/crudtool/services/mysqlsrv"
	"github.com/igorbalden/crudtool/utils"
)

func init() {
}

//LoginFunc implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	viewData := viewDataInit()
	viewData["title"] = "Login"
	sess, err := ibsession.GetSession(w, r)
	//err = errors.New("New Session error")
	if err != nil {
		log.Println("error identifying session")
		viewData["Message"] = "Error identifying session."
		loginTemplate.Execute(w, viewData)
		return
	}

	switch r.Method {
	case "GET":
		if sess.Get("loggedin") == "true" {
			http.Redirect(w, r, "/dbselect", 302)
			return
		}
		viewData["csrfToken"] = utils.MakeCsrf(w, r)
		loginTemplate.Execute(w, viewData)
	case "POST":
		r.ParseForm()
		if !utils.ChkCsrf(w, r) {
			viewData["Message"] = "Form data error"
			viewData["csrfToken"] = utils.MakeCsrf(w, r)
			loginTemplate.Execute(w, viewData)
			return
		}
		username := r.Form.Get("username")
		dbname := r.Form.Get("dbname")
		myConn := new(mysqlsrv.Mysqlconn)
		sess.Set("MyConn", myConn)
		err = myConn.Connect(w, r)
		if err != nil {
			viewData["Message"] = "Authentication failed for user " + username
			viewData["csrfToken"] = utils.MakeCsrf(w, r)
			loginTemplate.Execute(w, viewData)
			return
		}

		sess.Set("loggedin", "true")
		sess.Set("username", username)
		log.Print("Connected, session: ", sess.SessionID())

		if dbname == "" {
			http.Redirect(w, r, "/dbselect", 302)
		} else {
			http.Redirect(w, r, "/dbname/"+dbname, 302)
		}
		return
	default:
		http.Redirect(w, r, "/login", http.StatusUnauthorized)

	}
}

//LogoutFunc logout and destroy session
func LogoutFunc(w http.ResponseWriter, r *http.Request) {
	sess, _ := ibsession.GetSession(w, r)
	log.Print("Disconnected, session: ", sess.SessionID())
	ibsession.GlobalSessions.SessionDestroy(w, r)
	http.Redirect(w, r, "/login", 302) //redirect to login irrespective of error or not
}
