package main

import (
	"boomtown/game"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type handler struct {
	games []game.Game
	move  []chan<- game.Move
}

type Move struct {
	SiteID int  `json:"site"`
	Done   bool `json:"done"`
}

func main() {
	h := &handler{}

	host := "0.0.0.0"
	port := 8888
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("HTTP server listening to %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, h.newRouter()))
}

func (h *handler) newRouter() *mux.Router {
	type route struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}

	var routes = []route{
		route{"POST", "/game/", h.postGame},
		route{"POST", "/game/{gid:[0-9]+}/", h.postGameID},
		route{"POST", "/game/{gid:[0-9]+}/player/{pid:[0-9]}/", h.postPlayerID},
		route{"GET", "/game/{gid:[0-9]+}/player/{pid:[0-9]}/", h.getPlayerID},
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	for _, route := range routes {
		r.
			Methods(route.Method).
			Path(route.Path).
			Handler(route.Handler)
	}
	return r
}

// create game
func (h *handler) postGame(w http.ResponseWriter, r *http.Request) {
	gameID := len(h.games)

	g := game.New()
	h.games = append(h.games, g)

	if _, err := w.Write([]byte(fmt.Sprintf("%d", gameID))); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Created game", gameID)
}

// join game
func (h *handler) postGameID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, _ := strconv.Atoi(vars["gid"])

	var name string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&name)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	playerID := h.games[gameID].Join(name)
	if _, err := w.Write([]byte(fmt.Sprintf("%d", playerID))); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("\"%s\" joined game %d as player %d", name, gameID, playerID)
}

// move making... starting, surveying, drilling, selling
//
// start -> week
// survey -> probability, tax, cost
// drill -> bit, wet
// selling -> well revenue
func (h *handler) postPlayerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, _ := strconv.Atoi(vars["gid"])
	playerID, _ := strconv.Atoi(vars["pid"])

	var mv Move
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mv)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	g := h.games[gameID]
	update := g.Move(game.Move{playerID, mv.SiteID, mv.Done})

	js, err := json.Marshal(update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *handler) getPlayerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, _ := strconv.Atoi(vars["gid"])
	playerID, _ := strconv.Atoi(vars["pid"])

	g := h.games[gameID]
	if g == nil {
		http.Error(w, "game not found", http.StatusBadRequest)
		return
	}

	state := g.State(playerID)

	js, err := json.Marshal(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
