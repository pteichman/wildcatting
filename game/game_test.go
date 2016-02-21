package game

import (
	"fmt"
	"testing"
	"time"
)

type testGame struct {
	f     *field
	joins []string
	weeks []testWeek
}

type testWeek struct {
	surveys []int
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
			drills:  []int{3, 5, 0, 0},
			sells:   [][]int{{13}, {}, {}, {}},
		},
		testWeek{
			surveys: []int{12, 33, 36, 1913},
			drills:  []int{3, 5, 0, 0},
			sells:   [][]int{{12}, {33}, {}, {}},
		},
		testWeek{
			surveys: []int{11, 30, 35, 300},
			drills:  []int{0, 0, 5, 5},
			sells:   [][]int{{}, {}, {}, {}},
		},
		testWeek{
			surveys: []int{21, 40, 45, 400},
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
			sells:   [][]int{{11, 21, 31, 41}, {30, 40, 50, 60}, {35, 45, 55, 65}, {300, 400, 500, 600}},
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
	g.Move(Move{Done: true})

	for i, tw := range tg.weeks {
		time.Sleep(time.Millisecond)

		w := i + 1

		if g.week != w {
			t.Errorf("expect week %d; got %d", w, g.week)
			return
		}

		// surveys
		for p, s := range tw.surveys {
			fmt.Printf("move <- Move{PlayerID: %d, SiteID: %d}\n", p, s)
			g.Move(Move{PlayerID: p, SiteID: s})

		}

		for p, s := range tw.surveys {
			deed := g.deeds[s]
			if deed.player != p {
				t.Errorf("surveying (week %d player %d site %d): expect owner %d; got %d", g.week, p, s, p, deed.player)
			}
		}

		// drilling
		for p, n := range tw.drills {
			s := tw.surveys[p]
			var oil bool
			for i := 0; i < n; i++ {
				fmt.Printf("move <- Move{PlayerID: %d,}\n", p)
				g.Move(Move{PlayerID: p})

				if n == g.f.oil[s] {
					oil = true
					break
				}
			}

			if !oil && n < maxOil {
				fmt.Printf("move <- Move{PlayerID: %d, Done: true}\n", p)
				g.Move(Move{PlayerID: p, Done: true})
			}

			if g.deeds[s].bit != n {
				t.Errorf("drilling (week %d player %d site %d): expect bit %d; got %d", g.week, p, s, n, g.deeds[s].bit)
			}

			if tw.drills[p] > 0 && g.deeds[s].start != g.week {
				t.Errorf("drilling (week %d player %d site %d): expect start %d; got %d", g.week, p, s, g.week, g.deeds[s].start)
			}
		}

		// selling
		for p, sells := range tw.sells {
			for _, s := range sells {
				fmt.Printf("move <- Move{PlayerID: %d, SiteID: %d}\n", p, s)
				g.Move(Move{PlayerID: p, SiteID: s})

				if g.deeds[s].stop != g.week {
					t.Errorf("selling (week %d player %d site %d): expect stop %d; got %d", g.week, p, s, g.week, g.deeds[s].stop)
				}
			}

			fmt.Printf("move <- Move{PlayerID: %d, Done: true}\n", p)
			g.Move(Move{PlayerID: p, Done: true})

		}
	}
}
