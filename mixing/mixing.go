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
	numSamples = 4096
)

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot init SDL")
		return
	}
	defer sdl.Quit() // TODO: wrap this with func: recover from panic (pause audio stream?) then sdl.Quit

	var (
		start = time.Now().Add(1 * time.Second) // 1 second delay before start
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

	spec := atomix.Spec(sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_U16,
		Channels: 2,
		Samples:  numSamples,
	})

	t := start
	loopDur := 16 * step
	totalDur := time.Duration(0)
	for n := 0; n < loops; n++ {
		atomix.Play(kick1, t, 1)
		atomix.Play(marac, t.Add(1 * step), 0.5)
		atomix.Play(snare, t.Add(4 * step), 0.8)
		atomix.Play(marac, t.Add(6 * step), 0.5)
		atomix.Play(kick2, t.Add(7 * step), 0.9)
		atomix.Play(marac, t.Add(10 * step), 0.5)
		atomix.Play(kick2, t.Add(10 * step), 0.9)
		atomix.Play(snare, t.Add(12 * step), 0.8)
		atomix.Play(marac, t.Add(14 * step), 0.5)
		t = t.Add(loopDur)
		totalDur += loopDur
	}
	runLength := totalDur + 2 * time.Second

	sdl.OpenAudio(spec, nil)
	sdl.PauseAudio(false)
	log.WithFields(log.Fields{
		"spec": spec,
	}).Info("")

	time.Sleep(runLength)
}
