// Copyright ©️ 2021 oddstream.games

package tetra

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100
	baseFreq   = 220
)

var audioContext = audio.NewContext(sampleRate)

// pianoAt returns an i-th sample of piano with the given frequency.
func pianoAt(i int, freq float64) float64 {
	// Create piano-like waves with multiple sin waves.
	amp := []float64{1.0, 0.8, 0.6, 0.4, 0.2}
	x := []float64{4.0, 2.0, 1.0, 0.5, 0.25}
	v := 0.0
	for j := 0; j < len(amp); j++ {
		// Decay
		a := amp[j] * math.Exp(-5*float64(i)*freq/baseFreq/(x[j]*sampleRate))
		v += a * math.Sin(2.0*math.Pi*float64(i)*freq*float64(j+1)/sampleRate)
	}
	return v / 5.0
}

// toBytes returns the 2ch little endian 16bit byte sequence with the given left/right sequence.
func toBytes(l, r []int16) []byte {
	if len(l) != len(r) {
		panic("len(l) must equal to len(r)")
	}
	// linear PCM (16bits little endian, 2 channel stereo) without a header (e.g. RIFF header)
	b := make([]byte, len(l)*4)
	for i := range l {
		b[4*i] = byte(l[i])
		b[4*i+1] = byte(l[i] >> 8)
		b[4*i+2] = byte(r[i])
		b[4*i+3] = byte(r[i] >> 8)
	}
	return b
}

var (
	pianoNoteSamples = map[int][]byte{}
)

func init() {
	start := time.Now()
	// Create a reference data and use this for other frequency.
	const refFreq = 110
	length := 4 * sampleRate * baseFreq / refFreq
	refData := make([]int16, length)
	for i := 0; i < length; i++ {
		refData[i] = int16(pianoAt(i, refFreq) * math.MaxInt16)
	}

	for i := 0; i < 32; i++ {
		freq := baseFreq * math.Exp2(float64(i-1)/12.0)

		// Calculate the wave data for the freq.
		length := 4 * sampleRate * baseFreq / int(freq)
		l := make([]int16, length)
		r := make([]int16, length)
		for i := 0; i < length; i++ {
			idx := int(float64(i) * freq / refFreq)
			if len(refData) <= idx {
				break
			}
			l[i] = refData[idx]
		}
		copy(r, l)
		n := toBytes(l, r)
		pianoNoteSamples[int(freq)] = n
	}
	elapsed := time.Since(start)
	println("elapsed", elapsed/1000.0/1000.0)
}

// playNote plays piano sound with the given frequency.
func playNote(freq float64) {
	f := int(freq)
	p := audio.NewPlayerFromBytes(audioContext, pianoNoteSamples[f])
	p.SetVolume(0.25)
	p.Play()
}

// PlayPianoNote plays a synthesised note
func PlayPianoNote(i int) {
	i = i % len(pianoNoteSamples)
	playNote(baseFreq * math.Exp2(float64(i-1)/12.0))
}
