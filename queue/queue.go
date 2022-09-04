package queue

import (
	zmq "github.com/pebbe/zmq4"
)

// endpoint = "tcp://localhost:5556"
func CreateSubscriber(endpoint string) *zmq.Socket {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	err := subscriber.Connect(endpoint)
	if err != nil {
		println("error connecting to proxy")
		return nil
	}

	if err != nil {
		println("error subscribing: " + err.Error())
		return nil
	}

	return subscriber
}

// endpoint = "tcp://*:5556"
func CreatePublisher(endpoint string) *zmq.Socket {
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		println(err.Error())
	}
	serverClientErr := publisher.Bind(endpoint)
	if serverClientErr != nil {
		println(serverClientErr.Error())
	}

	return publisher
}

func CreateMessage(entity string, action string, detail string) []byte {
	return []byte(entity + "," + action + "," + detail)
}
