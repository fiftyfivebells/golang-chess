package board

type Board interface {
	SetBoardFromFEN(fen string)
	GetFENRepresentation() string
	SetPieceAtPosition(p Piece, coord string)
	GetPieceAtSquare(sq byte) Piece
}
