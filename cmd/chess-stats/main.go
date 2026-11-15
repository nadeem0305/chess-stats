// Command chess-stats summarises a chess.com player's recent games:
// results overall, by colour, by time control and by opening.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nadeem0305/chess-stats/stats"
)

func main() {
	months := flag.Int("months", 3, "how many recent months of games to include")
	top := flag.Int("openings", 10, "how many openings to list")
	flag.Usage = func() { fmt.Fprintln(os.Stderr, "usage: chess-stats [-months N] [-openings N] <username>") }
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	s, err := stats.Fetch(flag.Arg(0), *months)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%s, last %d month(s)\n\nOverall   %s\n", s.User, s.Months, format(s.Overall))
	printRows("Colour", s.Colour)
	printRows("Time control", s.TimeClass)
	printRows("Openings", s.Openings[:min(*top, len(s.Openings))])
}

func format(r stats.Record) string {
	return fmt.Sprintf("%4d games  %3dW %3dL %3dD  %3.0f%%", r.Total(), r.W, r.L, r.D, 100*float64(r.W)/float64(max(r.Total(), 1)))
}

func printRows(title string, rows []stats.Row) {
	fmt.Printf("\n%s\n", title)
	for _, r := range rows {
		fmt.Printf("  %-44.44s %s\n", r.Name, format(r.Record))
	}
}
