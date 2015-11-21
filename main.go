package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mafredri/go-trueskill"
)

// RoundedFloat64 provides a JSON encodable rounded float64
type RoundedFloat64 float64

// MarshalJSON marhsals a RoundedFloat64
func (f RoundedFloat64) MarshalJSON() ([]byte, error) {
	rounded := float64(math.Floor(float64(f)*1000+.5)) / 1000
	return json.Marshal(rounded)
}

// PlayerRequest .
type PlayerRequest struct {
	Mu    float64 `json:"mu"`
	Sigma float64 `json:"sigma"`
}

// RateRequest .
type RateRequest struct {
	Mu       float64         `json:"mu"`
	Sigma    float64         `json:"sigma"`
	Beta     float64         `json:"beta"`
	Tau      float64         `json:"tau"`
	DrawProb float64         `json:"draw_probability"`
	Players  []PlayerRequest `json:"players"`
}

// PlayerResponse .
type PlayerResponse struct {
	Mu        RoundedFloat64 `json:"mu"`
	Sigma     RoundedFloat64 `json:"sigma"`
	TrueSkill int64          `json:"trueskill"`
}

// RatedResponse .
type RatedResponse struct {
	Players     []PlayerResponse `json:"players"`
	Probability RoundedFloat64   `json:"probability_of_outcome"`
}

// Index .
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "trueskilld v0.1!\n")
}

// Rate .
func Rate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var req RateRequest
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}

	var ts trueskill.Config
	if req.Mu == 0 {
		ts, err = trueskill.NewDefault(req.DrawProb)
	} else {
		ts, err = trueskill.New(req.Mu, req.Sigma, req.Beta, req.Tau, req.DrawProb)
	}
	if err != nil {
		panic(err)
	}

	var players trueskill.Players
	for _, p := range req.Players {
		var player trueskill.Player
		if p.Mu == 0 && p.Sigma == 0 {
			player = ts.NewDefaultPlayer()
		} else {
			player = trueskill.NewPlayer(p.Mu, p.Sigma)
		}
		players = append(players, player)
	}

	adjustedPlayers, probability := ts.AdjustSkills(players, false)

	var newPlayers []PlayerResponse
	for _, p := range adjustedPlayers {
		newPlayers = append(newPlayers, PlayerResponse{
			Mu:        RoundedFloat64(p.Mu()),
			Sigma:     RoundedFloat64(p.Sigma()),
			TrueSkill: ts.TrueSkill(p),
		})
	}

	resp := RatedResponse{
		Players:     newPlayers,
		Probability: RoundedFloat64(probability * 100),
	}

	str, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")

	fmt.Fprintf(w, string(str))
}

func main() {
	port := flag.Int("p", 8495, "Port to run server on")
	flag.Parse()

	router := httprouter.New()

	router.GET("/", Index)
	router.POST("/rate", Rate)

	listen := fmt.Sprintf(":%d", *port)
	log.Printf("Starting trueskilld on: '%s'", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}