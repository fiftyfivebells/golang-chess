package board

import "nsdb-go-edition/engine"

type Board interface {
	SetBoardFromFEN(fen string)
	GetFENRepresentation() string
	SetPieceAtPosition(p engine.Piece, coord string)
	GetPieceAtSquare(sq byte) engine.Piece
}
