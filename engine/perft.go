package engine

import "fmt"

func Perft(state GameState, depth int) int64 {
	if depth == 0 {
		return 1
	}

	moves := state.GetMovesForPosition()
	nodes := int64(0)

	for i := 0; i < len(moves); i++ {

		move := moves[i]

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

	moves := state.GetMovesForPosition()
	nodes := int64(0)

	for i := 0; i < len(moves); i++ {

		move := moves[i]

		if state.ApplyMove(move) {
			trace = append(trace, move)
			nodes += PerftTrace(state, depth-1, trace)
		}

		state.UnapplyMove(move)
	}

	return nodes
}

func PerftDivide(state GameState, depth int) int64 {
	moves := state.GetMovesForPosition()

	var total int64 = 0

	for i := 0; i < len(moves); i++ {
		move := moves[i]
		if state.ApplyMove(move) {

			count := Perft(state, depth-1)
			fmt.Printf("%s: %d\n", move.String(), count)
			total += count
		}

		state.UnapplyMove(move)
	}

	fmt.Printf("Total nodes: %d\n", total)
	return total
}
