package handler

import (
	"muhammadyasir-dev/cmd/apis"
	"net/http"
)

func PsuedoTerminal(w http.ResponseWriter, r *http.Request) {
	apis.Streampty(w, r)
}
func Signup(w http.ResponseWriter, r *http.Request) {
	apis.Signup(w, r)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	apis.LoginHandler(w, r)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	apis.CallbackHandler(w, r)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	apis.GetUserHandler(w, r)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	apis.LogoutHandler(w, r)
}
