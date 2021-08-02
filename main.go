package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
)

type Game struct {
	Player Player
	World  World
	Window Window
}

type Window struct {
	width  int
	height int
}

type World struct {
	BorderGround BorderGround
	Gravity      float64
	Drag         float64
	GroundDrag   float64
	Ground       float64
}

type BorderGround struct {
	TopY    float64
	BottomY float64
	LeftX   float64
	RightX  float64
}

type Player struct {
	Img       image.Image
	Name      string
	PosX      float64
	PosY      float64
	Size      image.Point
	DeltaX    float64
	DeltaY    float64
	OnGround  bool
	JumpPower float64
	MoveSpeed float64
}

func main() {
	game, err := NewGame()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ebiten.RunGame(game)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func NewGame() (*Game, error) {
	window := Window{width: 800, height: 600}
	ebiten.SetWindowResizable(false)
	ebiten.SetWindowSize(window.width, window.height)
	playerSprite, err := ReadImage("resources/mario_48x.png")
	if err != nil {
		return nil, err
	}
	player := Player{
		Img:       playerSprite,
		Name:      "Mario",
		PosX:      float64((window.width - 48) / 2),
		PosY:      float64(window.height - 48*2),
		Size:      playerSprite.Bounds().Size(),
		DeltaX:    0,
		DeltaY:    0,
		OnGround:  true,
		JumpPower: -14,
		MoveSpeed: 6,
	}
	world := World{
		BorderGround: BorderGround{
			TopY:    0 + 48,
			BottomY: float64(window.height - 48),
			LeftX:   0 + 48,
			RightX:  float64(window.width - 48),
		},
		Gravity:    0.9,
		Drag:       1.0,
		GroundDrag: 0.9,
	}
	game := &Game{Player: player, World: world, Window: window}

	return game, nil
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.Player.DeltaX = -g.Player.MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.Player.DeltaX = g.Player.MoveSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if g.Player.OnGround {
			g.Player.DeltaY = g.Player.JumpPower
		}
	}

	g.Player.DeltaY += g.World.Gravity
	g.Player.DeltaY *= g.World.Drag
	if g.Player.OnGround {
		g.Player.DeltaX *= g.World.GroundDrag
	} else {
		g.Player.DeltaX *= g.World.Drag
	}
	g.Player.PosX += g.Player.DeltaX
	g.Player.PosY += g.Player.DeltaY

	if g.Player.PosY+float64(g.Player.Size.Y) >= g.World.BorderGround.BottomY {
		g.Player.PosY = g.World.BorderGround.BottomY - float64(g.Player.Size.Y)
		g.Player.DeltaY = 0
		g.Player.OnGround = true
	} else {
		g.Player.OnGround = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Azure)
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(g.Player.PosX, g.Player.PosY)
	screen.DrawImage(g.Player.Img.(*ebiten.Image), options)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func ReadImage(filename string) (image.Image, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
