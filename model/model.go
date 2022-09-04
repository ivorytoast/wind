package model

type SourceType string

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	GridSize     = 10

	Server SourceType = "server"
	Client SourceType = "client"
)

const (
	DirNone = iota
	DirLeft
	DirRight
	DirDown
	DirUp
)

type Position struct {
	X             int
	Y             int
	MoveDirection int
}

type State struct {
	Timer    int
	MoveTime int
	Lag      float64
	Player   Position
}
