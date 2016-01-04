/** Author: Charney Kaye */

package main

// typedef unsigned char Uint16;
// void AudioCallback(void *userdata, Uint16 *stream, int len);
import "C"
import (
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/outrightmental/go-atomix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	sampleHz   = 44100
	numSamples = 22050
)

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot init SDL")
		return
	}
    defer func() {
        if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"recover": r,
			}).Warn("Player Recovered")
        }
		sdl.PauseAudio(true)
		atomix.Teardown()
		sdl.Quit()
    }()

	var (
		step = 125 * time.Millisecond
		loops = 4
	)

	var (
		p808  = "assets/sounds/percussion/808/"
		kick1 = p808 + "kick1.wav"
		kick2 = p808 + "kick2.wav"
		snare = p808 + "snare.wav"
		marac = p808 + "maracas.wav"
	)

	atomix.Debug(true)
	atomix.Configure(sdl.AudioSpec{
			Freq:     sampleHz,
			Format:   sdl.AUDIO_U16,
			Channels: 2,
			Samples:  numSamples,
		})

	t := 1 * time.Second // padding before music
	for n := 0; n < loops; n++ {
        atomix.Play(kick1, t,            4 *step,  1.0)
        atomix.Play(marac, t + 1 *step,  1 *step,  0.5)
        atomix.Play(snare, t + 4 *step,  4 *step,  0.8)
        atomix.Play(marac, t + 6 *step,  1 *step,  0.5)
        atomix.Play(kick2, t + 7 *step,  4 *step,  0.9)
        atomix.Play(marac, t + 10 *step, 1 *step,  0.5)
        atomix.Play(kick2, t + 10 *step, 4 *step,  0.9)
        atomix.Play(snare, t + 12 *step, 4 *step,  0.8)
        atomix.Play(marac, t + 14 *step, 1 *step,  0.5)
		t += 16 * step
	}

	spec := atomix.Spec()
	sdl.OpenAudio(spec, nil)
	sdl.PauseAudio(false)
	log.WithFields(log.Fields{
		"spec": spec,
	}).Info("SDL OpenAudio > Atomix")

	time.Sleep(t + 1 * time.Second) // padding after music
}
