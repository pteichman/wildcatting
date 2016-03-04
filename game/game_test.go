package game

import (
	"fmt"
	"testing"
)

type testGame struct {
	f     *field
	joins []string
	weeks []testWeek
}

type testWeek struct {
	surveys []int
	reports []int
	drills  []int
	sells   [][]int
}

var tg = testGame{
	f: &field{
		prob: make([]int, 24*80),
		cost: make([]int, 24*80),
		oil:  make([]int, 24*80),
		tax:  make([]int, 24*80),
	},
	joins: []string{"bob", "peter", "joe", "jbz"},
	weeks: []testWeek{
		testWeek{
			surveys: []int{13, 42, 37, 1914},
			reports: []int{yes, yes, no, no},
			drills:  []int{3, 5, 0, 0},
			sells:   [][]int{{}, {}, {}, {}},
		},
		testWeek{
			surveys: []int{12, 33, 36, 1913},
			reports: []int{yes, yes, no, no},
			drills:  []int{3, 5, 0, 0},
			sells:   [][]int{{}, {}, {}, {}},
		},
		testWeek{
			surveys: []int{11, 30, 35, 300},
			reports: []int{no, no, yes, yes},
			drills:  []int{0, 0, 5, 5},
			sells:   [][]int{{}, {}, {}, {}},
		},
		testWeek{
			surveys: []int{21, 40, 45, 400},
			reports: []int{yes, yes, yes, yes},
			drills:  []int{9, 9, 9, 9},
			sells:   [][]int{{}, {}, {}, {}},
		},
		testWeek{
			surveys: []int{31, 50, 55, 500},
			drills:  []int{9, 9, 9, 9},
			sells:   [][]int{{}, {}, {}, {}},
		},
		testWeek{
			surveys: []int{41, 60, 65, 600},
			drills:  []int{9, 9, 9, 9},
			sells:   [][]int{{}, {}, {}, {}},
		},
	},
}

func TestGame(t *testing.T) {
	g := New().(*game)
	g.f = tg.f

	for p, name := range tg.joins {
		playerID := g.Join(name)

		if playerID != p {
			t.Errorf("join: expect playerID %d; got %d", p, playerID)
		}

		if g.players[p] != name {
			t.Errorf("join: expect player %s; got %s", name, g.players[p])
		}
	}

	// start the game
	g.Move(0, done)

	for i, tw := range tg.weeks {
		week := i + 1

		if g.week != week {
			t.Errorf("expect week %d; got %d", week, g.week)
			return
		}

		// surveys
		for p, s := range tw.surveys {
			if p > 0 {
				<-g.view[p]
			}

			g.Move(p, s)
			deed := g.deeds[s]
			if deed.player != p {
				t.Errorf("surveying (week %d player %d site %d): expect owner %d; got %d", g.week, p, s, p, deed.player)
			}
		}

		// surveyor's reports
		for p, yesNo := range tw.reports {
			g.Move(p, yesNo)
		}

		// drilling
		for p, n := range tw.drills {
			for i := 0; i < n; i++ {
				g.Move(p, 1)
			}
			s := tw.surveys[p]
			if g.deeds[s].bit != n {
				t.Errorf("drilling (week %d player %d site %d): expect bit %d; got %d", g.week, p, s, n, g.deeds[s].bit)
			}
			if tw.drills[p] > 0 && g.deeds[s].week != g.week {
				t.Errorf("drilling (week %d player %d site %d): expect start %d; got %d", g.week, p, s, g.week, g.deeds[s].week)
			}
		}

		// selling
		for p, sells := range tw.sells {
			for _, s := range sells {
				g.Move(p, s)

				if g.deeds[s].stop != g.week {
					t.Errorf("selling (week %d player %d site %d): expect stop %d; got %d", g.week, p, s, g.week, g.deeds[s].stop)
				}
			}
			fmt.Printf("sending sell done to %d\n", p)
			g.Move(p, done)
			fmt.Printf("SENT sell done to %d\n", p)
		}

		// finish week
		for p := range tg.joins {
			g.Move(p, done)
		}
	}
}
