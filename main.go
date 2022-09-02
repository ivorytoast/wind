package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"wind/api"
	"wind/windmq"
)

func handleRequests() {
	http.HandleFunc("/", api.BaseHandler)
	http.HandleFunc("/create", api.CreateHandler)
	http.HandleFunc("/join", api.JoinHandler)
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
			time.Sleep(1 * time.Minute)
		}
	}()

	go func() {
		println("Starting Subscriber...")
		pubAddr, err := net.ResolveTCPAddr("tcp", "45.77.153.58:8080")
		if err != nil {
			panic(err)
		}
		sub := windmq.NewSubscriber(pubAddr, 1024)
		sub.Start()
		defer sub.Close()

		for {
			message := sub.EnsureReceived()
			fmt.Println(string(message))
		}
	}()

	println("Starting REST API...")
	handleRequests()
}
