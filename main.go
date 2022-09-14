package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
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
var lastRequest = 0

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
	println("Starting REST API v0.2...")
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

			message := string(msg)

			msgType := determineMessageType(message)
			if msgType == model.Join {
				handleSetting(message)
			} else if msgType == model.Move {
				handleEvent(message)
			}
		}
	})

	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				sendGameState()
			}
		}
	}()

	newGame()
	handleRequests()
}

func sendGameState() {
	lastProcessedClientRequest := strconv.Itoa(lastRequest)
	x := strconv.Itoa(serverGame.State.Player.X)
	y := strconv.Itoa(serverGame.State.Player.Y)
	message := publisher.CreateMessage(lastProcessedClientRequest, "player", "move", x+"|"+y)
	publisher.Send(message)
}

func determineMessageType(msg string) model.MessageType {
	messageParts := strings.Split(msg, ",")

	if len(messageParts) == 2 {
		return model.Join
	} else if len(messageParts) == 4 {
		return model.Move
	}

	return model.Unknown
}

func handleSetting(message string) {
	messageParts := strings.Split(message, ",")

	entity := messageParts[0]
	action := messageParts[1]

	if entity == "player" {
		if action == "join" {
			x := strconv.Itoa(serverGame.State.Player.X)
			y := strconv.Itoa(serverGame.State.Player.Y)
			response := publisher.SendJoinResponse("player", x+"|"+y)
			publisher.Send(response)
		}
	}
}

func handleEvent(msg string) {
	messageParts := strings.Split(msg, ",")

	count := messageParts[0]
	entity := messageParts[1]
	action := messageParts[2]
	detail := messageParts[3]

	lastSeenClientRequest, _ := strconv.Atoi(count)
	lastRequest = lastSeenClientRequest

	if entity == "player" {
		if action == "move" {
			switch detail {
			case "left":
				pos := serverGame.State.Player
				pos.X--
				serverGame.State.Player = pos
			case "right":
				pos := serverGame.State.Player
				pos.X++
				serverGame.State.Player = pos
			case "down":
				pos := serverGame.State.Player
				pos.Y++
				serverGame.State.Player = pos
			case "up":
				pos := serverGame.State.Player
				pos.Y--
				serverGame.State.Player = pos
			case "none":
				pos := serverGame.State.Player
				serverGame.State.Player = pos
			}
		}
	}
}
