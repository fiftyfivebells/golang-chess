package engine

type Color byte

const (
	White Color = 0
	Black Color = 1
	Blank Color = 2
)

type PieceType byte

const (
	Pawn   PieceType = 0
	Knight PieceType = 1
	Bishop PieceType = 2
	Rook   PieceType = 3
	Queen  PieceType = 4
	King   PieceType = 5
	None   PieceType = 6
)

type Piece struct {
	PieceType PieceType
	Color     Color
}
