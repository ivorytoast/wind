package windmq

import (
	"net"
	"sync"
	"wind/windmq/dto"
)

type Consumer struct {
	listener            *net.TCPListener
	rxBufferSize        int
	consumerWorkers     []*ConsumerWorker
	consumerWorkerMutex sync.RWMutex
	rxFrames            chan *dto.Frame
	isClosed            bool
}

func NewConsumer(
	listenAddress *net.TCPAddr,
	rxBufferSize int,
) *Consumer {
	listener, err := net.ListenTCP("tcp", listenAddress)
	if err != nil {
		panic(err)
	}
	this := &Consumer{
		listener:     listener,
		rxBufferSize: rxBufferSize,
		rxFrames:     make(chan *dto.Frame),
		isClosed:     false,
	}
	return this
}

func (this *Consumer) Start() {
	go func() {
		for {
			conn, err := this.listener.AcceptTCP()
			if err != nil {
				if this.isClosed {
					return
				}
				panic(err)
			}
			this.consumerWorkerMutex.Lock()
			consumerWorker := NewConsumerWorker(
				this,
				conn,
				this.rxBufferSize,
			)
			consumerWorker.Start()
			this.consumerWorkers = append(this.consumerWorkers, consumerWorker)
			this.consumerWorkerMutex.Unlock()
		}
	}()
}

func (this *Consumer) Close() {
	if this.isClosed {
		return
	}
	this.isClosed = true
	this.listener.Close()
	this.consumerWorkerMutex.Lock()
	for _, consumerWorker := range this.consumerWorkers {
		consumerWorker.Close()
	}
	this.consumerWorkers = make([]*ConsumerWorker, 0)
	this.consumerWorkerMutex.Unlock()
	close(this.rxFrames)
}

func (this *Consumer) ProducerClosed(consumerWorker *ConsumerWorker) {
	this.consumerWorkerMutex.Lock()
	for i, cw := range this.consumerWorkers {
		if cw == consumerWorker {
			this.consumerWorkers = append(this.consumerWorkers[:i], this.consumerWorkers[i+1:]...)
			break
		}
	}
	this.consumerWorkerMutex.Unlock()
}

func (this *Consumer) ReceivedFrame(frame *dto.Frame) {
	if frame == nil {
		panic("frame is nil")
	}

	if this.isClosed {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			// this object has been closed, ignore this frame
			return
		}
	}()
	this.rxFrames <- frame
}

// return nil if closed
func (this *Consumer) Receive() []byte {
	if this.isClosed {
		// object has been closed
		return nil
	}
	frame := <-this.rxFrames
	if this.isClosed {
		// channel has been closed
		return nil
	}
	if frame == nil {
		panic("Received nil frame")
	}
	return frame.GetMessage()
}
