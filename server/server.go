package server

//import (
//	zmq "github.com/pebbe/zmq4"
//	"time"
//	"wind/model"
//	"wind/queue"
//)
//
//func main() {
//	publisher := queue.CreatePublisher("tcp://*:5556")
//	defer publisher.Close()
//	go func() {
//		for {
//			message := queue.CreateMessage(
//				"entity",
//				"action",
//				"detail",
//			)
//			publisher.Send(string(model.Server), zmq.SNDMORE)
//			publisher.Send("serverID", zmq.SNDMORE)
//			publisher.Send("serverLocation", zmq.SNDMORE)
//			publisher.Send(message, 0)
//			time.Sleep(1 * time.Second)
//		}
//	}()
//
//	//subscriber := queue.CreateSubscriber("tcp://localhost:5556")
//	//defer subscriber.Close()
//	//_ = subscriber.SetSubscribe("client")
//	//
//	//for {
//	//	address, _ := subscriber.Recv(0)
//	//	destination, _ := subscriber.Recv(0)
//	//	contents, _ := subscriber.Recv(0)
//	//	fmt.Printf("[%s] %s %s\n", address, destination, contents)
//	//}
//}
