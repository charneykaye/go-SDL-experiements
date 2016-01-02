/** Author: Charney Kaye */

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"math"
	"os"
	"runtime"
)

var (
	fireWidth, fireHeight int = 300, 300
	firePointSize int = 2
	fireDecay float64 = 0.99
)

/* the raster is in a
███████╗██╗██████╗ ███████╗
██╔════╝██║██╔══██╗██╔════╝
█████╗  ██║██████╔╝█████╗
██╔══╝  ██║██╔══██╗██╔══╝
██║     ██║██║  ██║███████╗
╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝*/

func NewFire() *Fire {
	r := &Fire{}
	r.Initialize()
	return r
}

type Fire struct {
	/* private */
	Points [][]float64
}

func (r *Fire) Initialize() {
	// Allocate raster of fire points
	r.Points = make([][]float64, fireHeight)
	for y := 0; y < fireHeight; y++ {
		r.Points[y] = make([]float64, fireWidth)
	}
}

func (r *Fire) RenderToSurface(surface *sdl.Surface) {
	sBox := sdl.Rect{0, 0, int32(firePointSize), int32(firePointSize)}
	for y := 0; y < fireLimitY; y++ {
		sBox.Y = int32(y * firePointSize)
		for x := 0; x < fireWidth; x++ {
			r.PointLife(y, x)
			sBox.X = int32(x * firePointSize)
			surface.FillRect(&sBox, colorBrightness(r.Points[y][x]))
		}
	}
	for x := 0; x < fireWidth; x++ {
		r.PointBirth(fireLimitY, x)
	}
}

func(r *Fire) PointLife(y int, x int) {
	// each row inherits from higher row
	if x == 0 {
		// lower x-limit
		r.Points[y][x] = fireDecay * (r.Points[y+1][x] + r.Points[y+1][x + 1]) / 2
	} else if x == fireLimitX {
		// upper x-limit
		r.Points[y][x] = fireDecay * (r.Points[y+1][x] + r.Points[y+1][x - 1]) / 2
	} else {
		// all x between
		r.Points[y][x] = fireDecay * (r.Points[y+1][x - 1] + r.Points[y+1][x] + r.Points[y+1][x + 1]) / 3
	}
}

func(r *Fire) PointBirth(y int, x int) {
	// bottom row generates pixels that are on/off
	// chance of being on (c) is inversely proportional to distance from center
	if rand.Float64() < 1 - math.Abs(float64(x - fireCenterX)) / float64(fireCenterX) {
		r.Points[y][x] = 1
	} else {
		r.Points[y][x] = 0
	}
}

/* there is one fire for the whole
 ██████╗  █████╗ ███╗   ███╗███████╗
██╔════╝ ██╔══██╗████╗ ████║██╔════╝
██║  ███╗███████║██╔████╔██║█████╗
██║   ██║██╔══██║██║╚██╔╝██║██╔══╝
╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗
 ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝*/

func NewGame() *Game {
	g := &Game{
		Name: "fire",
	}
	g.Initialize()
	return g
}

type Game struct {
	/* public */
	Name string
	/* private objects */
	m_Fire *Fire
	/* private */
	m_State StateEnum
	nowMs   uint32
	lastMs  uint32
	/* private: SDL */
	sdlRenderer      *sdl.Renderer
	sdlScreenSurface *sdl.Surface
	sdlScreenTexture *sdl.Texture
	sdlStarColor     sdl.Color
	sdlBgColor       sdl.Color
	sdlWindow        *sdl.Window
}

func (g *Game) Initialize() {
	var err error
	log.SetLevel(log.DebugLevel)

	log.WithFields(log.Fields{
		"name": g.Name,
	}).Info("Initialize Game")

	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to Init Simple DirectX Layer")
	}

	g.sdlStarColor = sdl.Color{255, 255, 255, 255}
	g.sdlBgColor = sdl.Color{0, 0, 0, 255}

	g.sdlWindow, err = sdl.CreateWindow(
		g.Name,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int(winWidth), int(winHeight),
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create window")
	}

	g.sdlRenderer, err = sdl.CreateRenderer(g.sdlWindow, -1,
		sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create renderer")
	}

	g.sdlScreenSurface, err = sdl.CreateRGBSurface(0, int32(winWidth), int32(winHeight), int32(32), 0x00FF0000, 0x0000FF00, 0x000000FF, 0xFF000000)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create screen surface")
	}

	g.m_Fire = NewFire()

	g.ChangeState(STATE_LOADING)
}

