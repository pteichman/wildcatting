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
	f          *field
	week       int
	deeds      map[int]*deed // site id to deed
	price      int           // oil price in cents
	players    []string      // id indexed player names
	surveyTurn int           // next survey turn (must happen in order)
	lobbyMove  chan int
	clientMove []chan int
	clientView []chan View
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
		f:         newField(),
		deeds:     make(map[int]*deed),
		lobbyMove: make(chan int),
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
	playerID := len(g.players)
	g.players = append(g.players, name)
	g.clientMove = append(g.clientMove, make(chan int))
	g.clientView = append(g.clientView, make(chan View))

	return playerID
}

func (g *game) Move(playerID, index int) View {
	// this is goofy but works for now. before the game is started
	// the lobby state is listening for moves on a single
	// game wide channel. this should be the game owner channel.
	if g.week == 0 {
		g.lobbyMove <- playerID
	} else {
		g.clientMove[playerID] <- index
	}
	return <-g.clientView[playerID]
}

// View returns a JSON serializable object representing the player's current game state.
func (g *game) View(playerID int) View {
	// noop returns the view without advancing the state
	go func() {
		g.clientMove[playerID] <- noop
	}()
	return <-g.clientView[playerID]
}

// State returns the high-level state of the game: who has joined and has it started.
func (g *game) State() View {
	return struct {
		Players []string `json:"players"`
		Started bool     `json:"started"`
	}{g.players, g.week > 0}
}

// game state machine func
type stateFn func(*game) stateFn

// lobby is the game state machine function for handling joins before the start of the game.
func lobby(g *game) stateFn {
	// wait for player zero to start the game
	for {
		playerID := <-g.lobbyMove

		if playerID != 0 {
			log.Printf("ignoring non-owner player %d", playerID)
			g.clientView[playerID] <- nil
			continue
		}

		log.Printf("started game with %d players", len(g.players))
		g.nextWeek()

		break
	}

	return playWeek
}

// week is the game state machine function for handling a single week's gameplay.
func playWeek(g *game) stateFn {
	var wg sync.WaitGroup
	wg.Add(len(g.players))

	// run a state machine for each player in individual go routines
	// var move []chan int
	for p := 0; p < len(g.players); p++ {
		// move = append(move, make(chan int))
		go func(playerID int) {
			for state := survey; state != nil; {
				state = state(g, playerID)
			}
			wg.Done()
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

	// Oil and gas wells usually reach their maximum output shortly after completion.
	// From that time, other than wells completed in water-drive reservoirs, they decline
	// in production, the rapidity of decline depending on the output of the wells and on
	// other factors governing their productivity. The production decline curve shows the
	// amount of oil and gas produced per unit of time for several consecutive periods;
	// if the conditions affecting the rate of production are not changed, the curve may
	// be fairly regular, and, if projected, will furnish useful knowledge as to the future
	// production of the well. By the aid of this knowledge the value of a property may be
	// judged, and proper depletion and depreciation charges may be made on the books of
	// the operating company.(Lewis 1918)
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

			// pressure diminishes 33% per pump site week. with a large enough
			// reservoir this is subtle but for a small reservoir it's devastating
			tot -= 1.0 - math.Pow(0.666, float64(until-d.week))
		}
		pressure := tot / float64(len(res))

		// FIXME aquafiers would need to be rolled into the previous loop,
		// counteracting the pressure decreases as pumping continue. this should be
		// tuned so without strong aquifiers, pressure reductions are devastating

		// The aquifer strength also refers to how well the aquifer mitigates the reservoir's
		// normal pressure decline. A strong aquifer refers to one in which the water-influx
		// rate approaches the reservoir's fluid withdrawal rate at reservoir conditions.

		// Reservoir engineers have often used pressure contour maps or some approximate
		// methods to determine field average reservoir pressure for p/z analysis.
		// Usually, however, individual well pressures are based on extrapolation of
		// pressure buildup tests or from long shut-in periods. In either case, the
		// average pressure measured does not represent a point value, but rather is
		// the average value within the well’s effective drainage volume

		// ramp up: capacity approaches 100 barrels per site @ 1.0 pressure
		capacity := 100 * (1 - math.Pow(0.5, float64(g.week-d.week)))
		output := int(math.Floor(pressure * capacity * float64(len(res))))
		log.Printf("reservoir %d size %d capacity %f pressure %f output %d", res, len(res), capacity, pressure, output)

		d.output = output
		d.pnl += int(float64(d.output*g.price)/100) - g.f.tax[s]
	}
}
