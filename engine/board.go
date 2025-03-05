package engine

type Board interface {
	SetBoardFromFEN(fen string)
	GetFENRepresentation() string
	SetPieceAtPosition(p Piece, coord string)
	GetPieceAtSquare(sq Square) Piece
	SquareIsUnderAttack(sq Square, activeSide Color) bool
}
