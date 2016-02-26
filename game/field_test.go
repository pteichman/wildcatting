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

// FIXME table driven reservoir tests something like this
var reservoirTests = []struct {
	depth int
	oil   []int
}{
	{6, []int{118, 119, 120, 121, 122, 40, 200}},
}

func TestReservoir(t *testing.T) {
	f := &field{
		oil: make([]int, 80*24),
	}

	depth := 6
	f.oil[118] = depth
	f.oil[119] = depth
	f.oil[120] = depth
	f.oil[121] = depth
	f.oil[122] = depth

	f.oil[40] = depth
	f.oil[200] = depth

	expect := map[int]bool{
		40:  true,
		118: true,
		119: true,
		120: true,
		121: true,
		122: true,
		200: true,
	}
	res := f.reservoir(120)
	if len(res) != 7 {
		t.Errorf("expected 7 neighbors; got %d", len(res))
	}
	for _, nbr := range res {
		if _, ok := expect[nbr]; !ok {
			t.Errorf("unexpected neighbor %d", nbr)
		}
	}
}
