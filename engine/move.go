package engine

import (
	"fmt"
	"slices"
)

// A move is a 32-bit unsigned integer, where groups of bits represent
// different parts of the move:
// - bits 0 - 5: from square
// - bits 6 - 11: to square
// - bits 12 - 14: piece type
// - bits 15 - 17: promotion piece type
// - bits 18 - 20: move type
const (
	SquareBits = 63
	PieceBits  = 7
	MoveBits   = 15

	ToSquareOffset       = 6
	PieceTypeOffset      = 12
	PromotionPieceOffset = 15
	MoveTypeOffset       = 18
)

type MoveType byte

const (
	NoFlag MoveType = iota
	Quiet
	Capture
	CastleKingside
	CastleQueenside
	EnPassant
	Promotion
	CapturePromotion
)

type Move uint32

func NewMove(from, to Square, pieceType PieceType, moveType MoveType) Move {
	return NewPromotionMove(from, to, pieceType, None, moveType)
}

func NewPromotionMove(from, to Square, pieceType, promotionPieceType PieceType, moveType MoveType) Move {
	newMove := uint32(0)

	newMove |= uint32(from)
	newMove |= (uint32(to) << ToSquareOffset)
	newMove |= (uint32(pieceType) << PieceTypeOffset)
	newMove |= (uint32(moveType) << MoveTypeOffset)
	newMove |= (uint32(promotionPieceType) << PromotionPieceOffset)

	return Move(newMove)
}

func IsAttackMove(moveType MoveType) bool {
	attackMoves := []MoveType{Capture, EnPassant, CapturePromotion}

	return slices.Contains(attackMoves, moveType)
}

func IsPromotionMove(moveType MoveType) bool {
	return moveType == Promotion
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

func (move Move) PromotionPieceType() PieceType {
	return PieceType((move >> PromotionPieceOffset) & PieceBits)
}

func (move Move) String() string {
	from := SquareToCoord(move.FromSquare())
	to := SquareToCoord(move.ToSquare())
	promotionPiece := move.PromotionPieceType()

	return fmt.Sprintf("%s%s%s", from, to, promotionPiece)
}
