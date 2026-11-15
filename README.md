# chess-stats

**Live:** https://chess-stats-chi.vercel.app

See how you actually do on chess.com: win rate overall, with white vs black, per time control,
and per opening. Uses the public [chess.com API](https://www.chess.com/news/view/published-data-api), no key needed.

```
go install github.com/nadeem0305/chess-stats/cmd/chess-stats@latest
chess-stats nadeem0305
chess-stats -months 6 -openings 15 nadeem0305
```

Example (`-openings 5`):

```
nadeem0305, last 2 month(s)

Overall    202 games   90W  99L  13D   45%

Colour
  white                                         103 games   50W  49L   4D   49%
  black                                          99 games   40W  50L   9D   40%

Time control
  rapid                                         116 games   48W  60L   8D   41%
  blitz                                          75 games   34W  37L   4D   45%
  daily                                          10 games    8W   1L   1D   80%
  bullet                                          1 games    0W   1L   0D    0%

Openings
  Alekhines Defense                              15 games    8W   6L   1D   53%
  Kings Pawn Opening Kings Knight Variation      15 games    7W   8L   0D   47%
  Scandinavian Defense                           15 games    6W   9L   0D   40%
  Queens Pawn Opening                            14 games    3W   9L   2D   21%
  Caro Kann Defense                              10 games    4W   4L   2D   40%
```

Variants (Chess960, bughouse, ...) are skipped. Openings come from chess.com's own classification,
trimmed to the named line (move numbers dropped).
