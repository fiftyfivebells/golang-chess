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

	F1G1Mask   = Bitboard(6)
	B1C1D1Mask = Bitboard(0x70)
	F8G8Mask   = Bitboard(0x600000000000000)
	B8C8D8Mask = Bitboard(0x7000000000000000)
)

var SquareMasks [65]Bitboard

func ReverseBitboard(bb Bitboard) Bitboard {
	asInt := uint64(bb)

	return Bitboard(bits.Reverse64(asInt))
}

func RankMaskForSquare(square Square) Bitboard {
	rank := square / 8

	switch rank {
	case 0:
		return Rank1
	case 1:
		return Rank2
	case 2:
		return Rank3
	case 3:
		return Rank4
	case 4:
		return Rank5
	case 5:
		return Rank6
	case 6:
		return Rank7
	case 7:
		return Rank8
	default:
		return 0
	}
}

func FileMaskForSquare(square Square) Bitboard {
	file := int(square % 8)

	switch file {
	case 7:
		return FileA
	case 6:
		return FileB
	case 5:
		return FileC
	case 4:
		return FileD
	case 3:
		return FileE
	case 2:
		return FileF
	case 1:
		return FileG
	case 0:
		return FileH
	default:
		return 0
	}
}

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
