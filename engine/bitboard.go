package engine

import (
	"fmt"
	"math/bits"
)

type Bitboard uint64

const (
	FileA Bitboard = 0x8080808080808080 >> iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH

	Rank1 Bitboard = 0xff
	Rank2 Bitboard = Rank1 << 8
	Rank3 Bitboard = Rank2 << 8
	Rank4 Bitboard = Rank3 << 8
	Rank5 Bitboard = Rank4 << 8
	Rank6 Bitboard = Rank5 << 8
	Rank7 Bitboard = Rank6 << 8
	Rank8 Bitboard = Rank7 << 8

	EmptyBitboard = Bitboard(0)
	FullBitboard  = Bitboard(0xffffffffffffffff)
	Diagonal      = Bitboard(0x0102040810204080)
	AntiDiagonal  = Bitboard(0x8040201008040201)
)

var SquareMasks [65]Bitboard

func (bb *Bitboard) setBitAtSquare(square Square) {
	*bb |= SquareMasks[square]
}

func (bb *Bitboard) clearBitAtSquare(square Square) {
	*bb &= ^SquareMasks[square]
}

func (bb Bitboard) lsb() Square {
	bit := bits.TrailingZeros64(uint64(bb))
	return Square(bit)
}

func (bb *Bitboard) PopLSB() Square {
	square := bb.lsb()
	bb.clearBitAtSquare(square)

	return square
}

func (bb Bitboard) String() string {
	bits := fmt.Sprintf("%064b\n", bb)
	bbString := ""
	for rank := H1; rank <= A8 && rank >= H1; rank += 8 {
		bbString += fmt.Sprintf("%v | ", 8-(rank/8))
		for i := rank; i < rank+8; i++ {
			square := bits[i]
			if square == '0' {
				square = '.'
			}
			bbString += fmt.Sprintf("%c ", square)
		}
		bbString += "\n"
	}

	bbString += "   ----------------"
	bbString += "\n    a b c d e f g h\n"

	return bbString
}

func InitializeBitMasks() {
	for i := H1; i <= A8; i++ {
		SquareMasks[i] = Bitboard(1) << i
	}
	SquareMasks[NoSquare] = 0
}
