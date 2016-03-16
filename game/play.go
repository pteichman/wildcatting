package game

import "log"

const (
	done = -1
	no   = 0
	yes  = 1
)

// a playFn represent's the players gameplay state within a week.
// calling the function transitions the state based on incoming
// player moves. it returns the playFn for the next state transition.
type playFn func(*game, entity) playFn

func survey(g *game, playerID entity) playFn {
	log.Printf("player %d survey state", playerID)
	var move int

Loop:
	for {
		select {
		case g.view[playerID] <- surveyView(g, playerID):
		case move = <-g.move[playerID]:
			if playerID != g.surveyTurn {
				log.Printf("waiting for player %d to survey; ignoring player %d", g.surveyTurn, playerID)
				break
			}

			if _, ok := g.deeds[move]; ok {
				log.Printf("site %d already surveyed; ignoring player %d", move, playerID)
				break
			}
			break Loop
		}
	}

	log.Printf("player %d surveying site %d", playerID, move)
	g.deeds[move] = &deed{player: playerID, week: g.week}
	g.surveyTurn = entity((int(g.surveyTurn) + 1) % len(g.players))

	return report(move)
}

func report(siteID int) playFn {
	// return this player's function for surveyor's report at specific site
	return func(g *game, playerID entity) playFn {
		log.Printf("player %d report state @ site %d", playerID, siteID)
		var move int
		for {
			select {
			case g.view[playerID] <- reportView(g, playerID, siteID):
			case move = <-g.move[playerID]:
				if move == no {
					return wells
				}
				if move == yes {
					return drill(siteID)
				}
				log.Printf("ignoring invalid report move from player %d move %d", playerID, move)
			}
		}
	}
}

func drill(siteID int) playFn {
	view := drillView(siteID)

	// return this player's function for drilling a specific site
	return func(g *game, playerID entity) playFn {
		log.Printf("player %d drill state @ site %d", playerID, siteID)
		oil := g.f.oil[siteID]
		deed := g.deeds[siteID]

	Loop:
		for {
			select {
			case g.view[playerID] <- view(g, playerID):
			case move := <-g.move[playerID]:
				if move == done {
					log.Printf("player %d done drilling site %d", playerID, siteID)
					break Loop
				}

				log.Printf("player %d drilling site %d with bit %d", playerID, siteID, deed.bit)
				deed.bit++
				deed.pnl -= g.f.cost[siteID]

				if deed.bit == oil || deed.bit == 9 {
					log.Printf("player %d done drilling site %d", playerID, siteID)
					break Loop
				}
			}
		}
		return wells
	}
}

func wells(g *game, playerID entity) playFn {
	log.Printf("player %d wells state", playerID)
Loop:
	for {
		select {
		case g.view[playerID] <- wellsView(g, playerID):
		case move := <-g.move[playerID]:

			if move == done {
				log.Printf("player %d done selling", playerID)
				break Loop
			}

			deed := g.deeds[move]
			if deed == nil || deed.player != playerID {
				log.Printf("ignoring sale for site %d; player %d does not own deed", move, playerID)
				break
			}
			if deed.stop > 0 {
				log.Printf("ignoring sale for site %d; already sold in week %d", move, deed.stop)
				break
			}
			log.Printf("player %d selling site %d", playerID, move)
			g.deeds[move].stop = g.week
		}
	}

	g.view[playerID] <- lobbyView(g)
	return nil
}
