# go-SDL-experiements

Experiments with **[go-sdl2](https://github.com/veandco/go-sdl2)** - Go bindings for the C++ library **[Simple DirectMedia Layer](https://www.libsdl.org/)** (SDL 2.0)

Author: [Charney Kaye](http://w.charney.io)

## Stars

    go run stars.go

![Stars](stars/screenshot.png)

## Radar Stars

    go run radar_stars.go

![Radar Stars](radar_stars/screenshot.png)

***TODO:*** Replace infinite `for alive {}` loop with a sdl.Event handler, that only performs the life/render work in between display frames?

## Fire

![Fire](fire/screenshot.png)

    go run fire.go

***TODO:*** Track mouse motion and generate fire on it

***TODO:*** Replace infinite `for alive {}` loop with a sdl.Event handler, that only performs the life/render work in between display frames?

## Stars in Goroutine (anti-pattern)

This performs poorly- it's included only to potentially investigate *why* it performs so poorly.

    go run stars_goroutine.go

# Tips

### Texture Garbage Collection 

Be sure to call `texture.Destroy` once you're done with a texture.

    defer g.sdlScreenTexture.Destroy()

### Don't overload the frame buffer

Enabling vsync for the renderer helped reduce cpu usage a good amount (so that it doesn't draw multiple times in the same frame).

    g.sdlRenderer, err = sdl.CreateRenderer(g.sdlWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
