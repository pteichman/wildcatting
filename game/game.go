package game

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Game interface {
	Join(string) int
	State() View
	Move(int, int) View
	View(int) View
}

type game struct {
	players    []string
	start      chan bool
	join       chan string
	joinID     chan int
	move       []chan int
	view       []chan View
	f          *field
	week       int
	deeds      map[int]*deed
	price      int
	surveyTurn int
}

type deed struct {
	player int
	week   int
	stop   int
	bit    int
	output int
	pnl    int
}

func New() Game {
	rand.Seed(time.Now().UTC().UnixNano())

	g := &game{
		f:      newField(),
		start:  make(chan bool),
		join:   make(chan string),
		joinID: make(chan int),
		deeds:  make(map[int]*deed),
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
	g.join <- name
	return <-g.joinID
}

func (g *game) Move(playerID, move int) View {
	if playerID == 0 && g.week == 0 {
		g.start <- true
	} else {
		g.move[playerID] <- move
	}
	return <-g.view[playerID]
}

// View returns a JSON serializable object representing the player's current game state.
func (g *game) View(playerID int) View {
	// noop returns the view without advancing the state
	go func() {
		g.move[playerID] <- noop
	}()
	return <-g.view[playerID]
}

// State returns the high-level state of the game: who has joined and has it started.
func (g *game) State() View {
	return view(g)
}

// game state machine func
type stateFn func(*game) stateFn

// lobby is the game state machine function for handling joins before the start of the game.
func lobby(g *game) stateFn {

Loop:
	for {
		select {
		case name := <-g.join:
			playerID := len(g.players)
			g.players = append(g.players, name)
			g.move = append(g.move, make(chan int))
			g.view = append(g.view, make(chan View))
			g.joinID <- playerID
			log.Printf("name %s joined as player %d", name, playerID)
		case <-g.start:
			break Loop
		}
	}
	close(g.join)

	log.Printf("starting game with %d players", len(g.players))
	g.nextWeek()

	return playWeek
}

// week is the game state machine function for handling a single week's gameplay.
func playWeek(g *game) stateFn {
	var wg sync.WaitGroup
	wg.Add(len(g.players))

	// run a state machine for each player in individual go routines
	for p := 0; p < len(g.players); p++ {
		go func(playerID int) {
			defer wg.Done()
			for state := survey; state != nil; {
				state = state(g, playerID)
			}
		}(p)
	}
	wg.Wait()

	log.Printf("all %d players completed week %d", len(g.players), g.week)
	g.nextWeek()

	if g.week == 13 {
		log.Println("game over!")
		return nil
	}

	return playWeek
}

func (g *game) nextWeek() {
	g.week++
	g.price = int(100 * math.Abs(1+rand.NormFloat64()))

	for s, d := range g.deeds {
		if d.bit == 0 || d.bit != g.f.oil[s] || d.stop > 0 {
			continue
		}

		// production considers reservoir pressure over time
		res := g.f.reservoir(s)
		tot := float64(len(res))
		for _, s := range res {
			d := g.deeds[s]
			if d == nil || d.bit == 0 || d.bit != g.f.oil[s] {
				continue
			}
			until := d.stop
			if until == 0 {
				until = g.week
			}
			// pressure diminishes 1/3 per pump site week. with a large enough
			// reservoir this is subtle but for a small reservoir it's devastating
			tot -= 1.0 - math.Pow(0.666, float64(until-d.week))
		}
		pressure := tot / float64(len(res))
		// ramp up: well capacity approaches 100 barrels per site @ 1.0 pressure
		capacity := 100 * (1 - math.Pow(0.5, float64(g.week-d.week)))
		output := int(math.Floor(pressure * capacity * float64(len(res))))
		log.Printf("reservoir %d size %d capacity %f pressure %f output %d", res, len(res), capacity, pressure, output)

		d.output = output
		d.pnl += int(float64(d.output*g.price)/100) - g.f.tax[s]
	}
}
