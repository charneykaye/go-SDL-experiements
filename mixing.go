/** Author: Charney Kaye */

package main

// typedef unsigned char Uint16;
// void AudioCallback(void *userdata, Uint16 *stream, int len);
import "C"
import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"math"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/cryptix/wav"
)

const (
	toneHz   = 260
	sampleHz = 44100
	dPhase   = 2 * math.Pi * toneHz / sampleHz
)

// read audio file

func ReadAudio() {

	file := "assets/sounds/percussion/808/kick1.wav"

	testInfo, err := os.Stat(file)
	if err != nil {
		panic(err)
	}

	testWav, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	wavReader, err := wav.NewReader(testWav, testInfo.Size())
	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"file":   file,
		"reader": wavReader,
	}).Info("Loaded")

	storedSample = make([]int32, 0)

sampleLoop:
	for {
		s, err := wavReader.ReadSample()
		if err == io.EOF {
			break sampleLoop
		} else if err != nil {
			panic(err)
		}
		storedSample = append(storedSample, s)
	}
}

var storedSample []int32

//export AudioCallback
func AudioCallback(userdata unsafe.Pointer, stream *C.Uint16, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint16)(unsafe.Pointer(&hdr))

	for i := 0; i < n; i += 1 {
		var s int32
		if i < len(storedSample) {
			//			fmt.Printf("i %d, storedSample[i] %+v\n",n, storedSample[i])
			s = storedSample[i]
			//			s = 0xFFFF
		} else {
			s = 0xFFFF
		}
		c := C.Uint16(uint16(s))
		buf[i] = c
//		buf[i+1] = c
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
