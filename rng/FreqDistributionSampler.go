/*Copyright 2019 Salah Mahmud. All rights reserved.*/
// Package rng

package rng

import (
	"math/rand"

	"github.com/MichaelTJones/pcg"
)

// FreqDistributionSampler provides constant time sampling from a discrete distribution
// modified from go-discrererand Gryski <damian@gryski.com>
// This is an implementation of Vose's labels method for
// choosing elements from a discrete distribution.
// For a full description of the algorithm, see http://www.keithschwarz.com/darts-dice-coins/
type FreqDistributionSampler struct {
	rnd    *rand.Rand
	labels []int
	prob   []float32
	size   int
	pcg    *pcg.PCG32
}

// NewFreqDistributionSampler constructs an FreqDistributionSampler  that will generate the discrete distribution given in probabilities.
// The probabilities array should sum to 1.
func NewFreqDistributionSampler(probabilities []float32, src rand.Source) FreqDistributionSampler {

	n := len(probabilities)

	fds := FreqDistributionSampler{}
	fds.labels = make([]int, n)
	fds.prob = make([]float32, n)
	fds.rnd = rand.New(src)
	fds.size = len(fds.labels)
	fds.pcg = pcg.NewPCG32().Seed(1, 1)

	p := make([]float32, n)
	for i := 0; i < n; i++ {
		p[i] = probabilities[i] * float32(n)
	}

	var small worklist
	var large worklist

	for i, pi := range p {
		if pi < 1 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	for len(large) > 0 && len(small) > 0 {
		l := small.pop()
		g := large.pop()
		fds.prob[l] = p[l]
		fds.labels[l] = g

		p[g] = (p[g] + p[l]) - 1
		if p[g] < 1 {
			small.push(g)
		} else {
			large.push(g)
		}
	}

	for len(large) > 0 {
		g := large.pop()
		fds.prob[g] = 1
	}

	for len(small) > 0 {
		l := small.pop()
		fds.prob[l] = 1
	}

	return fds
}

// Next returns the next random value from the discrete distribution
func (fds *FreqDistributionSampler) NextStdRand() int {
	i := fds.rnd.Intn(fds.size)
	if fds.rnd.Float32() < fds.prob[i] {
		return i
	}
	return fds.labels[i]
}

func (fds *FreqDistributionSampler) Next() int {
	i := fds.pcg.FastBounded(uint32(fds.size))
	if fds.pcg.Float32() < fds.prob[i] {
		return int(i)
	}
	return fds.labels[i]
}

var seed = 123455

func (fds *FreqDistributionSampler) NextInaccurate() int {
	i := fds.rnd.Intn(fds.size)
	if sfrand(&seed) < fds.prob[i] {
		return int(i)
	}
	return fds.labels[i]
}
