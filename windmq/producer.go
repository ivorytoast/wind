package windmq

import (
	"net"
	"time"
	"wind/windmq/dto"
)

type Producer struct {
	consumerAddress *net.TCPAddr
	consumer        *net.TCPConn
	isClosed        bool
}

func NewProducer(consumerAddress *net.TCPAddr) *Producer {
	this := &Producer{
		consumerAddress: consumerAddress,
		isClosed:        false,
	}
	return this
}

func (this *Producer) Start() {
	go this.ensureReconnected()
}

func (this *Producer) Close() {
	this.isClosed = true
	if this.consumer != nil {
		this.consumer.Close()
	}
}

// return: is successful
func (this *Producer) ensureReconnected() bool {
	if this.consumer != nil {
		this.consumer.Close()
	}
	isSucceeded := this.connect()
	for !isSucceeded {
		if this.isClosed {
			return false
		}
		time.Sleep(time.Second)
		isSucceeded = this.connect()
	}
	return true
}

// return: is successful
func (this *Producer) connect() bool {
	if this.isClosed {
		return false
	}
	conn, err := net.DialTCP("tcp", nil, this.consumerAddress)
	if err != nil {
		return false
	}
	err = conn.SetKeepAlive(true)
	if err != nil {
		panic(err)
	}
	this.consumer = conn
	return true
}

func (this *Producer) Send(message []byte) error {
	for this.consumer == nil {
		// wait for mqpub to be connected
		time.Sleep(time.Second)
	}
	frame := dto.NewFrame(message)
	frameBytes := frame.GetBytes()
	// notice: write might success even if the connection is closed by peer
	_, err := this.consumer.Write(frameBytes)
	if err != nil {
		return err
	}
	return nil
}

// return: is successful
func (this *Producer) EnsureSent(message []byte) bool {
	for {
		err := this.Send(message)
		if err != nil {
			ok := this.ensureReconnected()
			if !ok {
				// mqsub is closed
				return false
			}
			continue
		}
		return true
	}
}
