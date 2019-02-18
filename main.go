package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	trueskill "github.com/mafredri/go-trueskill"
)

// Player represent a player entry.
type Player struct {
	Mu        float64 `json:"mu"`
	Sigma     float64 `json:"sigma"`
	TrueSkill float64 `json:"trueskill"` // Used for response.
}

// RateRequest represents the JSON request that defines the game setup and players.
type RateRequest struct {
	Mu       float64  `json:"mu"`
	Sigma    float64  `json:"sigma"`
	Beta     float64  `json:"beta"`
	Tau      float64  `json:"tau"`
	DrawProb *float64 `json:"draw_probability"`
	Players  []Player `json:"players"`
}

// RatedResponse represents the response from the rate endpoint containing all rated players and the probability of match outcome.
type RatedResponse struct {
	Players     []Player `json:"players"`
	Probability float64  `json:"probability_of_outcome"`
}

// Index shows that the trueskilld service is running.
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "trueskilld v0.1!\n")
}

// Rate takes a request to rate players and calculates their true skills.
func Rate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req RateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Bad request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if len(req.Players) < 2 {
		message := fmt.Sprintf("A minimum of 2 players must be provided")
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	var opts []trueskill.Option
	if req.Mu != 0 {
		opts = append(opts, trueskill.Mu(req.Mu))
	}
	if req.Sigma != 0 {
		opts = append(opts, trueskill.Sigma(req.Sigma))
	}
	if req.Beta != 0 {
		opts = append(opts, trueskill.Beta(req.Beta))
	}
	if req.Tau != 0 {
		opts = append(opts, trueskill.Tau(req.Tau))
	}
	if req.DrawProb != nil {
		drawProb, err := trueskill.DrawProbability(*req.DrawProb)
		if err != nil {
			message := fmt.Sprintf("Could not calculate draw probability: %v", err)
			log.Println(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
		opts = append(opts, drawProb)
	}

	ts := trueskill.New(opts...)

	var players []trueskill.Player
	for _, p := range req.Players {
		var player trueskill.Player
		if p.Mu == 0 && p.Sigma == 0 {
			player = ts.NewPlayer()
		} else {
			player = trueskill.NewPlayer(p.Mu, p.Sigma)
		}
		players = append(players, player)
	}

	adjustedPlayers, probability := ts.AdjustSkills(players, false)

	var newPlayers []Player
	for _, p := range adjustedPlayers {
		newPlayers = append(newPlayers, Player{
			Mu:        p.Mu(),
			Sigma:     p.Sigma(),
			TrueSkill: ts.TrueSkill(p),
		})
	}

	resp := RatedResponse{
		Players:     newPlayers,
		Probability: probability * 100,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		message := fmt.Sprintf("Could not marshal response: %v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func main() {
	bind := flag.String("bind", "127.0.0.1", "Bind to interface")
	port := flag.Int("port", 8495, "Listen on port")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/rate", Rate)

	addr := fmt.Sprintf("%s:%d", *bind, *port)
	log.Printf("Starting trueskilld on: '%s'", addr)
	log.Fatal(http.ListenAndServe(addr, logRequest(mux)))
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
