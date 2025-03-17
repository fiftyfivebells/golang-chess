package engine

type Board interface {
	SetBoardFromFEN(fen string)
	ClearBoard()
	GetFENRepresentation() string
	SetPieceAtPosition(p Piece, coord Square)
	GetPieceAtSquare(sq Square) Piece
	RemovePieceFromSquare(sq Square)
	MovePiece(p Piece, from, to Square)
	CastleMove(kingFrom, kingTo Square)
	SquareIsUnderAttack(sq Square, activeSide Color) bool
	SquareIsUnderAttackByPawn(sq Square, activeSide Color) bool
}
