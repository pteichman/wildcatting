package game

import (
	"expvar"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

var stats = expvar.NewMap("game")

type Game interface {
	Join(string) int
	Status() View
	Move(int, int) View
	View(int) View
}

type game struct {
	players    []string
	join       chan string
	joinID     chan int
	move       []chan int
	status     chan View
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
		f:      newField(24, 80),
		join:   make(chan string),
		joinID: make(chan int),
		status: make(chan View),
		deeds:  make(map[int]*deed),
	}

	go g.run()

	stats.Add("Created", 1)
	return g
}

func (g *game) run() {
	for state := lobby; state != nil; {
		state = state(g)
	}
}

func (g *game) Join(name string) int {
	stats.Add("Joined", 1)
	g.join <- name
	return <-g.joinID
}

func (g *game) Move(playerID, move int) View {
	stats.Add("Moved", 1)
	g.move[playerID] <- move
	return <-g.view[playerID]
}

// View returns a JSON serializable object representing the player's current game state.
func (g *game) View(playerID int) View {
	stats.Add("Viewed", 1)
	return <-g.view[playerID]
}

// State returns the high-level state of the game: who has joined and has it started.
func (g *game) Status() View {
	return <-g.status
}

// game state machine func
type stateFn func(*game) stateFn

// lobby is the game state machine function for handling joins before the start of the game.
func lobby(g *game) stateFn {
	var start chan int

	stop := make(chan struct{})
	go func() {
	Loop:
		for {
			select {
			case g.status <- lobbyView(g):
			case <-stop:
				break Loop
			}
		}
	}()

Loop:
	for {
		// player 0 is the owner and her first move is the start signal
		if len(g.move) > 0 {
			start = g.move[0]
		}
		select {
		case name := <-g.join:
			playerID := len(g.players)
			g.players = append(g.players, name)
			g.move = append(g.move, make(chan int))
			g.view = append(g.view, make(chan View))
			g.joinID <- playerID
			log.Printf("name %s joined as player %d", name, playerID)
		case <-start:
			break Loop
		}
	}
	close(stop)

	log.Printf("starting week with %d players", len(g.players))
	g.nextWeek()

	return play
}

// week is the game state machine function for handling a single week's gameplay.
func play(g *game) stateFn {
	stop := make(chan struct{})
	go func() {
	Loop:
		for {
			select {
			case g.status <- playView(g):
			case <-stop:
				break Loop
			}
		}
	}()

	// run a state machine for each player in individual go routines
	var wg sync.WaitGroup
	wg.Add(len(g.players))
	for p := 0; p < len(g.players); p++ {
		go func(playerID int) {
			defer wg.Done()
			for state := survey; state != nil; {
				state = state(g, playerID)
			}
		}(p)
	}
	wg.Wait()
	close(stop)

	log.Printf("all %d players completed week %d", len(g.players), g.week)

	return lobby
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
