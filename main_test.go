package main

import "testing"

func TestOpeningName(t *testing.T) {
	cases := map[string]string{
		"https://www.chess.com/openings/Italian-Game-Two-Knights-Defense-4.d3": "Italian Game Two Knights Defense",
		"https://www.chess.com/openings/Nimzowitsch-Defense-Kennedy-Main-Line": "Nimzowitsch Defense Kennedy Main Line",
		"https://www.chess.com/openings/Sicilian-Defense-...2.Nc6":             "Sicilian Defense",
		"": "Unknown",
	}
	for in, want := range cases {
		if got := openingName(in); got != want {
			t.Errorf("openingName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSummarise(t *testing.T) {
	games := []game{
		{Rules: "chess", TimeClass: "rapid", White: side{"Me", "win"}, Black: side{"x", "resigned"}},
		{Rules: "chess", TimeClass: "rapid", White: side{"x", "win"}, Black: side{"me", "checkmated"}},
		{Rules: "chess", TimeClass: "blitz", White: side{"x", "agreed"}, Black: side{"me", "agreed"}},
		{Rules: "chess960", TimeClass: "blitz", White: side{"me", "win"}, Black: side{"x", "resigned"}},
	}
	s := summarise("me", games)
	if s.overall != (record{1, 1, 1}) {
		t.Errorf("overall = %+v", s.overall)
	}
	if s.colour["white"] != (record{W: 1}) || s.colour["black"] != (record{L: 1, D: 1}) {
		t.Errorf("colour = %+v", s.colour)
	}
	if s.timeClass["rapid"].total() != 2 || s.timeClass["blitz"] != (record{D: 1}) {
		t.Errorf("timeClass = %+v", s.timeClass)
	}
}
