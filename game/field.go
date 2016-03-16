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

type field struct {
	height, width        int
	prob, cost, oil, tax []int
}

func newField(height, width int) *field {
	prob := fill(height, width, 1+rand.Intn(4), minProb, maxProb, 0.05, 0.25, false) // a few well formed peaks
	oil := probFilter(fill(height, width, 1, minOil, maxOil, 0.1, 0.5, true), prob)  // hardship

	return &field{
		height: height,
		width:  width,
		prob:   prob,
		cost:   fill(height, width, 5+rand.Intn(5), minCost, maxCost, 0.1, 0.25, true), // many chaotic peaks
		oil:    oil,
		tax:    fill(height, width, 10+rand.Intn(10), minTax, maxTax, 0.1, 0.5, false), // local politics
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func closest(height, width, p int, peaks []int) (int, int) {
	var minIdx int
	minDist := height + width
	for i, q := range peaks {
		d := abs(p/width-q/width) + abs(p%width-q%width)
		if d < minDist {
			minDist = d
			minIdx = i
		}
	}
	return minIdx, minDist
}

func fill(height, width, n, min, max int, decay, fuzz float64, inverse bool) []int {
	var peaks []int
	for i := 0; i < n; i++ {
		peaks = append(peaks, rand.Intn(height*width))
	}
	values := make([]int, height*width, height*width)
	for i := 0; i < height*width; i++ {
		minIdx, minDist := closest(height, width, i, peaks)

		// ratio of the longest possible distance
		v := float64(minDist) / float64(height+width)

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

func (f *field) neighbors(s site) <-chan site {
	out := make(chan site)
	y, x := s/site(f.width), s%site(f.width)
	go func() {
		if y-1 >= 0 {
			out <- site(f.width)*(y-1) + x
		}
		if x-1 >= 0 {
			out <- site(f.width)*y + x - 1
		}
		if x+1 < site(f.width) {
			out <- site(f.width)*y + x + 1
		}
		if y+1 < site(f.height) {
			out <- site(f.width)*(y+1) + x
		}
		close(out)
	}()
	return out
}

func (f *field) reservoir(s site) []site {
	var res []site
	visited := make(map[site]bool)
	frontier := []site{s}
	for len(frontier) > 0 {
		cur := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		visited[cur] = true

		if f.oil[s] == 0 || f.oil[cur] != f.oil[s] {
			continue
		}

		for nbr := range f.neighbors(cur) {
			if _, ok := visited[nbr]; ok {
				continue
			}
			frontier = append(frontier, nbr)
		}
		res = append(res, cur)
	}
	return res
}
