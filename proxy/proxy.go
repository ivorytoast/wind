package proxy

import (
	zmq "github.com/pebbe/zmq4"
	"log"
)

func mainn() {
	println("Started Proxy...")

	frontend, _ := zmq.NewSocket(zmq.XSUB)
	defer frontend.Close()
	frontend.Connect("tcp://192.168.55.210:5556")

	backend, _ := zmq.NewSocket(zmq.XPUB)
	defer backend.Close()
	backend.Bind("tcp://10.1.1.0:8100")

	err := zmq.Proxy(frontend, backend, nil)
	log.Fatalln("Proxy Ended... [" + err.Error() + "]")
}
