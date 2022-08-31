package model

type SourceType string

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	GridSize     = 10

	Server SourceType = "server"
	Client SourceType = "client"
)

type Position struct {
	X int
	Y int
}

type State struct {
	ID            int
	ServerID      int
	Players       map[int]Position
	MoveDirection int
	Apple         Position
	Timer         int
	MoveTime      int
	Score         int
	BestScore     int
	Level         int
}