func (g *Game) Start() int {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"recover": r,
			}).Warn("Game Recovered")
		}
		g.Teardown()
	}()

	g.ChangeState(STATE_PLAYING)
	for g.Alive() {
		if g.NowMs() {
			g.PollEvents()
			g.Render()
		}
	}
	return 0
}

func (g *Game) Render() {
	var err error

	g.sdlScreenSurface.FillRect(nil, 0xFF000000)

	g.m_Fire.RenderToSurface(g.sdlScreenSurface)

	g.sdlScreenTexture, err = g.sdlRenderer.CreateTextureFromSurface(g.sdlScreenSurface)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Could not create texture from surface")
	}
	defer g.sdlScreenTexture.Destroy()

	g.sdlRenderer.Copy(g.sdlScreenTexture, nil, nil)

	g.sdlRenderer.Present()
}

func (g *Game) Stop() {
	g.ChangeState(STATE_FINISHED)
}

func (g *Game) Teardown() {
	log.Info("Teardown Game")
	g.sdlRenderer.Destroy()
	g.sdlWindow.Destroy()
}

func (g *Game) ChangeState(s StateEnum) {
	g.m_State = s
	log.WithFields(log.Fields{
		"state": g.StateName(),
	}).Info("Game changed")
	switch g.m_State {
	case STATE_LOADING:
	case STATE_PLAYING:
	case STATE_FINISHED:
	case STATE_FAILED:
	}
}

func (g *Game) StateName() string {
	switch g.m_State {
	case STATE_LOADING:
		return "Loading"
	case STATE_PLAYING:
		return "Playing"
	case STATE_FINISHED:
		return "Finished"
	case STATE_FAILED:
		return "Failed"
	}
	return ""
}

func (g *Game) PollEvents() {
	var e sdl.Event
	for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch t := e.(type) {
		case *sdl.QuitEvent:
			g.Stop()
		case *sdl.KeyUpEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				g.Stop()
			}
		}
	}
}

func (g *Game) Alive() bool {
	return g.m_State < STATE_FINISHED
}

func (g *Game) NowMs() bool {
	g.nowMs = sdl.GetTicks()
	if g.nowMs != g.lastMs {
		g.lastMs = g.nowMs
		return true
	}
	return false
}

/* the game is instantiated from
███╗   ███╗ █████╗ ██╗███╗   ██╗
████╗ ████║██╔══██╗██║████╗  ██║
██╔████╔██║███████║██║██╔██╗ ██║
██║╚██╔╝██║██╔══██║██║██║╚██╗██║
██║ ╚═╝ ██║██║  ██║██║██║ ╚████║
╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝*/

func main() {
	runtime.LockOSThread()
	game := NewGame()
	os.Exit(game.Start())
}

type StateEnum uint

const (
	STATE_LOADING StateEnum = 3
	STATE_PLAYING StateEnum = 5
	// it can be assumed that all alive states are < STATE_FINISHED
	STATE_FINISHED StateEnum = 6
	STATE_FAILED   StateEnum = 7
)

var (
	winWidth = fireWidth * firePointSize
	winHeight = fireHeight * firePointSize
	fireLimitX = fireWidth - 1
	fireLimitY = fireHeight - 1
	fireCenterX = fireWidth / 2
)

var palette = []uint32{
	0xFF000000,
	0xFF25120c,
	0xFF3b1b06,
	0xFF5a2211,
	0xFF722400,
	0xFF9c3e0a,
	0xFFbc490a,
	0xFFe16205,
	0xFFf4700b,
	0xFFfc8409,
	0xFFff9315,
	0xFFffb234,
	0xFFffe14f,
	0xFFffff53,
	0xFFfffeab,
	0xFFe16205,
}

func colorBrightness(b float64) uint32 {
	return palette[int(b*float64(15))]
}


