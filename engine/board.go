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
	ReverseCastleMove(kingFrom, kingTo Square)
	KingIsUnderAttack(color Color) bool
	SquareIsUnderAttack(sq Square, activeSide Color) bool
	SquareIsUnderAttackByPawn(sq Square, activeSide Color) bool
}
