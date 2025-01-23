package engine

import "strings"

type Color byte

const (
	White Color = 0
	Black Color = 1
	Blank Color = 2
)

func (c Color) String() string {
	switch c {
	case White:
		return "w"
	case Black:
		return "b"
	default:
		return ""
	}
}

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

type CastleAvailability byte

const (
	KingsideWhiteCastle  CastleAvailability = 0b0001
	QueensideWhiteCastle CastleAvailability = 0b0010
	KingsideBlackCastle  CastleAvailability = 0b0100
	QueensideBlackCastle CastleAvailability = 0b1000
)

func (ca CastleAvailability) String() string {
	var availability strings.Builder

	if (ca & KingsideWhiteCastle) != 0 {
		availability.WriteString("K")
	}
	if (ca & QueensideWhiteCastle) != 0 {
		availability.WriteString("Q")
	}
	if (ca & KingsideBlackCastle) != 0 {
		availability.WriteString("k")
	}
	if (ca & QueensideBlackCastle) != 0 {
		availability.WriteString("q")
	}

	return availability.String()
}
