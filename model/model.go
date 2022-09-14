package model

type SourceType string
type EntityType string
type MessageType string

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	GridSize     = 10

	Server SourceType = "server"
	Client SourceType = "client"

	Player EntityType = "player"
	Apple  EntityType = "apple"

	Join    MessageType = "join"
	Move    MessageType = "move"
	Unknown MessageType = "unknown"
)

const (
	DirNone = iota
	DirLeft
	DirRight
	DirDown
	DirUp
)

type Entity struct {
	X             int
	Y             int
	MoveDirection int
	Type          EntityType
}

type State struct {
	Timer    int
	MoveTime int
	Lag      float64
	Player   Entity
	Follow   Entity
}
