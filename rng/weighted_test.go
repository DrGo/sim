package rng

import (
	"fmt"
	"testing"
	"time"
)

const (
	nbGeneration = 100000
	setSize      = 10
)

func TestGeneratorDistribution(t *testing.T) {
	g := newGenerator(time.Now().Unix(), setSize)

	numbers := map[float32]float32{}

	for i := 0; i < nbGeneration; i++ {
		r := g.Random()
		if _, ok := numbers[r]; !ok {
			numbers[r] = 0
		}
		numbers[r]++
	}
	startBound := float32(0)
	for index, value := range g.values {

		got := numbers[value]
		want := (g.bounds[index] - startBound) * nbGeneration
		delta := want * 5 / 100

		if got > want+delta || got < want-delta {
			t.Errorf("distribution not correct for value %f, expected %f +/- 5%%, but got %f", value, want, got)
		}
		startBound = g.bounds[index]
	}
}

func TestGeneratorSeeding(t *testing.T) {

	seed := int64(1)
	firstRun := generate(seed, nbGeneration)
	secondRun := generate(seed, nbGeneration)

	for i := 0; i < nbGeneration; i++ {
		if firstRun[i] != secondRun[i] {
			t.Errorf("expected same sequence of number, but at pos %d, got %f and %f", i, firstRun[i], secondRun[i])
		}
	}
}

func generate(seed int64, size int) []float32 {
	g := newGenerator(seed, setSize)
	runVal := make([]float32, 0, size)
	for i := 0; i < size; i++ {
		runVal = append(runVal, g.Random())
	}
	return runVal
}

func newGenerator(seed int64, size int) *Generator {
	values := make([]float32, 0, size)
	weight := make([]float32, 0, size)

	p := float32(1) / float32(size)
	for i := 0; i < size; i++ {
		values = append(values, float32(i))
		weight = append(weight, p)
	}
	g, _ := NewGenerator(seed, values, weight)
	fmt.Printf("%+v", g)
	return g
}

func BenchmarkScalingImp2(b *testing.B) {
	for size := 2; size <= 1024; size = size * 2 {
		name := fmt.Sprintf("numberSet_size_%d", size)
		b.Run(name, func(b *testing.B) {
			g := newGenerator(1, size)
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				g.Random()
			}
		})
	}
}
