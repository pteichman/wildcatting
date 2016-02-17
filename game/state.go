package game

import (
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
	turn    int           // next surveying turn (which much happens in order)
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

	return week
}

func week(g *game) stateFn {
	g.week++
	g.price = append(g.price, g.price[len(g.price)-1])

	var wg sync.WaitGroup
	wg.Add(len(g.players))

	var move []chan Move
	for p := 0; p < len(g.players); p++ {
		move = append(move, make(chan Move))
		go func(playerID int) {
			for state := surveyTurn; state != nil; {
				state = state(g, move[playerID])
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
	log.Printf("All %d players completed week %d", len(g.players), g.week)

	if g.week == 13 {
		log.Println("Game over!")
		return nil
	}
	return week
}

type playerTurn func(g *game, move <-chan Move) playerTurn

func surveyTurn(g *game, move <-chan Move) playerTurn {
	for {
		mv := <-move

		if mv.PlayerID != g.turn {
			log.Printf("Waiting for player %d to survey; ignoring player %d", g.turn, mv.PlayerID)
			continue
		}

		if _, ok := g.deeds[mv.SiteID]; ok {
			log.Printf("Site %d already surveyed; ignoring player %d", mv.SiteID, mv.PlayerID)
			continue
		}

		log.Printf("Player %d surveying site %d", mv.PlayerID, mv.SiteID)
		g.deeds[mv.SiteID] = &deed{player: mv.PlayerID}
		g.turn = (g.turn + 1) % len(g.players)
		return drillTurn(mv.SiteID)
	}
}

func drillTurn(siteID int) playerTurn {
	return func(g *game, move <-chan Move) playerTurn {
		oil := g.f.oil[siteID]
		deed := g.deeds[siteID]
		for mv := range move {
			if mv.Done {
				log.Printf("Player %d done drilling site %d", mv.PlayerID, siteID)
				break
			}

			log.Printf("Player %d drilling site %d with bit %d", mv.PlayerID, siteID, deed.bit)
			deed.start = g.week
			deed.bit++

			if oil > 0 && deed.bit == oil {
				log.Printf("Player %d struck oil at depth %d", mv.PlayerID, deed.bit)
				break
			}
			if deed.bit == maxOil {
				log.Printf("DRY HOLE for player %d", mv.PlayerID)
				break
			}
		}
		return sellTurn
	}
}

func sellTurn(g *game, move <-chan Move) playerTurn {
	for mv := range move {
		if mv.Done {
			log.Printf("Player %d done selling", mv.PlayerID)
			break
		}
		deed := g.deeds[mv.SiteID]
		if deed == nil || deed.player != mv.PlayerID {
			log.Printf("Ignoring sale for site %d; player %d does not own deed", mv.SiteID, mv.PlayerID)
			continue
		}
		if deed.stop > 0 {
			log.Printf("Ignoring sale for site %d; already sold in week %d", mv.SiteID, deed.stop)
			continue
		}
		log.Printf("Player %d selling site %d", mv.PlayerID, mv.SiteID)
		g.deeds[mv.SiteID].stop = g.week
	}
	return nil
}
