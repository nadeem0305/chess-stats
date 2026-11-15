// Package handler is the Vercel serverless function behind the web demo: GET /api/stats?user=NAME&months=N
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/nadeem0305/chess-stats/stats"
)

var validUser = regexp.MustCompile(`^[A-Za-z0-9_-]{3,25}$`)

func Handler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if !validUser.MatchString(user) {
		http.Error(w, "invalid username", http.StatusBadRequest)
		return
	}
	months, _ := strconv.Atoi(r.URL.Query().Get("months"))
	months = min(max(months, 1), 6) // keep each request to a handful of chess.com calls

	s, err := stats.Fetch(user, months)
	if errors.Is(err, stats.ErrNotFound) {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "chess.com lookup failed", http.StatusBadGateway)
		return
	}
	s.Openings = s.Openings[:min(10, len(s.Openings))]
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, s-maxage=600") // games change slowly; spare chess.com
	json.NewEncoder(w).Encode(s)
}
