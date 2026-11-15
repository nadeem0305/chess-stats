// Package stats fetches a chess.com player's recent games and summarises the results.
package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Side struct {
	Username string `json:"username"`
	Result   string `json:"result"`
}

type Game struct {
	TimeClass string `json:"time_class"`
	Rules     string `json:"rules"`
	ECO       string `json:"eco"`
	White     Side   `json:"white"`
	Black     Side   `json:"black"`
}

type Record struct {
	W int `json:"wins"`
	L int `json:"losses"`
	D int `json:"draws"`
}

func (r Record) Total() int { return r.W + r.L + r.D }

type Row struct {
	Name string `json:"name"`
	Record
}

type Summary struct {
	User      string `json:"user"`
	Months    int    `json:"months"`
	Overall   Record `json:"overall"`
	Colour    []Row  `json:"colour"`
	TimeClass []Row  `json:"timeClass"`
	Openings  []Row  `json:"openings"`
}

var ErrNotFound = errors.New("player not found")

var draws = map[string]bool{"agreed": true, "repetition": true, "stalemate": true, "insufficient": true, "50move": true, "timevsinsufficient": true}

var client = &http.Client{Timeout: 15 * time.Second}

// Fetch downloads the last `months` monthly archives for user and summarises them.
func Fetch(user string, months int) (Summary, error) {
	user = strings.ToLower(user)
	var arch struct{ Archives []string }
	if err := getJSON("https://api.chess.com/pub/player/"+url.PathEscape(user)+"/games/archives", &arch); err != nil {
		return Summary{}, err
	}
	urls := arch.Archives[max(0, len(arch.Archives)-months):]
	var games []Game
	for _, u := range urls {
		var month struct{ Games []Game }
		if err := getJSON(u, &month); err != nil {
			return Summary{}, err
		}
		games = append(games, month.Games...)
	}
	s := Summarise(user, games)
	s.Months = len(urls)
	return s, nil
}

func Summarise(user string, games []Game) Summary {
	var overall Record
	colour, timeClass, opening := map[string]Record{}, map[string]Record{}, map[string]Record{}
	for _, g := range games {
		if g.Rules != "chess" { // skip variants like chess960 or bughouse
			continue
		}
		me, c := g.White, "white"
		if strings.EqualFold(g.Black.Username, user) {
			me, c = g.Black, "black"
		}
		add := func(r Record) Record {
			switch {
			case me.Result == "win":
				r.W++
			case draws[me.Result]:
				r.D++
			default:
				r.L++
			}
			return r
		}
		overall = add(overall)
		colour[c] = add(colour[c])
		timeClass[g.TimeClass] = add(timeClass[g.TimeClass])
		op := OpeningName(g.ECO)
		opening[op] = add(opening[op])
	}
	return Summary{User: user, Overall: overall, Colour: rows(colour), TimeClass: rows(timeClass), Openings: rows(opening)}
}

// rows sorts by games played, most first, then by name.
func rows(m map[string]Record) []Row {
	out := make([]Row, 0, len(m))
	for k, r := range m {
		out = append(out, Row{k, r})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Total() != out[j].Total() {
			return out[i].Total() > out[j].Total()
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// OpeningName turns ".../openings/Italian-Game-Two-Knights-Defense-4.d3" into "Italian Game Two Knights Defense".
func OpeningName(eco string) string {
	slug := eco[strings.LastIndex(eco, "/")+1:]
	var words []string
	for _, w := range strings.Split(slug, "-") {
		if w == "" || (w[0] >= '0' && w[0] <= '9') || w[0] == '.' {
			break
		}
		words = append(words, w)
	}
	if len(words) == 0 {
		return "Unknown"
	}
	return strings.Join(words, " ")
}

func getJSON(u string, v any) error {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "chess-stats-cli") // chess.com blocks requests without one
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: HTTP %d", u, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
