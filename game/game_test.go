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

	var players []int
	for _, name := range tg.joins {
		playerID := g.Join(name)
		players = append(players, playerID)

		actual := g.world.Name(entity(playerID))
		if actual != name {
			t.Errorf("join: expect player %s; got %s", name, actual)
		}
	}

	// start the game
	g.Move(players[0], done)

	for i, tw := range tg.weeks {
		week := i + 1

		if g.week != week {
			t.Errorf("expect week %d; got %d", week, g.week)
			return
		}

		// surveys
		for i, s := range tw.surveys {
			p := players[i]
			g.Move(p, s)
			deed := g.deeds[site(s)]
			if int(deed.player) != players[i] {
				t.Errorf("surveying (week %d player %d site %d): expect owner %d; got %d", g.week, p, s, p, deed.player)
			}
		}

		// surveyor's reports
		for i, yesNo := range tw.reports {
			g.Move(players[i], yesNo)
		}

		// drilling
		for i, n := range tw.drills {
			p := players[i]
			for j := 0; j < n; j++ {
				g.Move(p, 0)
			}
			s := site(tw.surveys[i])
			if g.deeds[s].bit != n {
				t.Errorf("drilling (week %d player %d site %d): expect bit %d; got %d", g.week, p, s, n, g.deeds[s].bit)
			}
			if tw.drills[i] > 0 && g.deeds[s].week != g.week {
				t.Errorf("drilling (week %d player %d site %d): expect start %d; got %d", g.week, p, s, g.week, g.deeds[s].week)
			}
		}

		// stop drilling where we were
		for i, yesNo := range tw.reports {
			if yesNo == yes {
				g.Move(players[i], done)
			}
		}

		// selling
		for i, sells := range tw.sells {
			p := players[i]
			for _, s := range sells {
				g.Move(p, s)

				s := site(s)
				if g.deeds[s].stop != g.week {
					t.Errorf("selling (week %d player %d site %d): expect stop %d; got %d", g.week, p, s, g.week, g.deeds[s].stop)
				}
			}
			g.Move(p, done)
		}

		// begin next week
		g.Move(players[0], 0)
	}
}
