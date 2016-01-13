/** Author: Charney Kaye */

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	// "math/rand"
	"os"
	"runtime"
)

var (
	graphWidth, graphHeight int     = 300, 300
	graphPointSize         int     = 2
	graphGenRows           int     = 2
	graphDecay             float64 = 0.98
)

const (
	GOLDEN_RATIO float64 = 1.61803398875
	INVERSE_GOLDEN_RATIO float64 = 1 / GOLDEN_RATIO
)

/* there is one
  ▄▀  █▄▄▄▄ ██   █ ▄▄   ▄  █
▄▀    █  ▄▀ █ █  █   █ █   █
█ ▀▄  █▀▀▌  █▄▄█ █▀▀▀  ██▀▀█
█   █ █  █  █  █ █     █   █
 ███    █      █  █       █
       ▀      █    ▀     ▀
             ▀          */

func NewGraph(surface *sdl.Surface) *Graph {
	r := &Graph{
		surface: surface,
	}
	r.Initialize()
	return r
}

type Graph struct {
	surface *sdl.Surface
}

func (r *Graph) Initialize() {
}

func (r *Graph) Render() {
	for i := float64(-9); i < -1; i++ {
		r.RenderGuideV(i, 0.15)
	}
	for i := float64(9); i > 1; i-- {
		r.RenderGuideV(i, 0.15)
	}
	r.RenderGuideV(-1, 0.25)
	r.RenderGuideV(1, 0.25)
	r.RenderGuideV(0, 0.35)
	r.RenderGuideH(-0.618, 0.25)
	r.RenderGuideH(0.618, 0.25)
	r.RenderGuideH(0, 0.5)
	for i := float64(-10); i <= 10; i += 0.03 {
		r.RenderAlgorithm(i, 1)
	}
}

func (r *Graph) RenderAlgorithm(i float64, brightness float64) {
	x := r.CoordI(i)
	y := r.CoordO(r.Algorithm(i))
	sBox := sdl.Rect{x, y, int32(graphPointSize), int32(graphPointSize)}
	r.surface.FillRect(&sBox, 0xFFFFFFFF)
}

func (r *Graph) RenderGuideH(i float64, brightness float64) {
	y := r.CoordO(i)
	sBox := sdl.Rect{0, y, int32(graphWidth * graphPointSize), int32(graphPointSize)}
	r.surface.FillRect(&sBox, colorBrightness(brightness))
}

func (r *Graph) RenderGuideV(i float64, brightness float64) {
	x := r.CoordI(i)
	sBox := sdl.Rect{x, 0, int32(graphPointSize), int32(graphHeight * graphPointSize)}
	r.surface.FillRect(&sBox, colorBrightness(brightness))
}

func (r *Graph) CoordI(i float64) int32 {
	return int32(float64(graphWidth * graphPointSize) * (i / 20 + 0.5))
}

func (r *Graph) CoordO(i float64) int32 {
	return int32(float64(graphWidth * graphPointSize) * (i / 2 + 0.5))
}

func (r *Graph) Algorithm(i float64) float64 {
	if i < -1 {
		return -math.Log(-i - 0.85) / 14 - 0.75
	} else if i > 1 {
		return math.Log(i - 0.85) / 14 + 0.75
	} else {
		return i / GOLDEN_RATIO
	}
}

/*
the graph
is within the
██   █ ▄▄  █ ▄▄
█ █  █   █ █   █
█▄▄█ █▀▀▀  █▀▀▀
█  █ █     █
   █  █     █
  █    ▀     ▀
 ▀
*/

func NewApp() *App {
	g := &App{
		Name: "graph",
	}
	g.Initialize()
	return g
}

type App struct {
	/* public */
	Name string
	/* private objects */
	graph *Graph
	/* private */
	state StateEnum
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

func (g *App) Initialize() {
	var err error
	log.SetLevel(log.DebugLevel)

	log.WithFields(log.Fields{
		"name": g.Name,
	}).Info("Initialize App")

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
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
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

	g.graph = NewGraph(g.sdlScreenSurface)

	g.ChangeState(STATE_LOADING)
}

func (g *App) Start() int {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"recover": r,
			}).Warn("App Recovered")
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

func (g *App) Render() {
	var err error

	g.sdlScreenSurface.FillRect(nil, 0xFF000000)

	g.graph.Render()

	g.sdlScreenTexture, err = g.sdlRenderer.CreateTextureFromSurface(g.sdlScreenSurface)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Could not create texture from surface")
	}
	defer g.sdlScreenTexture.Destroy()

	g.sdlRenderer.Copy(g.sdlScreenTexture, graphRenderOffsetSrc, nil)

	g.sdlRenderer.Present()
}

func (g *App) Stop() {
	g.ChangeState(STATE_FINISHED)
}

func (g *App) Teardown() {
	log.Info("Teardown App")
	g.sdlRenderer.Destroy()
	g.sdlWindow.Destroy()
}

func (g *App) ChangeState(s StateEnum) {
	g.state = s
	log.WithFields(log.Fields{
		"state": g.StateName(),
	}).Info("App changed")
	switch g.state {
	case STATE_LOADING:
	case STATE_PLAYING:
	case STATE_FINISHED:
	case STATE_FAILED:
	}
}

func (g *App) StateName() string {
	switch g.state {
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

func (g *App) PollEvents() {
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

func (g *App) Alive() bool {
	return g.state < STATE_FINISHED
}

func (g *App) NowMs() bool {
	g.nowMs = sdl.GetTicks()
	if g.nowMs != g.lastMs {
		g.lastMs = g.nowMs
		return true
	}
	return false
}

/*
the app is
instantiated from
█▀▄▀█ ██   ▄█    ▄
█ █ █ █ █  ██     █
█ ▄ █ █▄▄█ ██ ██   █
█   █ █  █ ▐█ █ █  █
   █     █  ▐ █  █ █
  ▀     █     █   ██
       ▀          */

func main() {
	runtime.LockOSThread()
	app := NewApp()
	os.Exit(app.Start())
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
	winWidth            = graphWidth * graphPointSize
	winHeight           = graphHeight*graphPointSize - graphPointSize*graphGenRows
	graphLimitX          = graphWidth - 1
	graphLimitY          = graphHeight - 1
	graphCenterX         = graphWidth / 2
	graphRenderOffsetSrc = &sdl.Rect{0, 0, int32(winWidth), int32(winHeight - graphPointSize*graphGenRows)}
)

var palette = []uint32{
	0xFF000000,
	0xFF251d1a,
	0xFF3b2d23,
	0xFF5a372d,
	0xFF72432e,
	0xFF9c562f,
	0xFFbc5b26,
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
