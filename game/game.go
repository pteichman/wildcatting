package game

import (
	"math/rand"
	"time"
)

type Game interface {
	Join(string) int
	Move(Move) *State
	View(int) *State
}

// Update is a JSON serializable update containing the result to a player Move.
type Update interface{}

type Move struct {
	PlayerID int
	SiteID   int
	Done     bool
}

type Field struct {
	Prob []int `json:"prob"`
}

type game struct {
	f       *field
	week    int
	deeds   map[int]*deed // site id to ownership record
	price   int           // oil price in cents
	players []string      // id indexed player names
	turn    int           // next surveying turn (which much happens in order)
	move    chan Move
	update  []chan *State
}

type deed struct {
	player  int
	start   int
	stop    int
	bit     int
	barrels int
	pnl     int
}

func New() Game {
	rand.Seed(time.Now().UTC().UnixNano())

	g := &game{
		f:     newField(),
		deeds: make(map[int]*deed),
		move:  make(chan Move),
	}

	go g.run()
	return g
}

func (g *game) run() {
	for state := lobby; state != nil; {
		state = state(g)
	}
}

func (g *game) Join(name string) int {
	// we probably need to lock the players slice
	p := len(g.players)
	g.players = append(g.players, name)

	g.update = append(g.update, make(chan *State))

	return p
}

func (g *game) Move(mv Move) *State {
	g.move <- mv
	return <-g.update[mv.PlayerID]
}

// View returns a players' View of the oil field state.
func (g *game) View(playerID int) *State {
	var deeds []Deed
	for s, deed := range g.deeds {
		d := Deed{
			SiteID: s,
			Owner:  deed.player,
			Sold:   deed.stop > 0,
			Cost:   g.f.cost[s] * deed.bit, // cost is in cents and bit is in 100 ft increments so they cancel out
			Income: deed.barrels * g.price / 100,
			PNL:    deed.pnl,
		}
		// players only know about oil if it was reached with the bit
		if g.f.oil[s] > 0 && deed.bit == g.f.oil[s] {
			d.Oil = g.f.oil[s]
		}

		deeds = append(deeds, d)
	}

	return &State{
		Players: g.players,
		Week:    g.week,
		Prob:    g.f.prob,
		Cost:    g.f.cost,
		Tax:     g.f.tax,
		Oil:     g.f.oil,
		Deeds:   deeds,
		Fact:    facts[rand.Intn(len(facts))],
	}
}
