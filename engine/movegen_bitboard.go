package engine

type BitboardMoveGenerator struct {
	board BitboardBoard
	moves []Move
}

func NewBitboardMoveGenerator(board BitboardBoard) MoveGenerator {
	return &BitboardMoveGenerator{
		board: board,
		moves: []Move{},
	}
}

func (bmg BitboardMoveGenerator) GetMoves() []Move {
	return bmg.moves
}

func (bmg *BitboardMoveGenerator) GenerateMoves(activeSide Color, enPassant Square) {

}

func (bmg *BitboardMoveGenerator) generateMovesByPiece(pieceType PieceType, from Square) {

}

func (bmg *BitboardMoveGenerator) createMovesFromBitboard(from Square, moves, targets Bitboard, pieceType PieceType) {

}
