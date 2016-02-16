package game

import "math"

type Game interface {
	Join(string) int
	State(int) *PlayerState
}

func (g *game) Join(name string) int {
	// we probably need to lock the players slice
	p := len(g.players)
	g.players = append(g.players, name)

	return p
}

func (g *game) State(playerID int) *PlayerState {
	var sites []PlayerSite
	for s, deed := range g.deeds {
		ps := PlayerSite{
			ID:    s,
			Owner: deed.player,
			P:     g.f.p[s],
			Cost:  g.f.cost[s],
			Tax:   g.f.tax[s],
		}
		// players only know about oil if it was reached with the bit
		if deed.bit > 0 && deed.bit == g.f.oil[s] {
			ps.Oil = g.f.oil[s]
		}

		sites = append(sites, ps)
	}

	// var waiting []int
	// for p, play := range g.play {
	// 	if play != nil {
	// 		waiting = append(waiting, p)
	// 	}
	// }

	wellRev := make(map[int]int)
	for s, deed := range g.deeds {
		if deed.bit == 0 || deed.bit != g.f.oil[s] {
			continue
		}

		// costs
		cost := deed.bit * g.f.cost[s]
		if deed.stop > 0 {
			wellRev[s] += cost / 2
		}
		wellRev[s] -= cost

		// production considers reservoir pressure over time
		var res []int
		for r := range g.f.reservoir(s) {
			res = append(res, r)
		}
		// consider the pressure levels of this reservoir over time
		pres := g.pressure(res)
		for w := 0; w <= g.week; w++ {
			barrels := int(math.Floor(pres[w] * 10.0 * float64(len(res))))
			wellRev[s] += barrels * g.price[w]
			wellRev[s] -= g.f.tax[s]
		}
	}

	revenue := make([]int, len(g.players))
	for s, rev := range wellRev {
		revenue[g.deeds[s].player] += rev
	}

	return &PlayerState{
		Players: g.players,
		Sites:   sites,
		Revenue: revenue,
		// Waiting: waiting,
	}
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
func (g *game) pressure(res []int) []float64 {
	pressure := []float64{1.0}

	tot := float64(len(res))
	for w := 1; w <= g.week; w++ {
		// as pumping proceeds reservoir pressure is adversely affected
		neg := 0.0
		for r := range res {
			d := g.deeds[r]
			if d == nil || d.bit == 0 || d.bit != g.f.oil[r] {
				continue
			}
			stop := g.deeds[r].stop
			if stop == 0 {
				stop = g.week
			}
			if w > g.deeds[r].start && w <= stop {
				// every pumpweek subtracts 1/3 of the pressure at its site
				neg += 0.333
			}
		}
		tot -= neg
		pressure = append(pressure, tot/float64(len(res)))
	}

	// FIXME aquifiers would need to be rolled into the previous loop,
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

	return pressure
}

// FIXME maybe we can put json tags on "game" (becomes GameState) and serialize that
// although we'd have to have a function which copies the available state for a given player
// so not sure if a new GameState simplifies anything
type PlayerSite struct {
	ID    int `json:"id"`
	P     int `json:"p"`
	Cost  int `json:"cost"`
	Oil   int `json:"oil"`
	Tax   int `json:"tax"`
	Owner int `json:"owner"`
}

type PlayerState struct {
	Players []string     `json:"players"`
	Sites   []PlayerSite `json:"sites"`
	Revenue []int        `json:"revenue"`
	Waiting []int        `json:"waiting"`
}
