package main

import (
	"net"
	"time"
	"wind/windmq"
)

func main() {
	pubListenAddr, err := net.ResolveTCPAddr("tcp", ":8080")
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
}
