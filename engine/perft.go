package engine

import (
	"fmt"
	"time"
)

const MaxPerftDepth = 10

type PerftState struct {
	gameState   *GameState
	moveBuffers [MaxPerftDepth][256]Move
	moveCounts  [MaxPerftDepth]int
}

func Perft(ps *PerftState, depth int) int64 {
	if depth == 0 {
		return 1
	}

	gameState := ps.gameState
	count := gameState.moveGen.GenerateMoves(&ps.moveBuffers[depth], gameState.ActiveSide, gameState.EPSquare, gameState.CastleRights)

	nodes := int64(0)
	buffer := ps.moveBuffers[depth]
	for i := range count {
		move := buffer[i]
		if gameState.ApplyMove(move) {
			nodes += Perft(ps, depth-1)
		}

		gameState.UnapplyMove(move)
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

func PerftDivide(ps *PerftState, depth int) int64 {

	gameState := ps.gameState
	count := gameState.moveGen.GenerateMoves(&ps.moveBuffers[depth], gameState.ActiveSide, gameState.EPSquare, gameState.CastleRights)
	buffer := ps.moveBuffers[depth]
	var total int64 = 0
	start := time.Now()

	for i := range count {
		move := buffer[i]
		if ps.gameState.ApplyMove(move) {
			count := Perft(ps, depth-1)
			fmt.Printf("%s: %d\n", move.String(), count)
			total += count
		}

		ps.gameState.UnapplyMove(move)
	}

	elapsed := time.Since(start)
	nps := int64(float64(total) / elapsed.Seconds())
	fmt.Printf("Total nodes: %d\n", total)
	fmt.Printf("Time: %s\n", elapsed)
	fmt.Printf("NPS: %d\n", nps)
	return total
}
