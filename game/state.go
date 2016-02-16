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
	PlayerID int
	SiteID   int
	Done     bool
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

		if mv.PlayerID != 0 {
			log.Printf("Game has not started; ignoring non-owner player %d", mv.PlayerID)
			continue
		}
		if !mv.Done {
			log.Println("Game has not started; owner must send Done=true")
			continue
		}
		break
	}

	log.Printf("Owner started game with %d players", len(g.players))

	return survey
}

func survey(g *game) stateFn {
	g.week++
	g.price = append(g.price, g.price[len(g.price)-1])

	drills := make([]int, 4)

	// eek realizing that his is unnecessarily synchronous...
	// while it is true that that the surveying must happen in order
	// once a player has surveyed they should immediately be able to
	// start drilling. waiting for all players to finish surveying
	// shouldn't be necessary. seems like completeTurn state is going
	// to go away which would bring us down to just two game states :/
	for p := 0; p < len(g.players); p++ {
		for {
			mv := <-g.move
			if mv.PlayerID != p {
				log.Printf("Waiting for player %d to survey; ignoring player %d", p, mv.PlayerID)
				continue
			}
			if _, ok := g.deeds[mv.SiteID]; ok {
				log.Printf("Site %d already surveyed; ignoring player %d", mv.SiteID, p)
				continue
			}
			log.Printf("Player %d surveying site %d", p, mv.SiteID)
			g.deeds[mv.SiteID] = &deed{player: p}
			drills[p] = mv.SiteID
			break
		}
	}
	log.Printf("All %d players finished surveying", len(g.players))
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
			go func(playerID int) {
				siteID := drills[p]
				oil := g.f.oil[siteID]
				deed := g.deeds[siteID]
				for mv := range move[p] {
					if mv.Done {
						log.Printf("Player %d done drilling site %d", playerID, siteID)
						break
					}

					log.Printf("Player %d drilling site %d with bit %d", playerID, siteID, deed.bit)
					deed.start = g.week
					deed.bit++

					if oil > 0 && deed.bit == oil {
						log.Printf("Player %d struck oil at depth %d", playerID, deed.bit)
						break
					}
					if deed.bit == maxOil {
						log.Printf("DRY HOLE for player %d", playerID)
						break
					}
				}

				for mv := range move[playerID] {
					if mv.Done {
						log.Printf("Player %d done selling", playerID)
						break
					}
					deed := g.deeds[mv.SiteID]
					if deed == nil || deed.player != playerID {
						log.Printf("Ignoring sale for site %d; player %d does not own deed", mv.SiteID, playerID)
						continue
					}
					if deed.stop > 0 {
						log.Printf("Ignoring sale for site %d; already sold in week %d", mv.SiteID, deed.stop)
						continue
					}
					log.Printf("Player %d selling site %d", playerID, mv.SiteID)
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
		log.Printf("All %d players completed week %d turn", len(g.players), g.week)

		// FIXME generate player game state here, right?
		// actually i think maybe it's on demand via Stater interface and the state channel goes away.
		// however, if there is long polling maybe there is an event channel from game that feeds it.

		return survey
	}
}
