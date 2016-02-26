package game

import "testing"

type pt struct{ y, x int }

func TestNeighbors(t *testing.T) {
	expect := map[int]bool{
		40:  true,
		119: true,
		121: true,
		200: true,
	}
	ct := 0
	for nbr := range neighbors(120) {
		if _, ok := expect[nbr]; !ok {
			t.Errorf("unexpected neighbor %d", nbr)
		}
		ct++
	}
	if ct != 4 {
		t.Errorf("expected 4 neighbors; got %d", ct)
	}
}

var reservoirTests = []struct {
	depth int
	oil   []int
	site  int
}{
	{9, []int{117}, 117},
	{3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1},
	{4, []int{0, 80, 160, 240, 320}, 320},
	{6, []int{118, 119, 120, 121, 122, 40, 200}, 120},
}

func TestReservoir(t *testing.T) {

	for _, test := range reservoirTests {
		f := &field{
			oil: make([]int, 80*24),
		}
		expect := make(map[int]bool)
		for _, s := range test.oil {
			f.oil[s] = test.depth
			expect[s] = true
		}

		res := f.reservoir(test.site)
		if len(res) != len(expect) {
			t.Errorf("expected %d neighbors; got %d", len(expect), len(res))
		}
		for _, nbr := range res {
			if _, ok := expect[nbr]; !ok {
				t.Errorf("unexpected neighbor %d", nbr)
			}
		}
	}
}
