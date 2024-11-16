package middleware

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const adminLogin = "admin"
const adminPassword = "qwerty"
const CookieSalt = "dsadjifJKZ"

type Middleware struct {
	Store * sessions.CookieStore
}

func New(store * sessions.CookieStore) *Middleware {
	return &Middleware{Store : store }
}

func (mw *Middleware) Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/success" {
			logrus.Printf("user=%s successful logged", r.RemoteAddr)
		} else {
			logrus.Printf("user=%s entered wrong password", r.RemoteAddr)
		}
		next.ServeHTTP(w, r)
	})
}

func(mw *Middleware) IsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := r.FormValue("login")
		password := r.FormValue("password")
		
		if login == adminLogin && password == adminPassword {
			session, err := mw.Store.Get(r, "myApp")
			if err != nil {
				logrus.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			session.Values["admin"] = HashCookie("admin")
			if err := session.Save(r, w); err != nil {
				logrus.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
		
	})
}

func(mw *Middleware) AdminAccess(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := mw.Store.Get(r, "myApp")
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if session.Values["admin"] == HashCookie("admin") {
			logrus.Printf("User=%s successfully logged as admin", r.RemoteAddr)
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}

	})
}

func(mw *Middleware) AccessUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := mw.Store.Get(r, "myApp")
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if session.Values["user"] == HashCookie("user") {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

	})
}

func HashCookie(cookie string) string {
	hasher := sha1.New()
	hasher.Write([]byte(cookie))
	return fmt.Sprintf("%x", hasher.Sum([]byte(CookieSalt)))
}
