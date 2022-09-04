package main

import (
	"github.com/gorilla/websocket"
	"github.com/pebbe/zmq4"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wind/api"
	"wind/model"
	"wind/queue"
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
var publisher *zmq4.Socket

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
			Player: model.Position{
				X:             17,
				Y:             30,
				MoveDirection: model.DirNone,
			},
		},
	}
}

func handleRequests() {
	println("Starting REST API...")
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
	publisher = queue.CreatePublisher("tcp://*:5556")

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

// {entity},{action},{detail} (Phase 1)
// 10|10|20|20|20 (Phase 2)
func handleMessage(msg string) {
	messageParts := strings.Split(msg, ",")

	entity := messageParts[0]
	action := messageParts[1]
	detail := messageParts[2]

	if entity == "player" {
		if action == "move" {
			switch detail {
			case "left":
				pos := serverGame.State.Player
				pos.X--
				serverGame.State.Player = pos
				publisher.Send("Player went left to: "+strconv.Itoa(pos.X), 0)
				println("Player went left")
			case "right":
				pos := serverGame.State.Player
				pos.X++
				serverGame.State.Player = pos
				publisher.Send("Player went right to: "+strconv.Itoa(pos.X), 0)
				println("Player went right")
			case "down":
				pos := serverGame.State.Player
				pos.Y++
				serverGame.State.Player = pos
				publisher.Send("Player went down to: "+strconv.Itoa(pos.Y), 0)
				println("Player went down")
			case "up":
				pos := serverGame.State.Player
				pos.Y--
				serverGame.State.Player = pos
				publisher.Send("Player went up to: "+strconv.Itoa(pos.Y), 0)
				println("Player went up")
			}
		}
	}
}
