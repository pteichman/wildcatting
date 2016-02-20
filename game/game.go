package game

import (
	"math"
	"math/rand"
	"time"
)

type Game interface {
	Join(string) int
	Move(Move) Update
	View(int) *View
}

// Update is a JSON serializable update containing the result to a player Move.
type Update interface{}

type Move struct {
	PlayerID int
	SiteID   int
	Done     bool
}

type View struct {
	Week    int      `json:"week"`
	Players []string `json:"players"`
	Deeds   []Deed   `json:"deeds"`
	Revenue []int    `json:"revenue"`
}

type Deed struct {
	SiteID int `json:"site"`
	Prob   int `json:"prob"`
	Cost   int `json:"cost"`
	Oil    int `json:"oil"`
	Tax    int `json:"tax"`
	Owner  int `json:"owner"`
}

type game struct {
	f       *field
	week    int
	deeds   map[int]*deed // site id to ownership record
	price   []int         // oil price each week in cents
	players []string      // id indexed player names
	turn    int           // next surveying turn (which much happens in order)
	move    chan Move
	update  []chan Update
}

type deed struct {
	player int
	start  int
	stop   int
	bit    int
}

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

func (g *game) Join(name string) int {
	// we probably need to lock the players slice
	p := len(g.players)
	g.players = append(g.players, name)

	g.update = append(g.update, make(chan Update))

	return p
}

func (g *game) Move(mv Move) Update {
	g.move <- mv
	return <-g.update[mv.PlayerID]
}

// View returns a players' View of the oil field state.
func (g *game) View(playerID int) *View {
	var deeds []Deed
	for s, deed := range g.deeds {
		ps := Deed{
			SiteID: s,
			Owner:  deed.player,
			Prob:   g.f.p[s],
			Cost:   g.f.cost[s],
			Tax:    g.f.tax[s],
		}
		// players only know about oil if it was reached with the bit
		if deed.bit > 0 && deed.bit == g.f.oil[s] {
			ps.Oil = g.f.oil[s]
		}

		deeds = append(deeds, ps)
	}

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

	return &View{
		Players: g.players,
		Week:    g.week,
		Deeds:   deeds,
		Revenue: revenue,
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

	return pressure
}
