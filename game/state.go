package game

import (
	"log"
	"math"
	"math/rand"
	"sync"
)

type State struct {
	Players []string `json:"players"`
	Week    int      `json:"week"`
	Prob    []int    `json:"prob"`
	Cost    []int    `json:"cost"`
	Tax     []int    `json:"tax"`
	Oil     []int    `json:"oil"`
	Deeds   []Deed   `json:"deeds"`
	Fact    string   `json:"fact"`
}

type Deed struct {
	SiteID int  `json:"site"`
	Owner  int  `json:"owner"`
	Sold   bool `json:"sold"`
	Oil    int  `json:"oil"`
	Cost   int  `json:"cost"`
	Income int  `json:"income"`
	PNL    int  `json:"pnl"`
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

		update := g.View(mv.PlayerID)
		g.update[mv.PlayerID] <- update
		break
	}

	return week
}

// week is the game state machine function for handling a single week's gameplay.
func week(g *game) stateFn {
	g.price = int(100 * math.Abs(1+rand.NormFloat64()))

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

	for s, deed := range g.deeds {
		if deed.bit == 0 || deed.bit != g.f.oil[s] || deed.stop > 0 {
			continue
		}

		// production considers reservoir pressure over time
		var res []int
		for r := range g.f.reservoir(s) {
			res = append(res, r)
		}
		// consider the pressure levels of this reservoir over time
		pressure := g.pressure(res)
		// 100 barrels per site @ 1.0 pressure
		deed.barrels = int(math.Floor(pressure * 100.0 * float64(len(res))))
		gross := deed.barrels * g.price / 100
		log.Printf("Well at site %d @ reservoir %d is outputting %d barrels with reservoir pressure %f with oil price %d grossing %d",
			s, res, deed.barrels, pressure, g.price, gross)
		deed.pnl += gross - g.f.tax[s]
	}

	if g.week == 13 {
		log.Println("Game over!")
		return nil
	}

	g.week++

	return week
}

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
func (g *game) pressure(res []int) float64 {
	tot := float64(len(res))

	for _, s := range res {
		d := g.deeds[s]
		if d == nil || d.bit == 0 || d.bit != g.f.oil[s] {
			continue
		}

		until := g.deeds[s].stop
		if until == 0 {
			until = g.week
		}

		dim := (1.0 - math.Pow(0.66, float64(until-g.deeds[s].start)))
		log.Printf("well at site %d has been active %d weeks; diminishing pressue %f", s, until-g.deeds[s].start, dim)

		// production diminishes 33% per pump site week
		tot -= dim
	}
	pressure := tot / float64(len(res))

	log.Printf("Total reservoir pressure for %d is %f", res, pressure)
	return pressure

	// FIXME aquafiers would need to be rolled into the previous loop,
	// counteracting the pressure decreases as pumping continue. this should be
	// tuned so without strong aquifiers, pressure reductions are devastating

	// The aquifer strength also refers to how well the aquifer mitigates the reservoir's
	// normal pressure decline. A strong aquifer refers to one in which the water-influx
	// rate approaches the reservoir's fluid withdrawal rate at reservoir conditions.

	// Reservoir engineers have often used pressure contour maps or some approximate
	//  methods to determine field average reservoir pressure for p/z analysis.
	// Usually, however, individual well pressures are based on extrapolation of
	// pressure buildup tests or from long shut-in periods. In either case, the
	// average pressure measured does not represent a point value, but rather is
	// the average value within the well’s effective drainage volume
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

		g.update[mv.PlayerID] <- g.View(mv.PlayerID)

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
				g.update[mv.PlayerID] <- g.View(mv.PlayerID)
				break
			}

			log.Printf("Player %d drilling site %d with bit %d", mv.PlayerID, siteID, deed.bit)
			deed.start = g.week
			deed.bit++

			if oil > 0 && deed.bit == oil {
				log.Printf("Player %d struck oil at depth %d", mv.PlayerID, deed.bit)
				g.update[mv.PlayerID] <- g.View(mv.PlayerID)
				break
			}
			if deed.bit == maxOil {
				log.Printf("DRY HOLE for player %d", mv.PlayerID)
				g.update[mv.PlayerID] <- g.View(mv.PlayerID)
				break
			}
			g.update[mv.PlayerID] <- g.View(mv.PlayerID)
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
			g.update[mv.PlayerID] <- g.View(mv.PlayerID)
			break
		}
		deed := g.deeds[mv.SiteID]
		if deed == nil || deed.player != mv.PlayerID {
			log.Printf("Ignoring sale for site %d; player %d does not own deed", mv.SiteID, mv.PlayerID)
			g.update[mv.PlayerID] <- g.View(mv.PlayerID)
			continue
		}
		if deed.stop > 0 {
			log.Printf("Ignoring sale for site %d; already sold in week %d", mv.SiteID, deed.stop)
			g.update[mv.PlayerID] <- g.View(mv.PlayerID)
			continue
		}
		log.Printf("Player %d selling site %d", mv.PlayerID, mv.SiteID)
		g.deeds[mv.SiteID].stop = g.week
		g.update[mv.PlayerID] <- g.View(mv.PlayerID)
	}
	return nil
}
