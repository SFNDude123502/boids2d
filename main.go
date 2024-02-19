package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	detectRadius     float64 = 100.0
	personalSpace    float64 = 30.0
	cohesionFactor   float64 = 0.005
	avoidanceFactor  float64 = 0.1
	allignmentFactor float64 = 0.05
	flockSize        int     = 100
	speedLimit       float64 = 15.0
	turnSpeed        float64 = 0.5
	margin           float64 = 300
)

var (
	//go:embed Boid2d.png
	boidBytes  []byte
	boidSprite *ebiten.Image
	backgnd    *ebiten.Image
)

type Game struct {
	Flock []*Boid
	op    ebiten.DrawImageOptions
}

func (g *Game) init() {
	g.Flock = make([]*Boid, 0)
	for range make([]int, flockSize) {
		var newBoid = &Boid{Loc: Vector{0, 0}, V: Vector{0, 0}}
		newBoid.Loc.X = float64(rand.Intn(1920-200) + 100) //avoiding placing boids near a wall
		newBoid.Loc.Y = float64(rand.Intn(1080-200) + 100)
		newBoid.V.X = float64(rand.Float64()*30 - 15)
		newBoid.V.Y = float64(rand.Float64()*30 - 15)
		g.Flock = append(g.Flock, newBoid)
	}
}

func (g *Game) Update() error {
	for _, b := range g.Flock {
		b.V.add(b.Cohesion(g.Flock))
		b.V.add(b.Seperation(g.Flock))
		b.V.add(b.Allignment(g.Flock))
		b.V.add(b.Bounds())

		b.EnforceSpeedLimit()

		(&b.Loc).add(b.V)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgnd, nil)
	w, h := backgnd.Bounds().Dx(), backgnd.Bounds().Dy()

	for _, b := range g.Flock {
		g.op.GeoM.Reset()

		g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)      //idk why i need to translate before rotating,
		g.op.GeoM.Rotate(math.Pi/2 + math.Atan2(b.V.Y, b.V.X)) // but it breaks if i dont
		g.op.GeoM.Translate(float64(w)/2, float64(h)/2)

		g.op.GeoM.Scale(
			float64(w)/(60*480.0), // sets X-dimention to 1/60 of the sceen width
			float64(h)/(30*480.0), // sets Y-dimention to 1/30 if the screen height
		)
		g.op.GeoM.Translate(b.Loc.X, b.Loc.Y)
		screen.DrawImage(boidSprite, &g.op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return backgnd.Bounds().Dx(), backgnd.Bounds().Dy()
}

func init() {
	backgnd = ebiten.NewImage(1920, 1080)
	backgnd.Fill(color.RGBA{63, 91, 112, 255})

	boid, _, err := image.Decode(bytes.NewReader(boidBytes))
	if err != nil {
		panic(err)
	}
	boidSprite = ebiten.NewImageFromImage(boid)
}

func main() {
	game := &Game{}
	game.init()
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Boid Simulation")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
