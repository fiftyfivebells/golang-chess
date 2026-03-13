package engine

import "testing"

func BenchmarkPerft(b *testing.B) {
	state := InitializeGameState(InitialStateFenString)
	ps := &PerftState{gameState: &state}
	b.ResetTimer()

	for b.Loop() {
		Perft(ps, 5)
	}
}
