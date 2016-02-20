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

	// if update is a slice indexed by player move probably shoud be too..
	move   chan Move
	update []chan Update
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

type JoinUpdate struct {
	PlayerID int
}

type stateFn func(*game) stateFn

func New() Game {
	rand.Seed(time.Now().UTC().UnixNano())

	g := &game{
		f:     newField(),
		deeds: make(map[int]*deed),
		price: []int{85 + rand.Intn(30)}, // $0.85 - $1.15
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

type StartUpdate struct {
	Players []string `json:"players"`
}

func lobby(g *game) stateFn {
	// FIXME maybe we should go to the lobby every round (showing score summary) but giving a chance for last joins to come in...
	// "start" then really becomes "begin week" and i guess the game owner would be responsible for it.

	// wait for player zero to start the game
	for {
		mv := <-g.move

		if mv.PlayerID != 0 {
			log.Printf("Game has not started; ignoring non-owner player %d", mv.PlayerID)
			g.update[mv.PlayerID] <- nil
			continue
		}
		if !mv.Done {
			log.Println("Game has not started; owner must set done to start")
			g.update[mv.PlayerID] <- nil
			continue
		}
		log.Printf("Owner started game with %d players", len(g.players))
		g.week++
		g.update[mv.PlayerID] <- &StartUpdate{Players: g.players}
		break
	}

	return week
}

func week(g *game) stateFn {
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

	g.week++

	return week
}

type playerTurn func(g *game, move <-chan Move) playerTurn

type SurveyUpdate struct {
	Prob int `json:"prob"`
	Cost int `json:"cost"`
	Tax  int `json:"tax"`
}

func surveyTurn(g *game, move <-chan Move) playerTurn {
	for {
		mv := <-move

		if mv.PlayerID != g.turn {
			log.Printf("Waiting for player %d to survey; ignoring player %d", g.turn, mv.PlayerID)
			g.update[mv.PlayerID] <- nil
			continue
		}

		if _, ok := g.deeds[mv.SiteID]; ok {
			log.Printf("Site %d already surveyed; ignoring player %d", mv.SiteID, mv.PlayerID)
			g.update[mv.PlayerID] <- nil
			continue
		}

		log.Printf("Player %d surveying site %d", mv.PlayerID, mv.SiteID)
		g.deeds[mv.SiteID] = &deed{player: mv.PlayerID}
		g.turn = (g.turn + 1) % len(g.players)

		update := SurveyUpdate{
			Prob: g.f.p[mv.SiteID],
			Cost: g.f.cost[mv.SiteID],
			Tax:  g.f.tax[mv.SiteID],
		}
		g.update[mv.PlayerID] <- update

		return drillTurn(mv.SiteID)
	}
}

type DrillUpdate struct {
	Depth int  `json:"depth"`
	Cost  int  `json:"cost"`
	Oil   bool `json:"oil"`
}

func drillTurn(siteID int) playerTurn {
	return func(g *game, move <-chan Move) playerTurn {
		oil := g.f.oil[siteID]
		deed := g.deeds[siteID]
		for mv := range move {
			if mv.Done {
				log.Printf("Player %d done drilling site %d", mv.PlayerID, siteID)
				g.update[mv.PlayerID] <- DrillUpdate{Depth: deed.bit, Cost: deed.bit * g.f.cost[siteID]}
				break
			}

			log.Printf("Player %d drilling site %d with bit %d", mv.PlayerID, siteID, deed.bit)
			deed.start = g.week
			deed.bit++

			update := &DrillUpdate{Depth: deed.bit, Cost: deed.bit * g.f.cost[siteID]}

			if oil > 0 && deed.bit == oil {
				log.Printf("Player %d struck oil at depth %d", mv.PlayerID, deed.bit)
				update.Oil = true
				g.update[mv.PlayerID] <- update
				break
			}
			if deed.bit == maxOil {
				log.Printf("DRY HOLE for player %d", mv.PlayerID)
				g.update[mv.PlayerID] <- update
				break
			}
			g.update[mv.PlayerID] <- update
		}
		return sellTurn
	}
}

type SellUpdate struct {
	Cost int `json:"cost"`
}

func sellTurn(g *game, move <-chan Move) playerTurn {
	for mv := range move {
		if mv.Done {
			log.Printf("Player %d done selling", mv.PlayerID)
			g.update[mv.PlayerID] <- &SellUpdate{}
			break
		}
		deed := g.deeds[mv.SiteID]
		if deed == nil || deed.player != mv.PlayerID {
			log.Printf("Ignoring sale for site %d; player %d does not own deed", mv.SiteID, mv.PlayerID)
			g.update[mv.PlayerID] <- nil
			continue
		}
		if deed.stop > 0 {
			log.Printf("Ignoring sale for site %d; already sold in week %d", mv.SiteID, deed.stop)
			g.update[mv.PlayerID] <- nil
			continue
		}
		log.Printf("Player %d selling site %d", mv.PlayerID, mv.SiteID)
		g.deeds[mv.SiteID].stop = g.week
		g.update[mv.PlayerID] <- &SellUpdate{}
	}
	return nil
}
