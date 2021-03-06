package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/9r33n/wildcatting/game"
	"github.com/gorilla/mux"
)

var (
	debug = flag.String("debug", "", "run expvar/pprof server (host:port)")
	stats = expvar.NewMap("wildcatting")
)

type handler struct {
	games []game.Game
}

func main() {
	flag.Parse()
	publishRuntime()

	if *debug != "" {
		go func() {
			log.Fatal(http.ListenAndServe(*debug, nil))
		}()
	}

	h := &handler{}

	host := "0.0.0.0"
	port := 8888
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("HTTP server listening to %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, statswrap(h.newRouter())))
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
		route{"GET", "/game/{gid:[0-9]+}/", h.getGameID},
		route{"POST", "/game/{gid:[0-9]+}/player/{pid:[0-9]+}/", h.postPlayerID},
		route{"GET", "/game/{gid:[0-9]+}/player/{pid:[0-9]+}/", h.getPlayerID},
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	for _, route := range routes {
		r.
			Methods(route.Method).
			Path(route.Path).
			Handler(route.Handler)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/")))

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

func (h *handler) postGameID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, err := strconv.Atoi(vars["gid"])
	if err != nil {
		// mux should guarantee a parsable int
		panic(err)
	}

	g := h.games[gameID]

	var name string
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&name)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	playerID := g.Join(name)
	if _, err := w.Write([]byte(fmt.Sprintf("%d", playerID))); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("\"%s\" joined game %d as player %d", name, gameID, playerID)
}

func (h *handler) getGameID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, err := strconv.Atoi(vars["gid"])
	if err != nil {
		// mux should guarantee a parsable int
		panic(err)
	}
	update := h.games[gameID].Status()
	js, err := json.Marshal(update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// move making... starting, surveying, drilling, selling, scoring
func (h *handler) postPlayerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, err := strconv.Atoi(vars["gid"])
	if err != nil {
		panic(err)
	}
	playerID, err := strconv.Atoi(vars["pid"])
	if err != nil {
		panic(err)
	}

	var mv int
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&mv)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	g := h.games[gameID]
	update := g.Move(playerID, mv)

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
	gameID, err := strconv.Atoi(vars["gid"])
	if err != nil {
		panic(err)
	}
	playerID, err := strconv.Atoi(vars["pid"])
	if err != nil {
		panic(err)
	}

	g := h.games[gameID]
	if g == nil {
		http.Error(w, "game not found", http.StatusBadRequest)
		return
	}

	state := g.View(playerID)

	js, err := json.Marshal(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func publishRuntime() {
	expvar.Publish("NumGoroutine", expvar.Func(
		func() interface{} { return runtime.NumGoroutine() },
	))

	start := time.Now().UnixNano()
	expvar.Publish("Uptime", expvar.Func(
		func() interface{} { return time.Now().UnixNano() - start },
	))
}

func statswrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		stats.Add("Requested", 1)
	})
}
