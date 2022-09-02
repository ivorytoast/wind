package api

import (
	"fmt"
	"net/http"
)

var server = make(map[int]int, 0)

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UP! Github Actions Worked! v2 :)"))
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	println("creating game")
	id := r.URL.Query().Get("id")
	println("ID: " + id)
	fmt.Println(server)
}

func JoinHandler(w http.ResponseWriter, r *http.Request) {
	if len(server) == 0 {
		println("no server available")
		return
	}
	println("finding server")

	// ID show be like a SSN (which has location, identifier, etc... built in)
	id := r.URL.Query().Get("id")

	println("ID: " + id)

	server[1] = 1
	fmt.Println(server)
}
