package utils

import (
	"github.com/igorbalden/crudtool/services/ibsession"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// GetRedirectURL comment
func GetRedirectURL(referer string) string {
	var redirectUrl string
	url := strings.Split(referer, "/")

	if len(url) > 4 {
		redirectUrl = "/" + strings.Join(url[3:], "/")
	} else {
		redirectUrl = "/"
	}
	return redirectUrl
}

//ARandomStr gives a random string of length l
func ARandomStr(l int) string {
	const charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"

	var b []byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		b = append(b, charSet[rand.Intn(len(charSet))])
	}
	return string(b)
}

//MakeCsrf sets the csrf token in session, and returns the token
func MakeCsrf(w http.ResponseWriter, r *http.Request) string {
	tkn := ARandomStr(32)
	sess, _ := ibsession.GetSession(w, r)
	sess.Set("csrfToken", tkn)
	return tkn
}

//ChkCsrf compares the csrf values of form, and session.
func ChkCsrf(w http.ResponseWriter, r *http.Request) bool {
	sess, _ := ibsession.GetSession(w, r)
	csrfSess := sess.Get("csrfToken")
	if csrfSess == r.Form.Get("csrfToken") {
		return true
	}
	return false
}
