package game

import "log"

const (
	done = -1
	noop = -2
)

const (
	no  = iota
	yes = iota
)

// a playFn represent's the players gameplay state within a week.
// calling the function transitions the state based on incoming
// player moves. it returns the playFn for the next state transition.
type playFn func(*game, int) playFn

func survey(g *game, playerID int) playFn {
	for {
		g.view[playerID] <- surveyView(g, playerID)

		mv := <-g.move[playerID]

		if mv == noop {
			continue
		}

		if playerID != g.surveyTurn {
			log.Printf("waiting for player %d to survey; ignoring player %d", g.surveyTurn, playerID)
			continue
		}

		if _, ok := g.deeds[mv]; ok {
			log.Printf("site %d already surveyed; ignoring player %d", mv, playerID)
			continue
		}

		log.Printf("player %d surveying site %d", playerID, mv)
		g.deeds[mv] = &deed{player: playerID, week: g.week}
		g.surveyTurn = (g.surveyTurn + 1) % len(g.players)

		return report(mv)
	}
}

func report(siteID int) playFn {
	// return this player's function for surveyor's report at specific site
	return func(g *game, playerID int) playFn {
		log.Print("report state player", playerID)
		for {
			g.view[playerID] <- reportView(g, playerID, siteID)
			mv := <-g.move[playerID]

			if mv == noop {
				continue
			}
			if mv == yes {
				return drill(siteID)
			}
			return wells
		}
	}
}

func drill(siteID int) playFn {
	view := drillView(siteID)

	// return this player's function for drilling a specific site
	return func(g *game, playerID int) playFn {
		oil := g.f.oil[siteID]
		deed := g.deeds[siteID]

		for {
			g.view[playerID] <- view(g, playerID)
			mv := <-g.move[playerID]

			if mv == noop {
				continue
			}
			if mv == done {
				log.Printf("player %d done drilling site %d", playerID, siteID)
				break
			}
			log.Printf("player %d drilling site %d with bit %d", playerID, siteID, deed.bit)
			deed.bit++

			if deed.bit == oil || deed.bit == 9 {
				log.Printf("player %d done drilling site %d", playerID, siteID)
				break
			}
		}
		return wells
	}
}

func wells(g *game, playerID int) playFn {
	for {
		g.view[playerID] <- wellsView(g, playerID)
		mv := <-g.move[playerID]

		if mv == noop {
			continue
		}
		if mv == done {
			log.Printf("player %d done selling", playerID)
			break
		}

		deed := g.deeds[mv]
		if deed == nil || deed.player != playerID {
			log.Printf("ignoring sale for site %d; player %d does not own deed", mv, playerID)
			continue
		}
		if deed.stop > 0 {
			log.Printf("ignoring sale for site %d; already sold in week %d", mv, deed.stop)
			continue
		}
		log.Printf("player %d selling site %d", playerID, mv)
		g.deeds[mv].stop = g.week
	}

	return score
}

func score(g *game, playerID int) playFn {
	g.view[playerID] <- scoreView(g, playerID)
	<-g.move[playerID]
	return nil
}
