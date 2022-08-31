package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
Create a graph showing the connections between Wind, Servers, and Clients
*/
var server map[int]int

func baseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UP"))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	println("creating game")
	id := r.URL.Query().Get("id")
	println("ID: " + id)
	fmt.Println(server)
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
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

func handleRequests() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/join", joinHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	println("Started Server...")
	server = make(map[int]int, 0)
	handleRequests()
}
