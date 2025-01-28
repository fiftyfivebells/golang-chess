package engine

import "fmt"

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

	Rank8 Bitboard = 0xff
	Rank7 Bitboard = Rank8 << 8
	Rank6 Bitboard = Rank7 << 8
	Rank5 Bitboard = Rank6 << 8
	Rank4 Bitboard = Rank5 << 8
	Rank3 Bitboard = Rank4 << 8
	Rank2 Bitboard = Rank3 << 8
	Rank1 Bitboard = Rank2 << 8

	EmptyBitboard = Bitboard(0)
	FullBitboard  = Bitboard(0xffffffffffffffff)
)

var SquareMasks [64]Bitboard

func (bb *Bitboard) setBitAtSquare(square Square) {
	*bb |= SquareMasks[square]
}

func (bb *Bitboard) clearBitAtSquare(square Square) {
	*bb &= ^SquareMasks[square]
}

func (bb Bitboard) String() string {
	bits := fmt.Sprintf("%064b\n", bb)
	bbString := ""
	for rank := A8; rank >= 0 && rank <= H8; rank -= 8 {
		bbString += fmt.Sprintf("%v | ", (rank/8)+1)
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
	for i := 0; i < len(SquareMasks); i++ {
		pieceIndex := len(SquareMasks) - 1 - i
		SquareMasks[pieceIndex] = Bitboard(1) << i
	}

}
