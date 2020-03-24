package rng

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/dgryski/go-xoroshiro"
)

func TestFreqDistributionSampler(t *testing.T) {

	p := []float32{0.125, 0.2, 0.1, 0.25, 0.1, 0.1, 0.125}
	counts := make([]int, len(p))

	a := NewFreqDistributionSampler(p, rand.NewSource(0))

	const rounds = 10e6

	for i := 0; i < rounds; i++ {
		n := a.Next()
		counts[n]++
	}

	const eps = 0.001
	for i := 0; i < len(counts); i++ {
		g := float32(counts[i]) / float32(rounds)
		fmt.Println(p[i], "=", g)
		if math.Abs(float64(g-p[i])) > eps {
			t.Errorf("failed element %d: got %f expected %f +/- %f", i, g, p[i], eps)
		}
	}
}

func TestFreqDistributionSamplerXoroshiro(t *testing.T) {

	p := []float32{0.125, 0.2, 0.1, 0.25, 0.1, 0.1, 0.125}
	counts := make([]int, len(p))
	seed := xoroshiro.SplitMix64(0x0ddc0ffeebadf00d)
	s := xoroshiro.State([2]uint64{seed.Next(), seed.Next()})
	a := NewFreqDistributionSampler(p, &s)

	const rounds = 10e6

	for i := 0; i < rounds; i++ {
		n := a.Next()
		counts[n]++
	}

	const eps = 0.001
	for i := 0; i < len(counts); i++ {
		g := float32(counts[i]) / float32(rounds)
		fmt.Println(p[i], "=", g)
		if math.Abs(float64(g-p[i])) > eps {
			t.Errorf("failed alias method test element %d: got %f expected %f +/- %f", i, g, p[i], eps)
		}
	}
}

// BenchmarkScalingAlias performance should not degrade with increasing distribution size
func BenchmarkScalingFreqDistributionSampler(b *testing.B) {
	for size := 2; size <= 1024; size = size * 2 {
		name := fmt.Sprintf("n levels=%d", size)
		b.Run(name, func(b *testing.B) {
			p := generateProbDist32(1, size)
			b.ResetTimer()
			a := NewFreqDistributionSampler(p, rand.NewSource(0))
			for n := 0; n < b.N; n++ {
				a.Next()
			}
		})
	}
}

func BenchmarkScalingFreqDistributionSamplerXoroshiro(b *testing.B) {
	seed := xoroshiro.SplitMix64(0x0ddc0ffeebadf00d)
	s := xoroshiro.State([2]uint64{seed.Next(), seed.Next()})
	for size := 2; size <= 1024; size = size * 2 {
		name := fmt.Sprintf("n levels=%d", size)
		b.Run(name, func(b *testing.B) {
			p := generateProbDist32(1, size)
			b.ResetTimer()
			a := NewFreqDistributionSampler(p, &s)
			for n := 0; n < b.N; n++ {
				a.NextStdRand()
			}
		})
	}
}

//utils

func generateProbDist64(seed int64, size int) []float64 {
	labels := make([]int, 0, size)
	sum := 0
	for i := 0; i < size; i++ {
		labels = append(labels, rand.Int())
		sum += labels[i]
	}
	probabilities := make([]float64, 0, size)
	for i := 0; i < size; i++ {
		probabilities = append(probabilities, float64(labels[i]/sum))
	}
	return probabilities
}

func generateProbDist32(seed int64, size int) []float32 {
	labels := make([]int, 0, size)
	sum := 0
	for i := 0; i < size; i++ {
		labels = append(labels, rand.Int())
		sum += labels[i]
	}
	probabilities := make([]float32, 0, size)
	for i := 0; i < size; i++ {
		probabilities = append(probabilities, float32(labels[i]/sum))
	}
	return probabilities
}
