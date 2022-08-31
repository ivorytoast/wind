package windmq

import (
	"net"
	"wind/windmq/dto"
	"wind/windmq/utils"
)

type ConsumerWorker struct {
	consumer     *Consumer
	producer     *net.TCPConn
	rxBuffer     []byte
	rxBytesSoFar []byte
}

func NewConsumerWorker(
	consumer *Consumer,
	producer *net.TCPConn,
	rxBufferSize int,
) *ConsumerWorker {
	this := &ConsumerWorker{
		consumer: consumer,
		producer: producer,
		rxBuffer: make([]byte, rxBufferSize),
	}
	return this
}

func (this *ConsumerWorker) Start() {
	go func() {
		for {
			frame, err := this.receive()
			if err != nil {
				this.consumer.ProducerClosed(this)
				return
			}
			if frame == nil {
				continue
			}
			this.consumer.ReceivedFrame(frame)
		}
	}()
}

func (this *ConsumerWorker) Close() {
	this.producer.Close()
}

// notice: the returned slice will not be reused in the next call
func (this *ConsumerWorker) receive() (*dto.Frame, error) {
	rxBytesSoFar, frame, err := utils.Receive(
		this.producer,
		this.rxBuffer,
		this.rxBytesSoFar,
	)
	this.rxBytesSoFar = rxBytesSoFar
	return frame, err
}
