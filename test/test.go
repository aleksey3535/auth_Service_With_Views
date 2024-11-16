package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("Secret_Key"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.New(r, "test")
	if err != nil {
		panic(err)
	}
	session.Values["name"] = "alexey"
	session.Values["age"] = "123"
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
	w.Write([]byte("<h1>Main handler</h1>"))


}

func sessionInfo(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "test")
	if err != nil {
		panic(err)
	}
	data := fmt.Sprint("Sessions values: ", session.Values["name"], session.Values["age"], session.Values["surname"])
	w.Write([]byte(data))

}


func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/info", sessionInfo)
	if err := http.ListenAndServe(":1234", mux); err != nil {
		panic(err)
	}
}