package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type game struct {
	f       *field
	week    int
	deeds   map[int]*deed // site id to ownership record
	price   []int         // oil price each week in cents
	players []string      // id indexed player names
	move    <-chan Move
	state   chan<- bool
}

type deed struct {
	player int
	start  int
	stop   int
	bit    int
}

type Move struct {
	PlayerID int  `json:"player"`
	SiteID   int  `json:"site"`
	Done     bool `json:"done"`
}

type stateFn func(*game) stateFn

func New(move <-chan Move) (Game, <-chan bool) {
	rand.Seed(time.Now().UTC().UnixNano())

	state := make(chan bool)

	g := &game{
		f:     newField(),
		deeds: make(map[int]*deed),
		price: []int{85 + rand.Intn(30)}, // $0.85 - $1.15
		move:  move,
		state: state,
	}

	go g.run()
	return g, state
}

func (g *game) run() {
	for state := lobby; state != nil; {
		state = state(g)

		select {
		case g.state <- true:
			fmt.Println("someone cares")
		default:
			fmt.Println("noone cares")
		}
	}
	close(g.state)
}

func lobby(g *game) stateFn {
	// FIXME maybe we should go to the lobby every round (showing score summary) but giving a chance for last joins to come in...
	// "start" then really becomes "begin week" and i guess the game owner would be responsible for it.

	// wait for player zero to start the game
	for {
		mv := <-g.move
		if mv.PlayerID == 0 && mv.Done {
			break
		}
	}

	fmt.Printf("started game with %d players\n", len(g.players))

	return survey
}

func survey(g *game) stateFn {
	g.week++
	g.price = append(g.price, g.price[len(g.price)-1])

	drills := make([]int, 4)

	for p := 0; p < len(g.players); p++ {
		for {
			mv := <-g.move
			if mv.PlayerID != p {
				log.Printf("waiting for player %d to survey; ignoring player %d", p, mv.PlayerID)
				continue
			}
			fmt.Printf("player %d surveying site %d\n", p, mv.SiteID)
			g.deeds[mv.SiteID] = &deed{player: p}
			drills[p] = mv.SiteID
			break
		}
	}
	// return a state transition based on the location of the drills
	return completeTurn(drills)
}

func completeTurn(drills []int) stateFn {
	return func(g *game) stateFn {
		var wg sync.WaitGroup
		wg.Add(len(g.players))

		move := make([]chan Move, len(g.players))

		for p := 0; p < len(g.players); p++ {
			move[p] = make(chan Move)
			go func(p int) {
				for mv := range move[p] {
					if mv.Done {
						fmt.Printf("player %d done drilling site %d\n", p, drills[p])
						break
					}
					fmt.Printf("player %d drilling site %d\n", p, drills[p])
					deed := g.deeds[drills[p]]
					deed.start = g.week
					deed.bit++
					if deed.bit == g.f.oil[drills[p]] || deed.bit == maxOil {
						break
					}
				}

				for mv := range move[p] {
					if mv.Done {
						fmt.Printf("player %d done selling\n", mv.PlayerID)
						break
					}
					fmt.Printf("player %d selling site %d\n", mv.PlayerID, mv.SiteID)
					g.deeds[mv.SiteID].stop = g.week
				}
				wg.Done()
			}(p)
		}

		done := make(chan struct{})
		go func() {
			for {
				select {
				case mv := <-g.move:
					move[mv.PlayerID] <- mv
				case <-done:
					return
				}
			}
		}()
		wg.Wait()
		close(done)

		// FIXME generate player game state here, right?
		// actually i think maybe it's on demand via Stater interface and the state channel goes away.
		// however, if there is long polling maybe there is an event channel from game that feeds it.

		return survey
	}
}