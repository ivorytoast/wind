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
	"time"
	"wind/model"
	"wind/windmq"
)

var subscriber *windmq.Subscriber

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
//var socketAddr = flag.String("socketAddr", "45.77.153.58:5556", "socket service address")

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

	subAddr, err := net.ResolveTCPAddr("tcp", "45.77.153.58:5556")
	if err != nil {
		panic(err)
	}
	subscriber = windmq.NewSubscriber(subAddr, 1024)
	subscriber.Start()
	defer subscriber.Close()

	go func() {
		for {
			message, _ := subscriber.Receive()
			fmt.Println("Received message: [" + string(message) + "]")
		}
	}()

	done := make(chan struct{})

	game := &Game{
		State: model.State{
			Lag: 0.0,
			Player: model.Position{
				X:             17,
				Y:             30,
				MoveDirection: model.DirNone,
			},
		},
		Messages: make(chan string, 5000),
	}

	go func() {
		defer close(game.Messages)
		println("Accepting messages to be written...")
		for {
			select {
			case <-done:
				return
			case t := <-game.Messages:
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

	startGame(game)
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
		g.Messages <- "player,move,left"
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.State.Player.MoveDirection = model.DirRight
		g.Messages <- "player,move,right"
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.State.Player.MoveDirection = model.DirDown
		g.Messages <- "player,move,down"
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.State.Player.MoveDirection = model.DirUp
		g.Messages <- "player,move,up"
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State.Player.MoveDirection = model.DirNone
		g.Messages <- "player,move,none"
	}

	switch g.State.Player.MoveDirection {
	case model.DirNone:
	case model.DirLeft:
		pos := g.State.Player
		pos.X--
		g.State.Player = pos
	case model.DirRight:
		pos := g.State.Player
		pos.X++
		g.State.Player = pos
	case model.DirDown:
		pos := g.State.Player
		pos.Y++
		g.State.Player = pos
	case model.DirUp:
		pos := g.State.Player
		pos.Y--
		g.State.Player = pos
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawPlayer(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Lag: %0.2f", g.State.Lag))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return model.ScreenWidth, model.ScreenHeight
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
