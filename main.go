package main

import (
	"fmt"
	"log"
	"net/http"
	"wind/api"

	"github.com/gorilla/websocket"
)

func handleRequests() {
	http.HandleFunc("/", api.BaseHandler)
	http.HandleFunc("/create", api.CreateHandler)
	http.HandleFunc("/join", api.JoinHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	println("Starting REST API...")
	handleRequests()
}

//func main() {
//	//go func() {
//	//	println("Starting Publisher...")
//	//	pubListenAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8080")
//	//	if err != nil {
//	//		panic(err)
//	//	}
//	//	pub := windmq.NewPublisher(pubListenAddr)
//	//	pub.Start()
//	//	defer pub.Close()
//	//
//	//	for {
//	//		pub.Send([]byte("Hello World!"))
//	//		time.Sleep(5 * time.Second)
//	//	}
//	//}()
//	//
//	//go func() {
//	//	println("Starting Subscriber...")
//	//	pubAddr, err := net.ResolveTCPAddr("tcp", "10.0.0.14:8080")
//	//	if err != nil {
//	//		panic(err)
//	//	}
//	//	sub := windmq.NewSubscriber(pubAddr, 1024)
//	//	sub.Start()
//	//	defer sub.Close()
//	//
//	//	for {
//	//		message := sub.EnsureReceived()
//	//		fmt.Println(string(message))
//	//	}
//	//}()
//
//
//	println("Starting REST API...")
//	handleRequests()
//}
