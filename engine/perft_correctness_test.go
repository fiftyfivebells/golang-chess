package engine

import (
	"testing"
)

func TestPerftCorrectness(t *testing.T) {
	tests := []struct {
		name  string
		fen   string
		depth int
		nodes int64
	}{
		{"startpos depth 1", InitialStateFenString, 1, 20},
		{"startpos depth 2", InitialStateFenString, 2, 400},
		{"startpos depth 3", InitialStateFenString, 3, 8902},
		{"startpos depth 4", InitialStateFenString, 4, 197281},
		{"startpos depth 5", InitialStateFenString, 5, 4865609},
		{"kiwipete depth 1", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1, 48},
		{"kiwipete depth 2", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2, 2039},
		{"kiwipete depth 3", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3, 97862},
		{"kiwipete depth 4", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4, 4085603},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := InitializeGameState(tt.fen)
			ps := &PerftState{gameState: &state}
			got := Perft(ps, tt.depth)
			if got != tt.nodes {
				t.Errorf("depth %d: got %d, want %d", tt.depth, got, tt.nodes)
			}
		})
	}
}
