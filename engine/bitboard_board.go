package engine

import (
	"fmt"
	"strconv"
	"strings"
)

type BitboardBoard struct {
	pieces  [2][6]Bitboard
	squares [64]Piece
}

func NewBitboardBoard(fen string) *BitboardBoard {
	board := BitboardBoard{}
	board.SetBoardFromFEN(fen)

	return &board
}

func (b *BitboardBoard) SetBoardFromFEN(fen string) {

	for i := range b.squares {
		b.squares[i] = Piece{
			PieceType: None,
			Color:     Blank,
		}
	}

	for i, square := 0, A8; i < len(fen); i++ {
		char := fen[i]
		switch byte(char) {
		case 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k':
			piece := CharToPiece[char]
			b.pieces[piece.Color][piece.PieceType].setBitAtSquare(square)
			b.squares[square] = piece
			square--
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			square -= Square(char - '0')
		}
	}
}

func (b BitboardBoard) GetFENRepresentation() string {
	var fenString strings.Builder

	for rank := H8; rank >= H1 && rank <= A8; rank -= 8 {
		emptySquares := 0
		for i := rank; i < rank+8; i++ {
			piece := b.squares[i]

			if piece.PieceType == None {
				emptySquares++
			} else {
				if emptySquares > 0 {
					fenString.WriteString(strconv.Itoa(emptySquares))
					emptySquares = 0
				}

				fenString.WriteRune(rune(PieceToChar[piece]))
			}
		}

		if emptySquares > 0 {
			fenString.WriteString(strconv.Itoa(emptySquares))
		}

		fenString.WriteString("/")
	}

	return strings.TrimSuffix(fenString.String(), "/")
}

func (b *BitboardBoard) SetPieceAtPosition(p Piece, coord string) {
	boardIndex := CoordToBoardIndex(coord)
	b.squares[boardIndex] = p

	color := p.Color
	pieceType := p.PieceType

	b.pieces[color][pieceType].setBitAtSquare(boardIndex)
}

func (b BitboardBoard) GetPieceAtSquare(sq Square) Piece {
	return b.squares[sq]
}

func (b BitboardBoard) SquareIsUnderAttack(sq Square, activeSide Color) bool {
	enemy := activeSide.EnemyColor()

	pawnAttacks := (PawnPushes[activeSide][sq] & b.getPiecesByColorAndType(enemy, Pawn)) != 0
	knightAttacks := (KnightMoves[sq] & b.getPiecesByColorAndType(enemy, Knight)) != 0
	kingAttacks := (KingMoves[sq] & b.getPiecesByColorAndType(enemy, King)) != 0

	return pawnAttacks || knightAttacks || kingAttacks
}

func (b BitboardBoard) getAllPieces() Bitboard {
	bb := Bitboard(0)

	for color := White; color < Blank; color++ {
		for pieceType := Pawn; pieceType < None; pieceType++ {
			bb |= b.pieces[color][pieceType]
		}
	}

	return bb
}

func (b BitboardBoard) getAllPiecesByColor(color Color) Bitboard {
	bb := Bitboard(0)

	for pieceType := Pawn; pieceType < None; pieceType++ {
		bb |= b.pieces[color][pieceType]
	}

	return bb
}

func (b BitboardBoard) getPiecesByColorAndType(color Color, pieceType PieceType) Bitboard {
	return b.pieces[color][pieceType]
}

func (b BitboardBoard) String() string {
	var stringRep string

	for rank := 8; rank > 0; rank-- {
		square := rank*8 - 1
		stringRep += fmt.Sprintf("%v | ", square/8+1)
		for i := square; i > square-8; i-- {
			piece := PieceToChar[b.squares[i]]
			stringRep += fmt.Sprintf("%s ", string(piece))
		}
		stringRep += "\n"
	}

	stringRep += "   ----------------"
	stringRep += "\n    a b c d e f g h\n"

	return stringRep
}
