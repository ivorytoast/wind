package main

import (
	"fmt"
	"net"
	"wind/windmq"
)

func main() {
	pubAddr, err := net.ResolveTCPAddr("tcp", "localhost:9090")
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
}
