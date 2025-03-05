package engine

import "strings"

type Color byte

const (
	White Color = 0
	Black Color = 1
	Blank Color = 2
)

func (c Color) EnemyColor() Color {
	if c != Blank {
		return c ^ 1
	}

	return Blank
}

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

func (p Piece) String() string {
	if p.Color == White {
		switch p.PieceType {
		case Pawn:
			return "P"
		case Knight:
			return "N"
		case Bishop:
			return "B"
		case Rook:
			return "R"
		case Queen:
			return "Q"
		case King:
			return "K"
		}
	} else if p.Color == Black {
		switch p.PieceType {
		case Pawn:
			return "p"
		case Knight:
			return "n"
		case Bishop:
			return "b"
		case Rook:
			return "r"
		case Queen:
			return "q"
		case King:
			return "k"
		}
	}

	return "none"
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
