package game

import "testing"

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
		height: 3,
		width:  3,
		prob: []int{
			50, 50, 50,
			50, 50, 50,
			50, 50, 50},
		cost: []int{
			10, 10, 10,
			10, 10, 10,
			10, 10, 10},
		oil: []int{
			0, 0, 0,
			0, 0, 0,
			0, 0, 0},
		tax: []int{
			100, 100, 100,
			100, 100, 100,
			100, 100, 100},
	},
	joins: []string{"bob", "peter", "joe"},
	weeks: []testWeek{
		testWeek{
			surveys: []int{0, 1, 2},
			reports: []int{yes, yes, yes},
			drills:  []int{1, 2, 3},
			sells:   [][]int{{}, {}, {}},
		},
		testWeek{
			surveys: []int{3, 4, 5},
			reports: []int{no, no, no},
			drills:  []int{0, 0, 0},
			sells:   [][]int{{}, {}, {}},
		},
		testWeek{
			surveys: []int{6, 7, 8},
			reports: []int{yes, yes, yes},
			drills:  []int{0, 0, 0},
			sells:   [][]int{{}, {}, {}},
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
				g.Move(p, 0)
			}
			s := tw.surveys[p]
			if g.deeds[s].bit != n {
				t.Errorf("drilling (week %d player %d site %d): expect bit %d; got %d", g.week, p, s, n, g.deeds[s].bit)
			}
			if tw.drills[p] > 0 && g.deeds[s].week != g.week {
				t.Errorf("drilling (week %d player %d site %d): expect start %d; got %d", g.week, p, s, g.week, g.deeds[s].week)
			}
		}

		// stop drilling where we were
		for p, yesNo := range tw.reports {
			if yesNo == yes {
				g.Move(p, done)
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
			g.Move(p, done)
		}

		// begin next week
		g.Move(0, 0)
	}
}
