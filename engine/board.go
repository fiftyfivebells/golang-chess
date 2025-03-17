package engine

type Board interface {
	SetBoardFromFEN(fen string)
	GetFENRepresentation() string
	SetPieceAtPosition(p Piece, coord Square)
	GetPieceAtSquare(sq Square) Piece
	RemovePieceFromSquare(sq Square)
	MovePiece(p Piece, from, to Square)
	SquareIsUnderAttack(sq Square, activeSide Color) bool
}
