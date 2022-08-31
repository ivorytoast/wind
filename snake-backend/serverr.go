//go:build example
// +build example

package main

//import (
//	"fmt"
//	zmq "github.com/pebbe/zmq4"
//	"math/rand"
//	"net/http"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/gorilla/websocket"
//)
//
//var game *Game
//var idCounter = 1
//var lobbies map[int]Lobby
//
//const (
//	dirNone = iota
//	dirLeft
//	dirRight
//	dirDown
//	dirUp
//)
//
//const (
//	screenWidth        = 640
//	screenHeight       = 480
//	gridSize           = 10
//	xGridCountInScreen = screenWidth / gridSize
//	yGridCountInScreen = screenHeight / gridSize
//)
//
//type Lobby struct {
//	Id   int
//	Game Game
//}
//
//type Game struct {
//	players       map[int]Position
//	listOfPlayers []int
//	moveDirection int
//	apple         Position
//	timer         int
//	moveTime      int
//	score         int
//	bestScore     int
//	level         int
//}
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//}
//
//type Position struct {
//	X int
//	Y int
//}
//
//func newGame() {
//	if game != nil {
//		return
//	}
//	game = &Game{
//		players:  make(map[int]Position),
//		apple:    Position{X: 3 * gridSize, Y: 3 * gridSize},
//		moveTime: 4,
//	}
//}
//
//func mainn() {
//
//	lobbies = make(map[int]Lobby)
//
//	context, _ := zmq.NewContext()
//	socket, _ := context.NewSocket(zmq.PUB)
//	defer socket.Close()
//	socket.Bind("tcp://*:5556")
//
//	// Seed the random number generator
//	rand.Seed(time.Now().UnixNano())
//
//	newGame()
//
//	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
//		conn, _ := upgrader.Upgrade(w, r, nil)
//
//		for {
//			_, msg, err := conn.ReadMessage()
//			if err != nil {
//				return
//			}
//
//			stringSlice := strings.Split(string(msg), ",")
//
//			if len(stringSlice) != 3 {
//				if len(stringSlice) == 2 {
//					println("HIT HERE!")
//					lobbyId := handleCreate()
//					appleX := 3 * gridSize
//					appleY := 3 * gridSize
//					// lobby, entity, noun, adjective
//					// Entity can be:
//					//	1. object
//					//  2. score
//					//  3. settings
//					socket.Send(strconv.Itoa(lobbyId)+",apple,position,"+strconv.Itoa(appleX)+"|"+strconv.Itoa(appleY), 0)
//					socket.Send(strconv.Itoa(lobbyId)+",settings,moveTime,4", 0)
//					socket.Send(strconv.Itoa(lobbyId)+",settings,level,1", 0)
//					socket.Send("1,lobby,"+strconv.Itoa(lobbyId), 0)
//				}
//				continue
//			}
//			player, _ := strconv.Atoi(stringSlice[0])
//			action := stringSlice[1]
//			value := stringSlice[2]
//
//			if action == "create" {
//				lobbyId := handleCreate()
//				appleX := 3 * gridSize
//				appleY := 3 * gridSize
//				// lobby, entity, noun, adjective
//				// Entity can be:
//				//	1. object
//				//  2. score
//				//  3. settings
//				socket.Send(strconv.Itoa(lobbyId)+",apple,position,"+strconv.Itoa(appleX)+"|"+strconv.Itoa(appleY), 0)
//				socket.Send(strconv.Itoa(lobbyId)+",settings,moveTime,4", 0)
//				socket.Send(strconv.Itoa(lobbyId)+",settings,level,1", 0)
//			}
//
//			if action == "join" {
//			}
//
//			if action == "start" {
//				responses := handleStart(player, value)
//				for _, resMes := range responses {
//					socket.Send(resMes, 0)
//					fmt.Printf("sent: %s\n", resMes)
//				}
//			}
//
//			if action == "move" {
//				responses := handleMove(player, value)
//				for _, resMes := range responses {
//					socket.Send(resMes, 0)
//					fmt.Printf("sent: %s\n", resMes)
//				}
//			}
//		}
//	})
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		http.ServeFile(w, r, "websockets.html")
//	})
//
//	err := http.ListenAndServe(":8080", nil)
//	if err != nil {
//		println("ERROR: " + err.Error())
//		return
//	}
//}
//
//func (g *Game) reset() {
//	g.apple.X = 3 * gridSize
//	g.apple.Y = 3 * gridSize
//	g.moveTime = 4
//	g.score = 0
//	g.level = 1
//	g.moveDirection = dirNone
//
//	for i := 1; i <= 3; i++ {
//		startPosition := Position{
//			X: xGridCountInScreen / 2,
//			Y: (yGridCountInScreen / 2) + i,
//		}
//		game.players[i] = startPosition
//		//messages <- strconv.Itoa(i) + ",start," + strconv.Itoa(game.players[i].X) + "|" + strconv.Itoa(game.players[i].Y)
//	}
//}
//
//func handleCreate() int {
//	lobbies[idCounter] = Lobby{
//		Id: idCounter,
//		Game: Game{
//			players:       make(map[int]Position),
//			listOfPlayers: make([]int, 0),
//			moveDirection: dirNone,
//			apple:         Position{X: 3 * gridSize, Y: 3 * gridSize},
//			timer:         0,
//			moveTime:      4,
//			score:         0,
//			bestScore:     0,
//			level:         1,
//		},
//	}
//	return idCounter
//}
//
//func collidesWithApple() bool {
//	for i := 1; i <= len(game.listOfPlayers); i++ {
//		if game.players[i].X == game.apple.X && game.players[i].Y == game.apple.Y {
//			return true
//		}
//	}
//	return false
//}
//
//func (g *Game) collidesWithWall() bool {
//	for i := 1; i <= len(g.listOfPlayers); i++ {
//		if g.players[i].X < 0 ||
//			g.players[i].Y < 0 ||
//			g.players[i].X >= xGridCountInScreen ||
//			g.players[i].Y >= yGridCountInScreen {
//			return true
//		}
//	}
//	return false
//}
//
//func handleStart(player int, value string) []string {
//	messages := make([]string, 0)
//
//	stringSlice := strings.Split(value, "|")
//
//	if len(stringSlice) != 2 {
//		return make([]string, 0)
//	}
//
//	x, _ := strconv.Atoi(stringSlice[0])
//	y, _ := strconv.Atoi(stringSlice[1])
//
//	game.players[player] = Position{X: x, Y: y}
//	game.listOfPlayers = append(game.listOfPlayers, player)
//
//	appleMoveMessage := "apple," + "moved," + strconv.Itoa(game.apple.X) + "|" + strconv.Itoa(game.apple.Y)
//	messages = append(messages, appleMoveMessage)
//
//	playerStartMessage := strconv.Itoa(player) + "," + "started," + strconv.Itoa(x) + "|" + strconv.Itoa(y)
//	messages = append(messages, playerStartMessage)
//
//	return messages
//}
//
//func handleMove(player int, value string) []string {
//	messages := make([]string, 0)
//
//	switch value {
//	case "left":
//		pos := game.players[player]
//		pos.X--
//		game.players[player] = pos
//	case "right":
//		pos := game.players[player]
//		pos.X++
//		game.players[player] = pos
//	case "down":
//		pos := game.players[player]
//		pos.Y++
//		game.players[player] = pos
//	case "up":
//		pos := game.players[player]
//		pos.Y--
//		game.players[player] = pos
//	}
//
//	if collidesWithApple() {
//		game.apple.X = rand.Intn(xGridCountInScreen - 1)
//		game.apple.Y = rand.Intn(yGridCountInScreen - 1)
//		game.score++
//		game.level = 1
//		if game.bestScore < game.score {
//			game.bestScore = game.score
//		}
//
//		appleMoveMessage := "apple," + "moved," + strconv.Itoa(game.apple.X) + "|" + strconv.Itoa(game.apple.Y)
//		messages = append(messages, appleMoveMessage)
//	}
//
//	if game.collidesWithWall() {
//		game.reset()
//	}
//
//	playerMoveMessage := strconv.Itoa(player) + "," + "moved," + strconv.Itoa(game.players[player].X) + "|" + strconv.Itoa(game.players[player].Y)
//	messages = append(messages, playerMoveMessage)
//	return messages
//}
