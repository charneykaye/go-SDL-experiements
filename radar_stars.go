/** Copyright 2015 Outright Mental, Inc. */

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"math"
	"os"
	"runtime"
)

var winWidth, winHeight int32 = 600, 600
var starRadius int32 = 2
var starBrightnessDecay float64 = 0.0012
var starBrightnessThreshold float64 = 0.05
var sweepDurationMs float64 = 10000
var numStars int = 10000
var centY, centX, maxR float64 = 300, 300, 300
var (
	twoPi float64 = math.Pi * 2
)

type Star struct {
	X int32
	Y int32
	B float64
}

func (s *Star) RenderToSurface(surface *sdl.Surface) {
	sBox := sdl.Rect{s.X - starRadius, s.Y - starRadius, starRadius * 2, starRadius * 2}
	//	var c = uint8(s.B * float64(255))
	// sColor := &sdl.Color{255, 255, 255, 255}
	surface.FillRect(&sBox, colorBrightness(s.B))
}

func (s *Star) Life() bool {
	s.B -= starBrightnessDecay
	return s.B > starBrightnessThreshold
}

func NewRadar() *Radar {
	r := &Radar{
		SweepPerTick: twoPi / sweepDurationMs,
	}
	r.Initialize()
	r.lastMs = sdl.GetTicks()
	return r
}

type Radar struct {
	SweepPerTick float64
	lastMs uint32
	winWidth int32
	NowMx float64
	NowMy float64
	NowSweep float64
	/* private */
	m_Stars []*Star
}

func (r *Radar) Initialize() {
	// Create stars
	for i := 0; i < numStars; i++ {
		s := &Star{}
		r.BirthStar(s)
		r.m_Stars = append(r.m_Stars, s)
	}
}

func (r *Radar) RenderToSurface(surface *sdl.Surface) {
	r.Life()
	for _, star := range r.m_Stars {
		if !star.Life() {
			r.BirthStar(star)
		}
		star.RenderToSurface(surface)
	}
}

func (r *Radar) Life() {
	nowMs := sdl.GetTicks()
	r.NowSweep += r.SweepPerTick * float64(nowMs - r.lastMs)
	r.lastMs = nowMs
	if r.NowSweep > twoPi {
		r.NowSweep -= twoPi
	}
	r.NowMy, r.NowMx = math.Sincos(r.NowSweep)
}

func (r *Radar) BirthStar(s *Star) {
	d := rand.Float64() * maxR
	s.X = int32(math.Max(0,math.Min(float64(winWidth), centX + d * r.NowMx)))
	s.Y = int32(math.Max(0,math.Min(float64(winHeight), centY + d * r.NowMy)))
	s.B = rand.Float64()
}

var palette = []uint32{
	0xFF000000,
	0xFF111111,
	0xFF222222,
	0xFF333333,
	0xFF444444,
	0xFF555555,
	0xFF666666,
	0xFF777777,
	0xFF888888,
	0xFF999999,
	0xFFAAAAAA,
	0xFFBBBBBB,
	0xFFCCCCCC,
	0xFFDDDDDD,
	0xFFEEEEEE,
	0xFFFFFFFF,
}

func colorBrightness(b float64) uint32 {
	return palette[int(b*float64(15))]
}

func main() {
	runtime.LockOSThread()
	game := NewGame()
	os.Exit(game.Start())
}

func NewGame() *Game {
	g := &Game{
		Name: "stars",
	}
	g.Initialize()
	return g
}

type Game struct {
	/* public */
	Name string
	/* private objects */
	m_Radar *Radar
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

	g.m_Radar = NewRadar()

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

	g.m_Radar.RenderToSurface(g.sdlScreenSurface)

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

type StateEnum uint

const (
	STATE_LOADING StateEnum = 3
	STATE_PLAYING StateEnum = 5
	// it can be assumed that all alive states are < STATE_FINISHED
	STATE_FINISHED StateEnum = 6
	STATE_FAILED   StateEnum = 7
)
