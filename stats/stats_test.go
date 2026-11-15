package stats

import "testing"

func TestOpeningName(t *testing.T) {
	cases := map[string]string{
		"https://www.chess.com/openings/Italian-Game-Two-Knights-Defense-4.d3": "Italian Game Two Knights Defense",
		"https://www.chess.com/openings/Nimzowitsch-Defense-Kennedy-Main-Line": "Nimzowitsch Defense Kennedy Main Line",
		"https://www.chess.com/openings/Sicilian-Defense-...2.Nc6":             "Sicilian Defense",
		"": "Unknown",
	}
	for in, want := range cases {
		if got := OpeningName(in); got != want {
			t.Errorf("OpeningName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSummarise(t *testing.T) {
	games := []Game{
		{Rules: "chess", TimeClass: "rapid", White: Side{"Me", "win"}, Black: Side{"x", "resigned"}},
		{Rules: "chess", TimeClass: "rapid", White: Side{"x", "win"}, Black: Side{"me", "checkmated"}},
		{Rules: "chess", TimeClass: "blitz", White: Side{"x", "agreed"}, Black: Side{"me", "agreed"}},
		{Rules: "chess960", TimeClass: "blitz", White: Side{"me", "win"}, Black: Side{"x", "resigned"}},
	}
	s := Summarise("me", games)
	if s.Overall != (Record{1, 1, 1}) {
		t.Errorf("overall = %+v", s.Overall)
	}
	// black has 2 games so it sorts first
	if len(s.Colour) != 2 || s.Colour[0] != (Row{"black", Record{L: 1, D: 1}}) || s.Colour[1] != (Row{"white", Record{W: 1}}) {
		t.Errorf("colour = %+v", s.Colour)
	}
	if s.TimeClass[0].Name != "rapid" || s.TimeClass[1] != (Row{"blitz", Record{D: 1}}) {
		t.Errorf("timeClass = %+v", s.TimeClass)
	}
}
