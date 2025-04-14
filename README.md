# Not So Deep Blue – Go Edition

![Go Version](https://img.shields.io/badge/Go-1.21+-blue)
![Status](https://img.shields.io/badge/status-work%20in%20progress-orange)

> A bitboard-based UCI chess engine written in Go.
> *It plays fast. It plays legal. It just doesn't play well.*

**Not So Deep Blue – Go Edition** is a modest attempt at a chess engine.
It’s written in Go, uses bitboards for board representation, and communicates using the [UCI protocol](https://gist.github.com/aliostad/02babb68c964ebf8d57c).
It's not the best, but it's... yours.

---

## Features

- Full board representation using 64-bit bitboards
- UCI protocol support (compatible with most chess GUIs)
- Legal move generation and validation
- Lightweight, fast startup, easy to debug
- Plays a complete game of chess from start to checkmate (or stalemate... or confusion)
- **Currently only plays random moves!!** (search is coming next)

---


## Getting Started

Work in progress...installation instructions coming soon.
<!-- <\!-- ### 📦 Installation -\-> -->

<!-- <\!-- ```bash -\-> -->
<!-- <\!-- git clone https://github.com/fiftyfivebells/golang-chess.git -\-> -->
<!-- <\!-- cd not-so-deep-blue-go -\-> -->
<!-- <\!-- go build -o nsdb -\-> -->
<!-- <\!-- ``` -\-> -->

<!-- Now you can run it directly or load it into a GUI like [Arena](http://www.playwitharena.de/), [CuteChess](https://cutechess.com/), or [Lucas Chess](https://lucaschess.pythonanywhere.com/). -->

<!-- --- -->

<!-- ### 🕹️ Using with a GUI -->

<!-- 1. Launch your GUI of choice -->
<!-- 2. Add a new UCI engine and point it to the `nsdb` binary -->
<!-- 3. Play -->

---

## Example UCI Session (CLI)

```bash
$ ./nsdb
uci
id name Not So Deep Blue - Go Edition
id author Stephen Bell
uciok
isready
readyok
position startpos moves e2e4 e7e5
go
bestmove g1f3
```

> Yes, it does respond. No, I can't promise it saw that fork.

---

## Philosophy

This project is:

- A playground for bitboard logic in Go
- A chance to implement UCI
- A reminder that even legal moves can be bad ones
- Not remotely a threat to Stockfish...but that's all right!

---

## Known Limitations

- Alpha-beta search (next on the list for implementation!)  
- No transposition table (yet)  
- No quiescence search (also yet)  
- Evaluation is... optimistic  
- Time management is very trusting  
- Loses most games, but politely  

---

## License

MIT — because making bad moves shouldn’t require a license.  
See [`LICENSE`](./LICENSE) for details.

---

## Credits

- Inspired by Deep Blue (but not *that* deep)  
- Bitboards drawn with love and bitwise ops  
- Thanks to everyone who writes better chess engines and shares what they learn  
