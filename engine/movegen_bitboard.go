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
	activePieces := bmg.board.GetAllPiecesByColor(activeSide)
	targets := bmg.board.GetAllPiecesByColor(^activeSide)

	for pieceType := Knight; pieceType < None; pieceType++ {
		pieceBoard := bmg.board.pieces[activeSide][pieceType]

		for pieceBoard != 0 {
			square := pieceBoard.PopLSB()
			bmg.generateMovesByPiece(pieceType, square, activePieces, targets)
		}
	}

	bmg.generatePawnMoves(activeSide, enPassant)
	bmg.generateCastlingMoves(activeSide)
}

func (bmg *BitboardMoveGenerator) generateMovesByPiece(pieceType PieceType, from Square, activePieces, targets Bitboard) {
	var moves Bitboard
	switch pieceType {
	case Knight:
		moves = (KnightMoves[from] & ^activePieces) & targets
	case King:
		moves = (KingMoves[from] & ^activePieces) & targets
	}

	bmg.createMovesFromBitboard(from, moves, targets, pieceType)
}

func (bmg *BitboardMoveGenerator) generatePawnMoves(activeSide Color, enPassant Square) {}

func (bmg *BitboardMoveGenerator) generateCastlingMoves(activeSide Color) {}

func (bmg *BitboardMoveGenerator) createMovesFromBitboard(from Square, moves, targets Bitboard, pieceType PieceType) {

	for moves != 0 {
		to := moves.PopLSB()
		toBoard := SquareMasks[to]

		moveType := Quiet

		if (toBoard & targets) != 0 {
			moveType = Capture
		}

		move := NewMove(from, to, pieceType, moveType)
		bmg.moves = append(bmg.moves, move)
	}
}
