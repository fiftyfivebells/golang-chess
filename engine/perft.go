package engine

import (
	"fmt"
	"time"
)

func Perft(state GameState, depth int) int64 {
	if depth == 0 {
		return 1
	}

	pseudoMoves := state.GetPseudoLegalMovesForPosition()
	n := copy(state.LegalMovesBuffer[:], pseudoMoves)
	moves := state.LegalMovesBuffer[:n]
	nodes := int64(0)

	for _, move := range moves {
		if state.ApplyMove(move) {
			nodes += Perft(state, depth-1)
		}

		state.UnapplyMove(move)
	}

	return nodes
}

func PerftTrace(state GameState, depth int, trace []Move) int64 {
	if depth == 0 {
		fmt.Printf("TRACE: %+v\n", trace)
		return 1
	}

	pseudoMoves := state.GetPseudoLegalMovesForPosition()
	n := copy(state.LegalMovesBuffer[:], pseudoMoves)
	moves := state.LegalMovesBuffer[:n]
	nodes := int64(0)

	for _, move := range moves {
		if state.ApplyMove(move) {
			trace = append(trace, move)
			nodes += PerftTrace(state, depth-1, trace)
		}

		state.UnapplyMove(move)
	}

	return nodes
}

func PerftDivide(state GameState, depth int) int64 {
	pseudoMoves := state.GetPseudoLegalMovesForPosition()
	n := copy(state.LegalMovesBuffer[:], pseudoMoves)
	moves := state.LegalMovesBuffer[:n]

	var total int64 = 0
	start := time.Now()

	for _, move := range moves {
		if state.ApplyMove(move) {
			count := Perft(state, depth-1)
			fmt.Printf("%s: %d\n", move.String(), count)
			total += count
		}

		state.UnapplyMove(move)
	}

	elapsed := time.Since(start)
	nps := int64(float64(total) / elapsed.Seconds())
	fmt.Printf("Total nodes: %d\n", total)
	fmt.Printf("Time: %s\n", elapsed)
	fmt.Printf("NPS: %d\n", nps)
	return total
}
