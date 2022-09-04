package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	"wind/model"
	"wind/windmq"
)

var clientGame *Game

type Game struct {
	State    model.State
	Messages chan string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var apiAddr = flag.String("addr", "45.77.153.58:8080", "http service address")
var socketAddr = flag.String("socketAddr", "45.77.153.58:5556", "socket service address")

//var apiAddr = flag.String("addr", "localhost:8080", "http service address")
//var socketAddr = flag.String("socketAddr", "localhost:5556", "socket service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *apiAddr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	clientGame = &Game{
		State: model.State{
			Timer:    0,
			MoveTime: 4,
			Lag:      0.0,
			Player: model.Position{
				X:             17,
				Y:             30,
				MoveDirection: model.DirNone,
			},
		},
		Messages: make(chan string, 5000),
	}

	go func() {
		defer close(clientGame.Messages)
		println("Accepting messages to be written...")
		for {
			select {
			case <-done:
				return
			case t := <-clientGame.Messages:
				err := c.WriteMessage(websocket.TextMessage, []byte(t))
				if err != nil {
					log.Println("write:", err)
					return
				}
				log.Println("Wrote message: " + t)
			case <-interrupt:
				log.Println("interrupt")
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()

	subAddr, err := net.ResolveTCPAddr("tcp", *socketAddr)
	if err != nil {
		panic(err)
	}
	subscriber := windmq.NewSubscriber(subAddr, 1024)
	subscriber.Start()
	defer subscriber.Close()

	go func() {
		for {
			message := subscriber.EnsureReceived()
			handleMessageToClient(string(message))
		}
	}()

	startGame(clientGame)
}

func handleMessageToClient(msg string) {
	messageParts := strings.Split(msg, ",")

	entity := messageParts[0]
	action := messageParts[1]
	detail := messageParts[2]

	if entity == "player" {
		if action == "move" {
			coordinates := strings.Split(detail, "|")

			x, _ := strconv.Atoi(coordinates[0])
			y, _ := strconv.Atoi(coordinates[1])

			clientGame.State.Player.X = x
			clientGame.State.Player.Y = y

			println("Received message from server: " + msg)
		}
	}
}

func startGame(game *Game) {
	ebiten.SetWindowSize(model.ScreenWidth, model.ScreenHeight)
	ebiten.SetWindowTitle("Multiplayer Snake")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.State.Player.MoveDirection = model.DirLeft
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.State.Player.MoveDirection = model.DirRight
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.State.Player.MoveDirection = model.DirDown
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.State.Player.MoveDirection = model.DirUp
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State.Player.MoveDirection = model.DirNone
	}

	if g.needsToMoveSnake() {
		switch g.State.Player.MoveDirection {
		case model.DirNone:
			g.Messages <- "player,move,none"
		case model.DirLeft:
			pos := g.State.Player
			pos.X--
			g.State.Player = pos
			g.Messages <- "player,move,left"
		case model.DirRight:
			pos := g.State.Player
			pos.X++
			g.State.Player = pos
			g.Messages <- "player,move,right"
		case model.DirDown:
			pos := g.State.Player
			pos.Y++
			g.State.Player = pos
			g.Messages <- "player,move,down"
		case model.DirUp:
			pos := g.State.Player
			pos.Y--
			g.State.Player = pos
			g.Messages <- "player,move,up"
		}
	}

	g.State.Timer++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawPlayer(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Lag: %0.2f", g.State.Lag))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return model.ScreenWidth, model.ScreenHeight
}

func (g *Game) needsToMoveSnake() bool {
	return g.State.Timer%g.State.MoveTime == 0
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(g.State.Player.X*model.GridSize),
		float64(g.State.Player.Y*model.GridSize),
		model.GridSize,
		model.GridSize,
		color.RGBA{0xFF, 0xFF, 0x00, 0xff},
	)
}
