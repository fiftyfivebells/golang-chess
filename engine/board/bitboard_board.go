package board

import (
	"fmt"
	"nsdb-go-edition/engine"
)

type BitboardBoard struct {
	pieces  [2][6]Bitboard
	squares [64]engine.Piece
}

func (b *BitboardBoard) SetBoardFromFEN(fen string) {
	InitializeBitMasks()

	for i := range b.squares {
		b.squares[i] = engine.Piece{
			PieceType: engine.None,
			Color:     engine.Blank,
		}
	}

	for i, square := 0, Square(56); i < len(fen); i++ {
		char := fen[i]
		switch char {
		case 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k':
			piece := engine.CharToPiece[char]
			b.pieces[piece.Color][piece.PieceType].setBitAtSquare(square)

			b.squares[square] = piece
			square++
		case '/':
			square -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			square += fen[i] - '0'
		}
	}
}

func (b *BitboardBoard) GetFENRepresentation() string {
	fenString := ""

	for rank := 56; rank > -1; rank -= 8 {
		fenString += fmt.Sprintf("%v | ", (rank/8)+1)
		for i := rank; i < rank+8; i++ {
			char := engine.PieceToChar[b.squares[i]]
			fenString += fmt.Sprintf("%c ", char)
		}
		fenString += "\n"
	}
	fenString += "    a b c d e f g h\n"

	return fenString
}

func (b *BitboardBoard) SetPieceAtPosition(p engine.Piece, coord string) {
	boardIndex := CoordToBoardIndex(coord)
	b.squares[boardIndex] = p

	color := p.Color
	pieceType := p.PieceType

	b.pieces[color][pieceType].setBitAtSquare(boardIndex)
}

func (b *BitboardBoard) GetPieceAtSquare(sq Square) engine.Piece {
	return b.squares[sq]
}
