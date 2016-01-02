/** Author: Charney Kaye */

package main

// typedef unsigned char Uint16;
// void AudioCallback(void *userdata, Uint16 *stream, int len);
import "C"
import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"math"
	"reflect"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

)

const (
	toneHz   = 260
	sampleHz = 44100
	dPhase   = 2 * math.Pi * toneHz / sampleHz
)

// read audio file

func ReadAudio() {

	file := "assets/sounds/percussion/808/kick1.wav"

	data, spec := sdl.LoadWAV(file, &sdl.AudioSpec{})

	log.WithFields(log.Fields{
		"spec":   spec,
	}).Info("Loaded")

	for n := 0; n < len(data); n += 2 {
		StoreSample(data[n:n+2])
	}
}

func StoreSample(s []byte) {
	storedAudio = append(storedAudio, C.Uint16(int32(s[0]) + int32(s[1])<<8))
}

var storedAudio []C.Uint16
var	(
	defaultSample = uint16(0xFFFF)
	defaultAudio = C.Uint16(defaultSample)
)

//export AudioCallback
func AudioCallback(userdata unsafe.Pointer, stream *C.Uint16, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint16)(unsafe.Pointer(&hdr))

	for i := 0; i < n; i += 1 {
		if i < len(storedAudio) {
			buf[i] = storedAudio[i]
		} else {
			buf[i] = defaultAudio
		}
	}
	fmt.Printf("AudioCallback length %d\n", n)
}

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot init SDL")
		return
	}
	defer sdl.Quit()

	ReadAudio()

	spec := &sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_U16,
		Channels: 2,
		Samples:  sampleHz,
		Callback: sdl.AudioCallback(C.AudioCallback),
	}
	sdl.OpenAudio(spec, nil)
	sdl.PauseAudio(false)

	time.Sleep(1 * time.Second)
}
