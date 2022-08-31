package dto

import "encoding/binary"

type Frame struct {
	message []byte
}

func RestoreFrame(frame []byte) (int, *Frame) {
	if len(frame) < 4 {
		return 0, nil
	}
	messageSize := binary.BigEndian.Uint32(frame[:4])
	if len(frame)-4 < int(messageSize) {
		return 0, nil
	}

	frameSize := messageSize + 4
	return int(frameSize), &Frame{
		message: frame[4:frameSize],
	}
}

func NewFrame(message []byte) *Frame {
	return &Frame{
		message: message,
	}
}

func (this *Frame) GetBytes() []byte {
	messageSize := uint32(len(this.message))
	bytes := make([]byte, 4+messageSize)
	binary.BigEndian.PutUint32(bytes[:4], messageSize)
	copy(bytes[4:], this.message)
	return bytes
}

func (this *Frame) GetMessage() []byte {
	return this.message
}
