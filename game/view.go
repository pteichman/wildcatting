package game

import "math/rand"

// View is a generic type for JSON serializable data representing the client state.
type View interface{}

func lobbyView(g *game) View {
	type player struct {
		Name string `json:"name"`
		PNL  int    `json:"pnl"`
	}

	players := make([]player, len(g.players))
	for p := range g.players {
		pnl := 0
		for _, deed := range g.deeds {
			if deed.player != entity(p) {
				continue
			}
			pnl += deed.pnl
		}
		players[p] = player{g.players[p], pnl}
	}

	return struct {
		Name    string   `json:"name"`
		Week    int      `json:"week"`
		Players []player `json:"players"`
	}{"lobby", g.week, players}
}

func playView(g *game) View {
	return struct {
		Name string `json:"name"`
		Week int    `json:"week"`
	}{"play", g.week}
}

type playerViewFn func(*game, entity) View

func surveyView(g *game, playerID entity) View {
	return struct {
		Name  string `json:"name"`
		Week  int    `json:"week"`
		Price int    `json:"price"`
		Prob  []int  `json:"prob"`
		Cost  []int  `json:"cost"`
		Tax   []int  `json:"tax"`
		Oil   []int  `json:"oil"`
		Fact  string `json:"fact"`
	}{"survey", g.week, g.price, g.f.prob, g.f.cost, g.f.tax, g.f.oil, facts[rand.Intn(len(facts))]}
}

func reportView(g *game, playerID entity, siteID site) View {
	return struct {
		Name string `json:"name"`
		Site site   `json:"site"`
		Prob int    `json:"prob"`
		Cost int    `json:"cost"`
		Tax  int    `json:"tax"`
	}{"report", siteID, g.f.prob[siteID], g.f.cost[siteID], g.f.tax[siteID]}
}

func drillView(siteID site) playerViewFn {
	return func(g *game, playerID entity) View {
		depth := g.deeds[siteID].bit * 100
		cost := g.deeds[siteID].bit * g.f.cost[siteID]
		return struct {
			Name  string `json:"name"`
			Depth int    `json:"depth"`
			Cost  int    `json:"cost"`
		}{"drill", depth, cost}
	}
}

type well struct {
	Week   int  `json:"week"`
	SiteID site `json:"site"`
	Sold   bool `json:"sold"`
	Depth  int  `json:"depth"`
	Cost   int  `json:"cost"`
	Tax    int  `json:"tax"`
	Income int  `json:"income"`
	PNL    int  `json:"pnl"`
}

func wellsView(g *game, playerID entity) View {
	wells := make([]well, g.week)
	for s, deed := range g.deeds {
		if deed.player != playerID {
			continue
		}

		var tax int
		if deed.bit > 0 {
			tax = g.f.tax[s]
		}

		well := well{
			Week:   deed.week,
			SiteID: s,
			Sold:   deed.stop > 0,
			Cost:   g.f.cost[s] * deed.bit, // cost is in cents and bit is in 100 ft increments so they cancel out
			Tax:    tax,
			Income: deed.output * g.price / 100,
			PNL:    deed.pnl,
		}
		// players only know about oil if it was reached with the bit
		if deed.bit == g.f.oil[s] {
			well.Depth = g.f.oil[s] * 100
		}
		wells[deed.week-1] = well
	}

	state := struct {
		Name   string `json:"name"`
		Player string `json:"player"`
		Week   int    `json:"week"`
		Price  int    `json:"price"`
		Wells  []well `json:"wells"`
	}{"wells", g.players[playerID], g.week, g.price, wells}
	return state
}

func scoreView(g *game, playerID entity) View {
	return struct {
		Name string `json:"name"`
	}{"score"}
}
