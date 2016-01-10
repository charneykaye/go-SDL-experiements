/** Author: Charney Kaye */

package main

// typedef unsigned char Uint8;
// void AudioCallback(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"encoding/binary"
	log "github.com/Sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
	"reflect"
	"time"
	"runtime/debug"
	"unsafe"
)

var (
	sampleFile string = "song.wav"
	samples   []uint16
	nowSample int
)

func LoadSample(file string) *sdl.AudioSpec {
	data, spec := sdl.LoadWAV(file, &sdl.AudioSpec{})
	for n := 0; n < len(data); n += 2 {
		samples = append(samples, binary.BigEndian.Uint16(data[n:n+2]))
	}
	return spec
}

func NextSample() (s uint16) {
	if nowSample < len(samples) {
		s = samples[nowSample]
	}
	nowSample++
	return
}

func NextSampleBytes() (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, NextSample())
	return
}

//export AudioCallback
func AudioCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))
	for i := 0; i < n; i += 2 {
		b := NextSampleBytes()
		buf[i] = C.Uint8(b[0])
		buf[i+1] = C.Uint8(b[1])
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot init SDL")
		return
	}
    defer func() {
        if r := recover(); r != nil {
			stk := debug.Stack()
			log.WithFields(log.Fields{
				"stack": string(stk[:]),
				"recover": r,
			}).Warn("Player Recovered")

        }
		sdl.PauseAudio(true)
		sdl.Quit()
    }()

	loadSpec := LoadSample(sampleFile)
	if loadSpec != nil {
		log.WithFields(log.Fields{
			"spec": loadSpec,
			}).Info("Loaded")
	} else {
		log.WithFields(log.Fields{
			"file": sampleFile,
			}).Fatal("Failed to load")
	}

	loadSpec.Callback = sdl.AudioCallback(C.AudioCallback)

	// playSpec := &sdl.AudioSpec{
	// 		Freq:     loadSpec.Freq,
	// 		Format:   loadSpec.Format,
	// 		Channels: loadSpec.Channels,
	// 		Samples:  loadSpec.Samples,
	// 		Callback: sdl.AudioCallback(C.AudioCallback),
	// 	}

	sdl.OpenAudio(loadSpec, nil)
	sdl.PauseAudio(false)

	time.Sleep(18 * time.Second)
}
