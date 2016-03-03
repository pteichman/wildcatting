package game

import (
	"math"
	"math/rand"
)

const (
	minProb = 1
	maxProb = 100
	minOil  = 1
	maxOil  = 9
	minCost = 10
	maxCost = 250
	minTax  = 100
	maxTax  = 550
)

type field struct{ prob, cost, oil, tax []int }

func newField() *field {
	prob := fill(1+rand.Intn(4), minProb, maxProb, 0.05, 0.25, false) // a few well formed peaks
	oil := probFilter(fill(1, minOil, maxOil, 0.1, 0.5, true), prob)  // hardship

	return &field{
		prob: prob,
		cost: fill(5+rand.Intn(5), minCost, maxCost, 0.1, 0.25, true), // many chaotic peaks
		oil:  oil,
		tax:  fill(10+rand.Intn(10), minTax, maxTax, 0.1, 0.5, false), // local politics
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func closest(p int, peaks []int) (int, int) {
	var minIdx int
	minDist := 24 + 80
	for i, q := range peaks {
		d := abs(p/80-q/80) + abs(p%80-q%80)
		if d < minDist {
			minDist = d
			minIdx = i
		}
	}
	return minIdx, minDist
}

func fill(n, min, max int, decay, fuzz float64, inverse bool) []int {
	var peaks []int
	for i := 0; i < n; i++ {
		peaks = append(peaks, rand.Intn(24*80))
	}
	values := make([]int, 24*80, 24*80)
	for i := 0; i < 24*80; i++ {
		minIdx, minDist := closest(i, peaks)

		// ratio of the longest possible distance
		v := float64(minDist) / (24 + 80)

		// Double the value for a better input into log :/
		// This should be distilled to some sane 0.0 to 1.0 param
		// which allows for controlling the steepness of the peaks
		v *= 2

		// Logarithmically adjust the value, shifting and dividing
		// to get a nice curve roughly in the range of 0 to 1.
		v = 1.0 - ((math.Log(v) + 4.0) / 4.0)

		// Adjust for subsequent peaks which are progressively lower.
		v *= math.Pow(1.0-decay, float64(minIdx))

		// Apply some random fuzz to keep everyone guessing.
		v += 2.0 * (rand.Float64() - 0.5) * fuzz

		// Contain the final value between zero and one.
		v = math.Min(math.Max(v, 0.0), 1.0)

		values[i] = int(math.Floor(float64(min) + float64(max-min)*v))

		if inverse {
			values[i] = min + max - values[i]
		}
	}

	return values
}

func probFilter(vals, p []int) []int {
	filtered := make([]int, len(vals))
	for i, v := range vals {
		if rand.Intn(100) > p[i] {
			continue
		}
		filtered[i] = v
	}
	return filtered
}

func neighbors(s int) <-chan int {
	out := make(chan int)
	y, x := s/80, s%80
	go func() {
		if x-1 >= 0 {
			out <- 80*y + x - 1
		}
		if y-1 >= 0 {
			out <- 80*(y-1) + x
		}
		if x+1 < 80 {
			out <- 80*y + x + 1
		}
		if y+1 < 24 {
			out <- 80*(y+1) + x
		}
		close(out)
	}()
	return out
}

func (f *field) reservoir(s int) []int {
	var res []int
	visited := make(map[int]bool)
	frontier := []int{s}
	for len(frontier) > 0 {
		cur := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		visited[cur] = true

		if f.oil[cur] != f.oil[s] {
			continue
		}

		for nbr := range neighbors(cur) {
			if _, ok := visited[nbr]; ok {
				continue
			}
			frontier = append(frontier, nbr)
		}
		res = append(res, cur)
	}
	return res
}
