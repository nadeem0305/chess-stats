// Command chess-stats summarises a chess.com player's recent games:
// results overall, by colour, by time control and by opening.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type side struct {
	Username string `json:"username"`
	Result   string `json:"result"`
}

type game struct {
	TimeClass string `json:"time_class"`
	Rules     string `json:"rules"`
	ECO       string `json:"eco"`
	White     side   `json:"white"`
	Black     side   `json:"black"`
}

type record struct{ W, L, D int }

func (r record) total() int { return r.W + r.L + r.D }

func (r record) String() string {
	return fmt.Sprintf("%4d games  %3dW %3dL %3dD  %3.0f%%", r.total(), r.W, r.L, r.D, 100*float64(r.W)/float64(max(r.total(), 1)))
}

var draws = map[string]bool{"agreed": true, "repetition": true, "stalemate": true, "insufficient": true, "50move": true, "timevsinsufficient": true}

var client = &http.Client{Timeout: 15 * time.Second}

func main() {
	months := flag.Int("months", 3, "how many recent months of games to include")
	top := flag.Int("openings", 10, "how many openings to list")
	flag.Usage = func() { fmt.Fprintln(os.Stderr, "usage: chess-stats [-months N] [-openings N] <username>") }
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	user := strings.ToLower(flag.Arg(0))

	var arch struct{ Archives []string }
	if err := getJSON("https://api.chess.com/pub/player/"+url.PathEscape(user)+"/games/archives", &arch); err != nil {
		fail(err)
	}
	urls := arch.Archives[max(0, len(arch.Archives)-*months):]
	var games []game
	for _, u := range urls {
		var month struct{ Games []game }
		if err := getJSON(u, &month); err != nil {
			fail(err)
		}
		games = append(games, month.Games...)
	}

	s := summarise(user, games)
	fmt.Printf("%s, last %d month(s)\n\nOverall   %v\n", user, len(urls), s.overall)
	printGroup("Colour", s.colour, 0)
	printGroup("Time control", s.timeClass, 0)
	printGroup("Openings", s.opening, *top)
}

type summary struct {
	overall                    record
	colour, timeClass, opening map[string]record
}

func summarise(user string, games []game) summary {
	s := summary{colour: map[string]record{}, timeClass: map[string]record{}, opening: map[string]record{}}
	for _, g := range games {
		if g.Rules != "chess" { // skip variants like chess960 or bughouse
			continue
		}
		me, colour := g.White, "white"
		if strings.EqualFold(g.Black.Username, user) {
			me, colour = g.Black, "black"
		}
		add := func(r record) record {
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
		s.overall = add(s.overall)
		s.colour[colour] = add(s.colour[colour])
		s.timeClass[g.TimeClass] = add(s.timeClass[g.TimeClass])
		op := openingName(g.ECO)
		s.opening[op] = add(s.opening[op])
	}
	return s
}

// openingName turns ".../openings/Italian-Game-Two-Knights-Defense-4.d3" into "Italian Game Two Knights Defense".
func openingName(eco string) string {
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

func printGroup(title string, m map[string]record, limit int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]].total() != m[keys[j]].total() {
			return m[keys[i]].total() > m[keys[j]].total()
		}
		return keys[i] < keys[j]
	})
	if limit > 0 && len(keys) > limit {
		keys = keys[:limit]
	}
	fmt.Printf("\n%s\n", title)
	for _, k := range keys {
		fmt.Printf("  %-44.44s %v\n", k, m[k])
	}
}

func getJSON(u string, v any) error {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "chess-stats-cli") // chess.com blocks requests without one
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("player not found")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: HTTP %d", u, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
