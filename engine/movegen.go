package engine

type ScoredMove struct {
	Move  Move
	Score int
}

type MoveGenerator interface {
	GenerateMoves(activeSide Color, enPassant Square)
	GetMoves() []Move
}
