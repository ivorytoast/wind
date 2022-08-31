package main

//
//import (
//	"flag"
//	zmq "github.com/pebbe/zmq4"
//	"log"
//	"net/url"
//	"os"
//	"os/signal"
//	"strconv"
//	"strings"
//	"time"
//	"wind/model"
//
//	"github.com/gorilla/websocket"
//)
//
///*
//Wind has servers, which has games, which has players...
//
//This is a load balancer between the clients and the servers
//
//Client sends a request to server
//Server sends a request to client
//*/
//
//type Server struct {
//	ID      int
//	Clients []Game
//}
//
//var servers = make([]Server, 10)
//var addr = flag.String("addr", "localhost:8080", "http service address")
//var player = flag.String("player", "1", "player")
//var messages = make(chan string, 5000)
//
//const (
//	xGridCountInScreen = model.ScreenWidth / model.GridSize
//	yGridCountInScreen = model.ScreenHeight / model.GridSize
//)
//
////func CreateNewServer() int {
////	serverId := rand.Intn(10000000)
////	newServer := Server{
////		ID:      serverId,
////		Clients: make([]Game, 0),
////	}
////	servers = append(servers, newServer)
////	return serverId
////}
//
////func receiveStart(player int, value string) string {
////	stringSlice := strings.Split(value, "|")
////
////	if len(stringSlice) != 2 {
////		return "ERROR (START)"
////	}
////
////	x, _ := strconv.Atoi(stringSlice[0])
////	y, _ := strconv.Atoi(stringSlice[1])
////
////	onePos := model.Position{
////		X: x,
////		Y: y,
////	}
////	localGame.players[player] = onePos
////
////	return strconv.Itoa(player) + "," + "started at" + strconv.Itoa(localGame.players[player].X) + "|" + strconv.Itoa(localGame.players[player].Y)
////}
//
////func receiveMove(player int, value string) string {
////	stringSlice := strings.Split(value, "|")
////
////	if len(stringSlice) != 2 {
////		return "ERROR (MOVE)"
////	}
////
////	x, _ := strconv.Atoi(stringSlice[0])
////	y, _ := strconv.Atoi(stringSlice[1])
////
////	onePos := Position{
////		X: x,
////		Y: y,
////	}
////	localGame.players[player] = onePos
////
////	return strconv.Itoa(player) + "," + "moved to" + strconv.Itoa(localGame.players[player].X) + "|" + strconv.Itoa(localGame.players[player].Y)
////}
//
///*
//Acts as a load balancer between the servers and the clients
//
//The messages it receives must be forwarded to the correct server and games
//You could create more of these load balancers by subscribing to only certain filters
//
//Therefore, no logic is created here. Just forwarding and compressing any messages
//to be as efficient as possible
//
//The message should be:
//	{source}|{location}|{entity}|{action}|{detail}
//	10|10|20|20|20
//*/
//func main() {
//
//	// -- Connecting to the game engine -- //
//	flag.Parse()
//	log.SetFlags(0)
//
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt)
//
//	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
//	log.Printf("connecting to %s", u.String())
//
//	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//	if err != nil {
//		log.Fatal("dial:", err)
//	}
//	defer c.Close()
//	// -- -- -- -- -- -- -- -- -- -- -- -- //
//
//	go func() {
//		context, _ := zmq.NewContext()
//		socket, _ := context.NewSocket(zmq.SUB)
//		defer socket.Close()
//
//		socket.SetSubscribe("")
//		socket.Connect("tcp://localhost:5556")
//
//		// Actually receiving the messages and sending them off to the game
//		for {
//			message, _ := socket.Recv(0)
//
//			// {sourceType},{sourceID},{location},{entity},{action},{detail} (Phase 1)
//			//  10|10|20|20|20 (Phase 2)
//
//			messageParts := strings.Split(message, ",")
//			if len(messageParts) != 5 {
//				println("error: invalid message length: [" + message + "]")
//				continue
//			}
//
//			/*
//
//				Wind
//					[]Server (Backend)
//						[]Games (Clients -- Xbox Connecting)
//							[]Players
//							[]Entities
//
//				Wind knows only about servers. All other info is encapsulated by the servers
//				You can start multiple Wind instances -- i.e. different locations around the world
//				Wind is simply a load balancer -- therefore, it can handle more load than Server or Game
//				You can start up multiple servers independent of one another
//				You can start up multiple games within servers independent of one another
//
//				All Traffic Goes Through Wind "Load-Balancers"
//					This allows for large scalability on both servers and games
//				Since Wind is simply parsing messages, it can handle the highest load
//				Since Server
//
//				Wind NY
//					Server 1
//						Client 1
//					Server 2
//						Client 1
//						Client 2
//
//				Wind NA
//					Server NY
//						Client 1
//						Client 2
//						Client 3
//					Server TX
//						Client 1
//						Client 2
//					Server CA
//						Client 1
//						Client 2
//						Client 3
//						Client 4
//
//				Wind ASIA
//					Server X
//						Client 1
//						Client 2
//						Client 3
//					Server Y
//						Client 1
//						Client 2
//					Server Z
//						Client 1
//						Client 2
//						Client 3
//						Client 4
//
//					Server 1 (123)
//						Game 1 (abc)
//							Player 1 (1A)
//							Player 2 (2B)
//							Apple (AA)
//							Health (BB)
//						Game 2 (def)
//							Player 1 (3C)
//							Player 2 (4D)
//							Apple (CC)
//					Server 2 (456)
//						Game 3 (ghi)
//							Player 1 (1E)
//							Player 2 (2F)
//						Game 4 (jkl)
//							Player 1 (1G)
//							Player 2 (2H)
//
//					Game 1 sends update message to server
//						client | abc | entity | action | detail
//
//					Find which server contains (abc). Server (123) contains (abc)
//						Location = 123
//
//					Since a client can only send to a server, no further location is required.
//					Server will be given new information to update game state
//
//					---
//
//					Server 2 sends update message to Game 3 (ghi)
//						client | abc | entity | action | detail
//
//					Find which server contains (abc). Server (123) contains (abc)
//						Location = 123
//
//					Since a client can only send to a server, no further location is required.
//					Server will be given new information to update game state
//			*/
//
//			sourceType := messageParts[0]
//			sourceID := messageParts[1]
//			entity := messageParts[3]
//			action := messageParts[4]
//			detail := messageParts[5]
//
//			if sourceType == "client" {
//				// For location, a client knows what server it is located on. It can only
//				// send messages to that specific server
//
//				location := getLocation(sourceType, sourceID)
//
//				if action == "join" {
//					servers[location].Clients = append(servers[location].Clients, Game{
//						ID:            0,
//						ServerID:      0,
//						players:       nil,
//						listOfPlayers: nil,
//						moveDirection: 0,
//						apple:         Position{},
//						timer:         0,
//						moveTime:      0,
//						score:         0,
//						bestScore:     0,
//						level:         0,
//					})
//				}
//			} else if sourceType == "server" {
//				// For location, a server knows all its clients. So can choose to hit all
//				// or select individually
//
//				location := getLocation(sourceType, sourceID)
//			}
//
//			if stringSlice[0] == "apple" {
//				println("APPLE received: " + stringSlice[1] + ", " + stringSlice[2])
//				appleSlice := strings.Split(stringSlice[2], "|")
//
//				if len(appleSlice) != 2 {
//					println("Invalid apple message. The length was: " + strconv.Itoa(len(appleSlice)))
//					continue
//				}
//
//				x, _ := strconv.Atoi(appleSlice[0])
//				y, _ := strconv.Atoi(appleSlice[1])
//				localGame.apple.X = x
//				localGame.apple.Y = y
//				continue
//			}
//			player, _ := strconv.Atoi(stringSlice[0])
//			key := stringSlice[1]
//			value := stringSlice[2]
//
//			/*
//				Join (Join a game)
//				Start (Game begins)
//				Move (Player moves the snake)
//				Collide (Snake collides with an object)
//				Win (Player wins)
//				Lose (Player loses)
//
//				State (The current state of the game is given)
//			*/
//			switch key {
//			case "joined":
//				//receivedJoined(player, value)
//			case "started":
//				receiveStart(player, value)
//			case "moved":
//				receiveMove(player, value)
//			case "collide":
//				//receiveCollide(player, value)
//			case "win":
//				//receiveWin(player, value)
//			case "lose":
//				//receiveLose(player, value)
//			case "state":
//				//receiveState(player, value)
//			}
//
//			log.Printf("recv: %s", message)
//		}
//	}()
//
//	done := make(chan struct{})
//
//	defer close(messages)
//
//	go func() {
//		for {
//			select {
//			case <-done:
//				return
//			case t := <-messages:
//				err := c.WriteMessage(websocket.TextMessage, []byte(t))
//				if err != nil {
//					log.Println("write:", err)
//					return
//				}
//				log.Println("Wrote message: " + t)
//			case <-interrupt:
//				log.Println("interrupt")
//
//				// Cleanly close the connection by sending a close message and then
//				// waiting (with timeout) for the server to close the connection.
//				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//				if err != nil {
//					log.Println("write close:", err)
//					return
//				}
//				select {
//				case <-done:
//				case <-time.After(time.Second):
//				}
//				return
//			}
//		}
//	}()
//}
