package game

import (
	"log"
	"math/rand"
	"sync"
)

// JoinUpdate is a JSON serializable update containing the result of a join Move.
type JoinUpdate struct {
	PlayerID int
}

// StartUpdate is a JSON serializable Update containing the result of a game start Move.
type StartUpdate struct {
	Week int    `json:"week"`
	Prob []int  `json:"prob"`
	Cost []int  `json:"cost"`
	Tax  []int  `json:"tax"`
	Oil  []int  `json:"oil"`
	Fact string `json:"fact"`
}

// SurveyUpdate is a JSON serializable Update containing the result of a survey Move.
type SurveyUpdate struct {
	Site int `json:"site"`
	Prob int `json:"prob"`
	Cost int `json:"cost"`
	Tax  int `json:"tax"`
}

// DrillUpdate is a JSON serializable Update containing the result of a drill Move.
type DrillUpdate struct {
	Depth int  `json:"depth"`
	Cost  int  `json:"cost"`
	Oil   bool `json:"oil"`
}

// SellUpdate is a JSON serializable Update container the result of a sell Move.
type SellUpdate struct {
	Cost int `json:"cost"`
}

// game state machine func
type stateFn func(*game) stateFn

// player turn state machine func
type playerFn func(g *game, move <-chan Move) playerFn

// lobby is the game state machine function for handling
// joins before the start of the game.
func lobby(g *game) stateFn {
	// FIXME maybe we should go to the lobby every round (showing score summary)
	// giving a chance for late joins to come in... "start" then becomes
	// "begin week" and i guess the game owner would be responsible for it

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
		update := &StartUpdate{
			Week: g.week,
			Prob: g.f.prob,
			Cost: g.f.cost,
			Tax:  g.f.tax,
			Oil:  g.f.oil,
			Fact: facts[rand.Intn(len(facts))],
		}
		g.update[mv.PlayerID] <- update
		break
	}

	return week
}

// week is the game state machine function for handling a single week's gameplay.
func week(g *game) stateFn {
	g.price = append(g.price, g.price[len(g.price)-1])

	var wg sync.WaitGroup
	wg.Add(len(g.players))

	// run a state machine for each player in individual go routines
	var move []chan Move
	for p := 0; p < len(g.players); p++ {
		move = append(move, make(chan Move))
		go func(playerID int) {
			for state := survey; state != nil; {
				state = state(g, move[playerID])
			}
			wg.Done()
		}(p)
	}

	// direct incoming Moves to a player specific channel
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

// survey is a player state machine function for handling player survey moves.
func survey(g *game, move <-chan Move) playerFn {
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
			Site: mv.SiteID,
			Prob: g.f.prob[mv.SiteID],
			Cost: g.f.cost[mv.SiteID],
			Tax:  g.f.tax[mv.SiteID],
		}
		g.update[mv.PlayerID] <- update

		return drillSite(mv.SiteID)
	}
}

// drill returns a player state machine function for drilling a specific site
func drillSite(siteID int) playerFn {
	return func(g *game, move <-chan Move) playerFn {
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
		return sell
	}
}

// sell is a player state machine function for handling sales
// of wells before the end of the turn
func sell(g *game, move <-chan Move) playerFn {
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
