package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"wind/windmq"
)

/*
Create a graph showing the connections between Wind, Servers, and Clients
*/
var server map[int]int

func baseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UP! Github Actions Worked! v2 :)"))
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
	go func() {
		println("Starting Publisher...")
		pubListenAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8080")
		if err != nil {
			panic(err)
		}
		pub := windmq.NewPublisher(pubListenAddr)
		pub.Start()
		defer pub.Close()

		for {
			pub.Send([]byte("Hello World!"))
			time.Sleep(1 * time.Second)
		}
	}()

	println("Starting REST API...")
	server = make(map[int]int, 0)
	handleRequests()
}
