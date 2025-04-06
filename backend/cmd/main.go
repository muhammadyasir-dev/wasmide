package main

import (
	"muhammadyasir-dev/cmd/dbs"
	"muhammadyasir-dev/cmd/routes"
	"net/http"
)

func main() {
	dbs.Initdb()
	r := routes.Router()
	http.ListenAndServe(":8080", r)

}
