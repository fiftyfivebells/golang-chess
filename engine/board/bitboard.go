package board

import "fmt"

type Bitboard uint64

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
	for rank := 56; rank > -1; rank -= 8 {
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
