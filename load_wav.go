/** Author: Charney Kaye */

package main

// typedef unsigned char Uint16;
// void AudioCallback(void *userdata, Uint16 *stream, int len);
import "C"
import (
	log "github.com/Sirupsen/logrus"
	// "io"
	"math"
	// "fmt"
	"os"
	"time"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	toneHz   = 260
	sampleHz = 44100
	dPhase   = 2 * math.Pi * toneHz / sampleHz
)

// read audio file

func ReadAudio() {

	file := os.Args[1]

	data, spec := sdl.LoadWAV(file, &sdl.AudioSpec{})

	log.WithFields(log.Fields{
		"spec":   spec,
	}).Info("Loaded")

	for n := 0; n < len(data); n += 2 {
		StoreSample(data[n:n+2])
	}
/*
	wavReader, err := wav.NewReader(testWav, testInfo.Size())
	if err != nil {
		panic(err)
	}

	storedAudio = make([]C.Uint16, 0)

sampleLoop:
	for {
		s, err := wavReader.ReadSample()
		if err == io.EOF {
			break sampleLoop
		} else if err != nil {
			panic(err)
		}
		storedAudio = append(storedSample, C.Uint16(s))
	}
*/
}

func StoreSample(s []byte) {
	storedAudio = append(storedAudio, C.Uint16(int32(s[0]) + int32(s[1])<<8))
}

var storedAudio []C.Uint16

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot init SDL")
		return
	}
	defer sdl.Quit()

	// spec := &sdl.AudioSpec{
	// 	Freq:     sampleHz,
	// 	Format:   sdl.AUDIO_U16,
	// 	Channels: 1,
	// 	Samples:  sampleHz,
	// 	Callback: sdl.AudioCallback(C.AudioCallback),
	// }
	// sdl.OpenAudio(spec, nil)
	// sdl.PauseAudio(false)

	ReadAudio()

	time.Sleep(1 * time.Second)
}
