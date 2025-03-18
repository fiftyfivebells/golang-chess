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

	b.ClearBoard()

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

func (b *BitboardBoard) ClearBoard() {
	for i := range b.squares {
		b.squares[i] = NoPiece
	}

	for color := White; color <= Black; color++ {
		for piece := Pawn; piece <= King; piece++ {
			b.pieces[color][piece] = 0
		}
	}
}

func (b BitboardBoard) GetFENRepresentation() string {
	var fenString strings.Builder

	for rank := A8; rank >= H1 && rank <= A8; rank -= 8 {
		emptySquares := 0
		for i := rank; i > rank-8; i-- {
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

func (b *BitboardBoard) SetPieceAtPosition(p Piece, square Square) {
	b.squares[square] = p

	color := p.Color
	pieceType := p.PieceType

	b.pieces[color][pieceType].setBitAtSquare(square)
}

func (b BitboardBoard) GetPieceAtSquare(sq Square) Piece {
	return b.squares[sq]
}

// RemovePieceFromSquare takes a square and removes the piece from that square, which involves clearing
// the piece from it's associated bitboard and setting the value at the square index in the square array
// to NoPiece
func (b *BitboardBoard) RemovePieceFromSquare(square Square) {
	piece := b.squares[square]

	if piece.PieceType != None {
		color := piece.Color
		pieceType := piece.PieceType

		b.squares[square] = NoPiece
		b.pieces[color][pieceType].clearBitAtSquare(square)
	}
}

// MovePiece takes a piece and two squares (from and to), and moves the given piece to the "to" square
// It first removes the pieces from the "from" and "to" squares, then places the given piece at the
// destination square (the "to" square)
func (b *BitboardBoard) MovePiece(piece Piece, from, to Square) {
	b.RemovePieceFromSquare(from)
	b.RemovePieceFromSquare(to)
	b.SetPieceAtPosition(piece, to)
}

func (b *BitboardBoard) CastleMove(kingFrom, kingTo Square) {
	rookFrom, rookTo := CastlingRookPositions(kingFrom, kingTo)

	color := White
	if kingFrom == E8 {
		color = Black
	}

	king := Piece{
		Color:     color,
		PieceType: King,
	}
	rook := Piece{
		Color:     color,
		PieceType: Rook,
	}

	b.RemovePieceFromSquare(kingFrom)
	b.RemovePieceFromSquare(rookFrom)
	b.SetPieceAtPosition(king, kingTo)
	b.SetPieceAtPosition(rook, rookTo)
}

func (b BitboardBoard) SquareIsUnderAttack(sq Square, activeSide Color) bool {
	enemy := activeSide.EnemyColor()

	bishopMoves := b.GetBishopMoves(sq, activeSide)
	rookMoves := b.GetRookMoves(sq, activeSide)

	pawnAttacks := (PawnAttacks[activeSide][sq] & b.getPiecesByColorAndType(enemy, Pawn)) != 0
	knightAttacks := (KnightMoves[sq] & b.getPiecesByColorAndType(enemy, Knight)) != 0
	bishopAttacks := (bishopMoves & b.getPiecesByColorAndType(enemy, Bishop)) != 0
	rookAttacks := (rookMoves & b.getPiecesByColorAndType(enemy, Rook)) != 0
	queenAttacks := ((bishopMoves | rookMoves) & b.getPiecesByColorAndType(enemy, Queen)) != 0
	kingAttacks := (KingMoves[sq] & b.getPiecesByColorAndType(enemy, King)) != 0

	return pawnAttacks || knightAttacks || bishopAttacks || rookAttacks || queenAttacks || kingAttacks
}

func (b BitboardBoard) SquareIsUnderAttackByPawn(sq Square, activeSide Color) bool {
	enemy := activeSide.EnemyColor()

	pawnAttacks := (PawnAttacks[activeSide][sq] & b.getPiecesByColorAndType(enemy, Pawn)) != 0

	return pawnAttacks
}

func (b BitboardBoard) generateSlidingMoves(square Square, activeSide Color, mask Bitboard) Bitboard {
	squareBoard := SquareMasks[square]
	occupied := b.getAllPieces()

	bottom := ((occupied & mask) - (squareBoard << 1)) & mask
	top := ReverseBitboard(ReverseBitboard((occupied & mask)) - 2*ReverseBitboard(squareBoard))

	allies := b.GetAllPiecesByColor(activeSide)

	return (bottom ^ top) & mask & ^allies
}

func (b BitboardBoard) GetBishopMoves(sq Square, activeSide Color) Bitboard {
	allies := b.GetAllPiecesByColor(activeSide)

	diagonalMask := DiagonalMasks[sq]
	antiDiagonalMask := AntiDiagonalMasks[sq]

	diagonal := b.generateSlidingMoves(sq, activeSide, diagonalMask)
	antiDiagonal := b.generateSlidingMoves(sq, activeSide, antiDiagonalMask)

	return (diagonal | antiDiagonal) & ^allies
}

func (b BitboardBoard) GetRookMoves(sq Square, activeSide Color) Bitboard {
	allies := b.GetAllPiecesByColor(activeSide)

	rank := RankMaskForSquare(sq)
	file := FileMaskForSquare(sq)

	horizontal := b.generateSlidingMoves(sq, activeSide, rank)
	vertical := b.generateSlidingMoves(sq, activeSide, file)

	return (horizontal | vertical) & ^allies
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

func (b BitboardBoard) GetAllPiecesByColor(color Color) Bitboard {
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
