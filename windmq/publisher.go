package windmq

import (
	"net"
	"sync"
	"wind/windmq/dto"
)

type Publisher struct {
	subscribers     []*net.TCPConn
	subscriberMutex sync.RWMutex
	listener        *net.TCPListener
	isClosed        bool
}

func NewPublisher(listenAddress *net.TCPAddr) *Publisher {
	listener, err := net.ListenTCP("tcp", listenAddress)
	if err != nil {
		panic(err)
	}
	this := &Publisher{
		subscribers: make([]*net.TCPConn, 0),
		listener:    listener,
		isClosed:    false,
	}
	return this
}

func (this *Publisher) Start() {
	go func() {
		for {
			conn, err := this.listener.AcceptTCP()
			if err != nil {
				if this.isClosed {
					return
				}
				panic(err)
			}
			this.subscriberMutex.Lock()
			this.subscribers = append(this.subscribers, conn)
			this.subscriberMutex.Unlock()
		}
	}()
}

func (this *Publisher) Close() {
	this.isClosed = true
	this.listener.Close()
	this.subscriberMutex.Lock()
	for _, subscriber := range this.subscribers {
		subscriber.Close()
	}
	this.subscribers = make([]*net.TCPConn, 0)
	this.subscriberMutex.Unlock()
}

func (this *Publisher) Send(message []byte) {
	frame := dto.NewFrame(message)
	frameBytes := frame.GetBytes()
	deadSubscribers := make([]int, 0)
	this.subscriberMutex.RLock()
	for i, subscriber := range this.subscribers {
		_, err := (*subscriber).Write(frameBytes)
		if err != nil {
			deadSubscribers = append(deadSubscribers, i)
		}
	}
	this.subscriberMutex.RUnlock()
	for _, i := range deadSubscribers {
		(*this.subscribers[i]).Close()
		this.subscriberMutex.Lock()
		this.subscribers = append(this.subscribers[:i], this.subscribers[i+1:]...)
		this.subscriberMutex.Unlock()
	}
}

func (publisher *Publisher) CreateMessage(entity string, action string, detail string) []byte {
	return []byte(entity + "," + action + "," + detail)
}
