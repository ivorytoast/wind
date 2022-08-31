package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	zmq "github.com/pebbe/zmq4"
	"image/color"
	"log"
	"math/rand"
	"time"
	"wind/model"
	"wind/queue"
)

const (
	dirNone = iota
	dirLeft
	dirRight
	dirDown
	dirUp
)

/*
Game has two responsibilities:
	1. Drawing to the screen
	2. Sending updates to Wind

It shares the same state as servers, and only knows the most basic state of all
*/
var windQueue *zmq.Socket

type Game struct {
	State model.State
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if g.State.MoveDirection != dirRight {
			g.State.MoveDirection = dirLeft
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if g.State.MoveDirection != dirLeft {
			g.State.MoveDirection = dirRight
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.State.MoveDirection != dirUp {
			g.State.MoveDirection = dirDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.State.MoveDirection != dirDown {
			g.State.MoveDirection = dirUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State.MoveDirection = dirNone
	}

	if g.allowSnakeDirectionToChange() {
		switch g.State.MoveDirection {
		case dirLeft:
			//messages <- *player + ",move,left"
		case dirRight:
			//messages <- *player + ",move,right"
		case dirDown:
			//messages <- *player + ",move,down"
		case dirUp:
			//messages <- *player + ",move,up"
		}
	}

	g.State.Timer++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawApple(screen)

	for playerId := 1; playerId <= len(g.State.Players); playerId++ {
		g.drawPlayer(screen, playerId)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f Level: %d Score: %d Best Score: %d", ebiten.CurrentFPS(), g.State.Level, g.State.Score, g.State.BestScore))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return model.ScreenWidth, model.ScreenHeight
}

func (g *Game) drawApple(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(g.State.Apple.X*model.GridSize),
		float64(g.State.Apple.Y*model.GridSize),
		model.GridSize, model.GridSize,
		color.RGBA{0xFF, 0x00, 0x00, 0xff},
	)
}

func (g *Game) drawPlayer(screen *ebiten.Image, playerId int) {
	// TODO: Can throw a fatal exception
	ebitenutil.DrawRect(
		screen,
		float64(g.State.Players[playerId].X*model.GridSize),
		float64(g.State.Players[playerId].Y*model.GridSize),
		model.GridSize,
		model.GridSize,
		color.RGBA{0xFF, 0xFF, 0x00, 0xff},
	)
}

func (g *Game) allowSnakeDirectionToChange() bool {
	return true
	//return g.State.Timer%g.State.MoveTime == 0
}

func main() {
	publisher := queue.CreatePublisher("tcp://*:5556")
	defer publisher.Close()
	go func() {
		for {
			message := queue.CreateMessage(
				"entity",
				"action",
				"detail",
			)
			publisher.Send(string(model.Client), zmq.SNDMORE)
			publisher.Send("clientID", zmq.SNDMORE)
			publisher.Send("clientLocation", zmq.SNDMORE)
			publisher.Send(message, 0)
			println("Sent message...")
			time.Sleep(1 * time.Second)
		}
	}()

	subscriber := queue.CreateSubscriber("tcp://localhost:5556")
	defer subscriber.Close()
	go func() {
		_ = subscriber.SetSubscribe("server")
		for {
			msg, _ := subscriber.Recv(0)
			println("[Received Message]: " + msg)
		}
	}()

	ebiten.SetWindowSize(model.ScreenWidth, model.ScreenHeight)
	ebiten.SetWindowTitle("Multiplayer Snake")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
