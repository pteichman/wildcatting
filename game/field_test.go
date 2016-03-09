package game

import (
	"fmt"
	"sort"
	"testing"
)

type pt struct{ y, x int }

func TestNeighbors(t *testing.T) {

	f := newField(3, 3)

	// 0 1 2
	// 3 4 5
	// 6 7 8
	expect := [][]int{
		{1, 3},
		{0, 2, 4},
		{1, 5},
		{0, 4, 6},
		{1, 3, 5, 7},
		{2, 4, 8},
		{3, 7},
		{4, 6, 8},
		{5, 7},
	}

	for i, expect := range expect {
		j := 0
		for nbr := range f.neighbors(i) {
			if nbr != expect[j] {
				t.Errorf("neigbors(%d) -> element at index %d is %d; expect %d", i, j, nbr, expect[j])
			}
			j++
		}
	}
}

var reservoirTests = []struct {
	oil    []int
	expect [][]int
}{
	{
		oil: []int{
			0, 0, 0,
			0, 0, 0,
			0, 0, 0},
		expect: [][]int{{}, {}, {}, {}, {}, {}, {}, {}, {}}},
	{
		oil: []int{
			0, 1, 0,
			1, 1, 1,
			0, 1, 0},
		expect: [][]int{
			{},
			{1, 3, 4, 5, 7},
			{},
			{1, 3, 4, 5, 7},
			{1, 3, 4, 5, 7},
			{1, 3, 4, 5, 7},
			{},
			{1, 3, 4, 5, 7},
			{},
		}},
	{
		oil: []int{
			1, 1, 1,
			2, 1, 2,
			3, 3, 3},
		expect: [][]int{
			{0, 1, 2, 4},
			{0, 1, 2, 4},
			{0, 1, 2, 4},
			{3},
			{0, 1, 2, 4},
			{5},
			{6, 7, 8},
			{6, 7, 8},
			{6, 7, 8},
		}},
}

func TestReservoir(t *testing.T) {
	for i, test := range reservoirTests {
		f := newField(3, 3)
		f.oil = test.oil
		for s := 0; s < 9; s++ {
			res := f.reservoir(s)
			sort.Sort(sort.IntSlice(res))
			if len(res) != len(test.expect[s]) {
				t.Errorf("reservoir test %d len(reservoir(%d)) -> %d; expect %d", i, s, len(res), len(test.expect[s]))
				continue
			}
			fmt.Println("reservoir: ", res)
			for r := 0; r < len(res); r++ {
				if res[r] != test.expect[s][r] {
					t.Errorf("reservoir test %d reservoir(%d)[%d] == %d; expect %d", i, s, r, res[r], test.expect[s][r])
				}
			}
		}
	}
}
