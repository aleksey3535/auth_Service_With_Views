package handler

import (
	"auth/internal/pkg"
	"auth/internal/pkg/middleware"
	"auth/internal/pkg/service"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)



type Handler struct {
	mw middleware.Middleware
	s *service.Service
}

func New(mw middleware.Middleware ,s *service.Service,) *Handler {
	return &Handler{mw: mw, s : s}
}

func(h *Handler) InitRoutes() *mux.Router {
	mux := mux.NewRouter()

	mux.HandleFunc("/", h.MainHandler)
	mux.HandleFunc("/login",h.Login)
	mux.HandleFunc("/loginHandler",h.mw.IsAdmin(h.LoginHandler)).Methods(http.MethodPost)
	mux.HandleFunc("/registry", h.Registry)
	mux.HandleFunc("/registryHandler",h.RegistryHandler).Methods(http.MethodPost)
	mux.HandleFunc("/success", h.mw.Logger(h.mw.AccessUser((h.Success))))
	mux.HandleFunc("/wrong", h.mw.Logger(h.Wrong))
	mux.HandleFunc("/delete", h.DeleteHandler).Methods(http.MethodPost)
	mux.HandleFunc("/admin", h.mw.AdminAccess(h.AdminHandler))
	return mux
}
type Data struct {
	Info string
}

func(h *Handler) MainHandler(w http.ResponseWriter, r *http.Request) {
	tm, err := template.ParseFiles("internal/pkg/views/index.html")
	getError(w, r, err)
	err = tm.Execute(w, nil)
	getError(w, r, err)
}

func(h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	tm, err := template.ParseFiles("internal/pkg/views/login.html")
	getError(w, r , err)
	err = tm.Execute(w, nil)
	getError(w, r , err)
}

func(h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	flag, _ := h.s.IsExist(login, password)
	if flag {
		session, err := h.mw.Store.Get(r, "myApp")
		getError(w, r, err)
		session.Values["user"] = middleware.HashCookie("user")
		err = session.Save(r, w)
		getError(w, r, err)
		http.Redirect(w, r, "/success", http.StatusFound)
	} else {
		http.Redirect(w, r, "/wrong", http.StatusFound)
	}
}

func(h *Handler) Registry(w http.ResponseWriter, r *http.Request) {
	tm, err := template.ParseFiles("internal/pkg/views/registry.html")
	getError(w, r , err)
	err = tm.Execute(w, nil)
	getError(w, r , err)
}

func(h *Handler) RegistryHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	err := h.s.CreateUser(login, password)
	if err != nil {
		if err.Error() == h.s.ValidateLogin || err.Error() == h.s.ValidatePassword {
			info := Data{Info: err.Error()}
			tm, _ := template.ParseFiles("internal/pkg/views/info.html")
			tm.Execute(w, info)
			return
			
		}
		getError(w, r , err)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func(h *Handler) Success(w http.ResponseWriter, r *http.Request) {
	info := Data{Info: "Success! You are logged in."}
	tm, err := template.ParseFiles("internal/pkg/views/info.html")
	getError(w, r , err)
	tm.Execute(w, info)
	getError(w, r , err)
}

func(h *Handler) Wrong(w http.ResponseWriter, r *http.Request) {
	info := Data{Info:"Wrong password or such user does not exist"}
	tm, err := template.ParseFiles("internal/pkg/views/info.html")
	getError(w, r , err)
	err = tm.Execute(w, info)
	getError(w, r , err)
	
}

type DataForAdmin struct {
	Users []pkg.User
}

func(h *Handler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.s.GetAllUsers()
	getError(w, r , err)
	data := DataForAdmin{Users: users}
	tm, err := template.ParseFiles("internal/pkg/views/admin.html")
	getError(w, r , err)
	err = tm.Execute(w, data)
	getError(w, r , err)
}

func(h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	err := h.s.DeleteUser(login)
	getError(w, r , err)
	http.Redirect(w, r, "/admin", http.StatusFound)
}


func getError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		logrus.Error(err.Error(), "---", r.URL)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}