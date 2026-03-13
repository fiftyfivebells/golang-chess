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

func (pt PieceType) String() string {
	switch pt {
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
	default:
		return ""
	}
}

type Piece byte

const (
	WhitePawn   = 0
	WhiteKnight = 1
	WhiteBishop = 2
	WhiteRook   = 3
	WhiteQueen  = 4
	WhiteKing   = 5
	BlackPawn   = 6
	BlackKnight = 7
	BlackBishop = 8
	BlackRook   = 9
	BlackQueen  = 10
	BlackKing   = 11
	NoPiece     = 12
)

func (p Piece) IsDefined() bool {
	return p < NoPiece
}

func (p Piece) Color() Color {
	if !p.IsDefined() {
		return Blank
	} else {
		return Color(p / 6)
	}
}

func (p Piece) Type() PieceType {
	if !p.IsDefined() {
		return None
	} else {
		return PieceType(p % 6)
	}
}

func (p Piece) String() string {
	switch p {
	case WhitePawn:
		return "P"
	case WhiteKnight:
		return "N"
	case WhiteBishop:
		return "B"
	case WhiteRook:
		return "R"
	case WhiteQueen:
		return "Q"
	case WhiteKing:
		return "K"
	case BlackPawn:
		return "p"
	case BlackKnight:
		return "n"
	case BlackBishop:
		return "b"
	case BlackRook:
		return "r"
	case BlackQueen:
		return "q"
	case BlackKing:
		return "k"

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

type CastleSide byte

const (
	Kingside  CastleSide = 0
	Queenside CastleSide = 1
)

var castleMask = [2][2]CastleAvailability{
	{KingsideWhiteCastle, QueensideWhiteCastle},
	{KingsideBlackCastle, QueensideBlackCastle},
}

func (ca *CastleAvailability) RemoveAllRights(color Color) {
	if color == White {
		*ca &= ^(KingsideWhiteCastle | QueensideWhiteCastle)
	} else {
		*ca &= ^(KingsideBlackCastle | QueensideBlackCastle)
	}
}

func (ca *CastleAvailability) Remove(color Color, side CastleSide) {
	*ca &^= castleMask[color][side]
}

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

	if ca == 0 {
		availability.WriteString("-")
	}

	return availability.String()
}
