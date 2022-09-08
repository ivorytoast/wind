package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"wind/api"
	"wind/model"
	"wind/windmq"
)

/*
Direct web socket connection from client to server

Client starts up -- asks for the server's IP address.
Will receive the server's IP address

Server then can interact with whatever it wants with the background

The HTTP calls should be on status and data -- not actual game logic
For example, might be to get the server's IP address
*/

type ServerGame struct {
	State model.State
}

var serverGame *ServerGame
var publisher *windmq.Publisher

const (
	screenWidth        = 640
	screenHeight       = 480
	gridSize           = 10
	xGridCountInScreen = screenWidth / gridSize
	yGridCountInScreen = screenHeight / gridSize
)

func newGame() {
	serverGame = &ServerGame{
		State: model.State{
			Lag: 0.0,
			Player: model.Entity{
				X:             17,
				Y:             30,
				MoveDirection: model.DirNone,
				Type:          model.Player,
			},
		},
	}
}

func handleRequests() {
	println("Starting REST API V0.1...")
	http.HandleFunc("/", api.BaseHandler)
	http.HandleFunc("/create", api.CreateHandler)
	http.HandleFunc("/join", api.JoinHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	pubListenAddr, err := net.ResolveTCPAddr("tcp", ":5556")
	if err != nil {
		panic(err)
	}
	publisher = windmq.NewPublisher(pubListenAddr)
	publisher.Start()
	println("Started Publisher")
	defer publisher.Close()

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			handleMessage(string(msg))
		}
	})

	newGame()
	handleRequests()
}

func handleMessage(msg string) {
	messageParts := strings.Split(msg, ",")

	count := messageParts[0]
	entity := messageParts[1]
	action := messageParts[2]
	detail := messageParts[3]

	if entity == "player" {
		if action == "move" {
			switch detail {
			case "left":
				pos := serverGame.State.Player
				pos.X--
				serverGame.State.Player = pos
				x := strconv.Itoa(pos.X)
				y := strconv.Itoa(pos.Y)
				message := publisher.CreateMessage(count, "player", "move", x+"|"+y)
				publisher.Send(message)
			case "right":
				pos := serverGame.State.Player
				pos.X++
				serverGame.State.Player = pos
				x := strconv.Itoa(pos.X)
				y := strconv.Itoa(pos.Y)
				message := publisher.CreateMessage(count, "player", "move", x+"|"+y)
				publisher.Send(message)
			case "down":
				pos := serverGame.State.Player
				pos.Y++
				serverGame.State.Player = pos
				x := strconv.Itoa(pos.X)
				y := strconv.Itoa(pos.Y)
				message := publisher.CreateMessage(count, "player", "move", x+"|"+y)
				publisher.Send(message)
			case "up":
				pos := serverGame.State.Player
				pos.Y--
				serverGame.State.Player = pos
				x := strconv.Itoa(pos.X)
				y := strconv.Itoa(pos.Y)
				message := publisher.CreateMessage(count, "player", "move", x+"|"+y)
				publisher.Send(message)
			case "none":
				pos := serverGame.State.Player
				serverGame.State.Player = pos
				x := strconv.Itoa(pos.X)
				y := strconv.Itoa(pos.Y)
				message := publisher.CreateMessage(count, "player", "move", x+"|"+y)
				publisher.Send(message)
			}
		}
	}
}
