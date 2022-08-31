package windmq

import (
	"net"
	"time"
	"wind/windmq/utils"
)

type Subscriber struct {
	publisher        *net.TCPConn
	publisherAddress *net.TCPAddr
	rxBuffer         []byte
	rxBytesSoFar     []byte
	isClosed         bool
}

func NewSubscriber(
	publisherAddress *net.TCPAddr,
	rxBufferSize int,
) *Subscriber {
	this := &Subscriber{
		publisherAddress: publisherAddress,
		rxBuffer:         make([]byte, rxBufferSize),
		isClosed:         false,
	}
	return this
}

func (this *Subscriber) Start() {
	go this.ensureReconnected()
}

func (this *Subscriber) Close() {
	this.isClosed = true
	if this.publisher != nil {
		this.publisher.Close()
	}
}

// return true if successful
func (this *Subscriber) ensureReconnected() bool {
	if this.publisher != nil {
		this.publisher.Close()
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
func (this *Subscriber) connect() bool {
	if this.isClosed {
		return false
	}
	conn, err := net.DialTCP("tcp", nil, this.publisherAddress)
	if err != nil {
		return false
	}
	err = conn.SetKeepAlive(true)
	if err != nil {
		panic(err)
	}
	this.publisher = conn
	return true
}

// notice: the returned slice will not be reused in the next call
func (this *Subscriber) Receive() ([]byte, error) {
	for this.publisher == nil {
		// wait for publisher to be connected
		time.Sleep(time.Second)
	}

	rxBytesSoFar, frame, err := utils.Receive(
		this.publisher,
		this.rxBuffer,
		this.rxBytesSoFar,
	)
	this.rxBytesSoFar = rxBytesSoFar
	if err != nil {
		return nil, err
	}
	if frame == nil {
		return nil, nil
	}
	return frame.GetMessage(), err
}

// notice: the returned slice will not be reused in the next call
// return nil if closed
func (this *Subscriber) EnsureReceived() []byte {
	for {
		message, err := this.Receive()
		if err != nil {
			ok := this.ensureReconnected()
			if !ok {
				// closed
				return nil
			}
			continue
		}
		if message == nil {
			continue
		}
		return message
	}
}
