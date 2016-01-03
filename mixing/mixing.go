/** Author: Charney Kaye */

package main

// typedef unsigned char Uint16;
// void AudioCallback(void *userdata, Uint16 *stream, int len);
import "C"
import (
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/outrightmental/go-atomix"
)

const (
	sampleHz = 44100
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
		beat = 500 * time.Millisecond
		loops = 4
		)

	var (
		p808 = "assets/sounds/percussion/808/"
		kick1 = p808 + "kick1.wav"
		kick2 = p808 + "kick2.wav"
		snare = p808 + "snare.wav"
		marac = p808 + "maracas.wav"
		)

	for n := 0; n < loops; n++ {
		atomix.Play(kick1, start, 1)
		atomix.Play(marac, start + 0.5 * beat, 0.5)
		atomix.Play(snare, start + 1 * beat, 0.8)
		atomix.Play(marac, start + 1.5 * beat, 0.5)
		atomix.Play(kick2, start + 1.75 * beat, 0.9)
		atomix.Play(marac, start + 2.5 * beat, 0.5)
		atomix.Play(kick2, start + 2.5 * beat, 0.9)
		atomix.Play(snare, start + 3 * beat, 0.8)
		atomix.Play(marac, start + 3.5 * beat, 0.5)
		start += 4 * beat
	}

	spec := atomix.Spec(&sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_U16,
		Channels: 2,
		Samples:  numSamples,
	})
	sdl.OpenAudio(spec, nil)
	sdl.PauseAudio(false)
	log.WithFields(log.Fields{
		"spec": spec,
		}).Info("")

	time.Sleep(1 * time.Second)
}
