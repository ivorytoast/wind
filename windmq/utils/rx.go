package utils

import (
	"net"
	"wind/windmq/dto"
)

// notice: the returned slice will not be reused in the next call
func Receive(
	conn *net.TCPConn,
	rxBuffer []byte,
	rxBytesSoFar []byte,
) ([]byte, *dto.Frame, error) {
	readByteCount := 0
	rxBytesSoFar, frame := restoreFrameFromRXBytesSoFar(rxBytesSoFar)
	if frame != nil {
		return rxBytesSoFar, frame, nil
	}
	nr, err := conn.Read(rxBuffer)
	readByteCount = nr
	if err != nil {
		return rxBytesSoFar, nil, err
	}
	// it will always perform a deep copy in this case
	rxBytesSoFar = append(rxBytesSoFar, rxBuffer[:readByteCount]...)
	rxBytesSoFar, frame = restoreFrameFromRXBytesSoFar(rxBytesSoFar)
	if frame == nil {
		return rxBytesSoFar, nil, nil
	}
	return rxBytesSoFar, frame, nil
}

func restoreFrameFromRXBytesSoFar(rxBytesSoFar []byte) ([]byte, *dto.Frame) {
	frameSize, frame := dto.RestoreFrame(rxBytesSoFar)
	if frame == nil {
		return rxBytesSoFar, nil
	}
	rxBytesSoFar = rxBytesSoFar[frameSize:]
	return rxBytesSoFar, frame
}
