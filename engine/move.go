package engine

import "fmt"

// A move is a 32-bit unsigned integer, where groups of bits represent
// different parts of the move:
// - bits 0 - 5: from square
// - bits 6 - 11: to square
// - bits 12 - 14: piece type
// - bits 15 - 17:
const (
	SquareBits = 63
	PieceBits  = 7
	MoveBits   = 15

	ToSquareOffset  = 6
	PieceTypeOffset = 12
	MoveTypeOffset  = 15
)

type MoveType byte

const (
	NoFlag MoveType = iota
	Quiet
	Capture
	CastleKingside
	CastleQueenside
	EnPassant
	PromotionKnight
	PromotionBishop
	PromotionRook
	PromotionQueen
	CapturePromotionKnight
	CapturePromotionBishop
	CapturePromotionRook
	CapturePromotionQueen
)

type Move uint32

func NewMove(from, to Square, pieceType PieceType, moveType MoveType) Move {
	newMove := uint32(0)

	newMove |= uint32(from)
	newMove |= (uint32(to) << ToSquareOffset)
	newMove |= (uint32(pieceType) << PieceTypeOffset)
	newMove |= (uint32(moveType) << MoveTypeOffset)

	return Move(newMove)
}

func (move Move) FromSquare() Square {
	return Square(move & SquareBits)
}

func (move Move) ToSquare() Square {
	return Square((move >> ToSquareOffset) & SquareBits)
}

func (move Move) PieceType() PieceType {
	return PieceType((move >> PieceTypeOffset) & PieceBits)
}

func (move Move) MoveType() MoveType {
	return MoveType((move >> MoveTypeOffset) & MoveBits)
}

func (move Move) String() string {
	from := SquareToCoord(move.FromSquare())
	to := SquareToCoord(move.ToSquare())

	return fmt.Sprintf("%s%s", from, to)
}
