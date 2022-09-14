package main

import (
	"flag"
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

var requests map[string]string
var latestClientRequest = 0
var clientGame *Game

type Game struct {
	State    model.State
	Messages chan string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var turnOnServerReconciliation = true

//var isServerReconciliationTurnedOn = flag.Bool("isServerReconciliationOn", true, "turn on server reconciliation")
//var isServerReconciliationTurnedOn = flag.Bool("isServerReconciliationOn", false, "turn off server reconciliation")

//var apiAddr = flag.String("addr", "45.77.153.58:8080", "http service address")
//var socketAddr = flag.String("socketAddr", "45.77.153.58:5556", "socket service address")

//var apiAddr = flag.String("addr", "67.219.107.162:8080", "http service address")
//var socketAddr = flag.String("socketAddr", "67.219.107.162:5556", "socket service address")

var apiAddr = flag.String("addr", "localhost:8080", "http service address")
var socketAddr = flag.String("socketAddr", "localhost:5556", "socket service address")

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

	requests = make(map[string]string, 0)
	clientGame = &Game{
		State: model.State{
			Timer:    0,
			MoveTime: 4,
			Lag:      0.0,
			Player: model.Entity{
				X:             38,
				Y:             15,
				MoveDirection: model.DirNone,
				Type:          model.Player,
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
				log.Println(t)
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
			handleGameStateUpdate(string(message))
		}
	}()

	startGame(clientGame)
	clientGame.sendJoinMessage()
}

func performServerReconciliation(latestServerRequest int, xStartPoint int, yStartPoint int) {
	if latestServerRequest < latestClientRequest {
		var xTotalDelta = 0
		var yTotalDelta = 0
		for i := latestServerRequest + 1; i <= latestClientRequest; i++ {
			// TODO: Can cause a concurrent map read and map write...
			// 		"fatal error: concurrent map read and map write"
			// TODO: Have to set a mutex or use a WriteLock
			messageAdjustment := requests[strconv.Itoa(i)]

			//handleMessageAdjustment(messageAdjustment)
			xTotalDelta, yTotalDelta = handleMessageAdjustment2(messageAdjustment, xTotalDelta, yTotalDelta)
		}

		if xTotalDelta == 0 || yTotalDelta == 0 {
			println("WHY U 0?")
		}

		// Apply delta to the start points
		xFinishPoint := xStartPoint + xTotalDelta
		yFinishPoint := yStartPoint + yTotalDelta

		// Apply the total deltas on the square
		pos := clientGame.State.Player
		pos.X = xFinishPoint
		pos.Y = yFinishPoint
		clientGame.State.Player = pos
		println(strconv.Itoa(latestServerRequest) + "  < " + strconv.Itoa(latestClientRequest))
		println("Went from (" + strconv.Itoa(xStartPoint) + ", " + strconv.Itoa(yStartPoint) + ") " +
			"to (" + strconv.Itoa(xFinishPoint) + "  < " + strconv.Itoa(yFinishPoint))
	}
}

/*
Unlike v1, I am only going to apply adjustments to the square until ALL
adjustments have been calculated together. This should help prevent the square
from being visually glitchy when changing directions (especially in 90 degree angles)
*/
func handleMessageAdjustment2(msg string, xDelta int, yDelta int) (int, int) {
	messageParts := strings.Split(msg, ",")

	if len(messageParts) != 4 {
		println("ERROR!!!")
		return 0, 0
	}

	_ = messageParts[0]
	entity := messageParts[1]
	action := messageParts[2]
	detail := messageParts[3]

	if entity == "player" {
		if action == "move" {
			switch detail {
			case "none":
			case "left":
				xDelta = xDelta - 1
			case "right":
				xDelta = xDelta + 1
			case "down":
				yDelta = yDelta + 1
			case "up":
				yDelta = yDelta - 1
			}
		}
	}

	return xDelta, yDelta
}

func handleGameStateUpdate(msg string) {
	println(msg)
	messageParts := strings.Split(msg, ",")

	if len(messageParts) != 4 {
		if len(messageParts) == 2 {
			settingEntity := messageParts[0]
			settingDetail := messageParts[1]

			if settingEntity == "player" {

				coordinates := strings.Split(settingDetail, "|")

				settingX, _ := strconv.Atoi(coordinates[0])
				settingY, _ := strconv.Atoi(coordinates[1])

				clientGame.State.Player.X = settingX
				clientGame.State.Player.Y = settingY
			}
		}
		println(msg)
		return
	}

	count := messageParts[0]
	entity := messageParts[1]
	action := messageParts[2]
	detail := messageParts[3]

	ct, _ := strconv.Atoi(count)

	if entity == "player" {
		if action == "move" {
			coordinates := strings.Split(detail, "|")

			x, _ := strconv.Atoi(coordinates[0])
			y, _ := strconv.Atoi(coordinates[1])

			if latestClientRequest == 50 {
				clientGame.State.Player.X = x
				clientGame.State.Player.Y = y
			}
			clientGame.State.Follow.X = x
			clientGame.State.Follow.Y = y

			xStartPoint := x
			yStartPoint := y

			if turnOnServerReconciliation {
				performServerReconciliation(ct, xStartPoint, yStartPoint)
			} else {
				clientGame.State.Player.X = x
				clientGame.State.Player.Y = y
			}
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
		latestClientRequest++
		ct := strconv.Itoa(latestClientRequest)
		switch g.State.Player.MoveDirection {
		case model.DirNone:
			g.createAndSendMessage(ct, "none")
		case model.DirLeft:
			x := clientGame.State.Player.X
			clientGame.State.Player.X = x - 1
			g.createAndSendMessage(ct, "left")
		case model.DirRight:
			x := clientGame.State.Player.X
			clientGame.State.Player.X = x + 1
			g.createAndSendMessage(ct, "right")
		case model.DirDown:
			y := clientGame.State.Player.Y
			clientGame.State.Player.Y = y + 1
			g.createAndSendMessage(ct, "down")
		case model.DirUp:
			y := clientGame.State.Player.Y
			clientGame.State.Player.Y = y - 1
			g.createAndSendMessage(ct, "up")
		}
	}

	g.State.Timer++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// g.drawLines(screen)
	g.drawPlayer(screen)
	g.drawFollow(screen)

	if turnOnServerReconciliation {
		ebitenutil.DebugPrint(screen, "Server Reconciliation Turned On")
	} else {
		ebitenutil.DebugPrint(screen, "Server Reconciliation Turned Off")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return model.ScreenWidth, model.ScreenHeight
}

func (g *Game) needsToMoveSnake() bool {
	return g.State.Timer%g.State.MoveTime == 0
}

func (g *Game) drawLines(screen *ebiten.Image) {
	for i := 0; i < model.ScreenHeight; i += 2 {
		ebitenutil.DrawRect(
			screen,
			0,
			float64(i*model.GridSize),
			model.ScreenWidth,
			model.GridSize,
			color.RGBA{0, 0, 0xff, 0xff},
		)
	}
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(g.State.Player.X*model.GridSize),
		float64(g.State.Player.Y*model.GridSize),
		model.GridSize,
		model.GridSize,
		color.RGBA{0xff, 0, 0, 0xff},
	)
}

func (g *Game) drawFollow(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(g.State.Follow.X*model.GridSize),
		float64(g.State.Follow.Y*model.GridSize),
		model.GridSize,
		model.GridSize,
		color.RGBA{0, 0xff, 0, 0xff},
	)
}

func (g *Game) createAndSendMessage(count string, direction string) {
	outboundMessage := count + ",player,move," + direction
	requests[count] = outboundMessage
	g.Messages <- outboundMessage
}

func (g *Game) sendJoinMessage() {
	g.Messages <- "player,join"
}
