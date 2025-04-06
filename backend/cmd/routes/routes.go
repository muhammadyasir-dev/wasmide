package routes

import (
	"muhammadyasir-dev/cmd/handler"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/stream", handler.PsuedoTerminal).Methods("POST")

	router.HandleFunc("/signup", handler.Signup).Methods("POST")

	router.HandleFunc("/login", handler.LoginHandler).Methods("GET")
	router.HandleFunc("/auth/callback", handler.CallbackHandler).Methods("GET")
	router.HandleFunc("/user", handler.GetUserHandler).Methods("GET")
	router.HandleFunc("/logout", handler.LogoutHandler).Methods("POST")
	return router
}
