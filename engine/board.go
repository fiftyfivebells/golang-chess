package engine

import (
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	pieces    [2][6]Bitboard
	squares   [64]Piece
	colorBB   [2]Bitboard
	occupancy Bitboard
	kingSq    [2]Square
}

func NewBoard(fen string) *Board {
	board := Board{}
	board.SetBoardFromFEN(fen)

	return &board
}

func (b *Board) SetBoardFromFEN(fen string) {

	b.ClearBoard()

	for i, square := 0, A8; i < len(fen); i++ {
		char := fen[i]
		switch byte(char) {
		case 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k':
			piece := CharToPiece[char]
			b.SetPieceAtPosition(piece, square)
			square--
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			square -= Square(char - '0')
		}
	}

	b.kingSq[White] = b.pieces[White][King].lsb()
	b.kingSq[Black] = b.pieces[Black][King].lsb()
}

func (b *Board) ClearBoard() {
	for i := range b.squares {
		b.squares[i] = NoPiece
	}

	for color := White; color <= Black; color++ {
		for piece := Pawn; piece <= King; piece++ {
			b.pieces[color][piece] = 0
		}
		b.colorBB[color] = 0
	}
	b.occupancy = 0
}

func (b Board) GetFENRepresentation() string {
	var fenString strings.Builder

	for rank := 8; rank > 0; rank-- {
		emptySquares := 0
		startingSquare := rank*8 - 1

		for i := startingSquare; i > startingSquare-8; i-- {
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

func (b *Board) SetPieceAtPosition(p Piece, square Square) {
	b.squares[square] = p

	if p == NoPiece {
		return
	}

	color := p.Color
	pieceType := p.PieceType

	b.pieces[color][pieceType].setBitAtSquare(square)
	b.colorBB[color].setBitAtSquare(square)
	b.occupancy.setBitAtSquare(square)

	if pieceType == King {
		b.kingSq[color] = square
	}
}

func (b Board) GetPieceAtSquare(sq Square) Piece {
	return b.squares[sq]
}

// RemovePieceFromSquare takes a square and removes the piece from that square, which involves clearing
// the piece from it's associated bitboard and setting the value at the square index in the square array
// to NoPiece
func (b *Board) RemovePieceFromSquare(square Square) {
	piece := b.squares[square]

	if piece.PieceType != None {
		b.pieces[piece.Color][piece.PieceType].clearBitAtSquare(square)
		b.colorBB[piece.Color].clearBitAtSquare(square)
		b.occupancy.clearBitAtSquare(square)
		b.squares[square] = NoPiece
	}
}

// MovePiece takes a piece and two squares (from and to), and moves the given piece to the "to" square
// It first removes the pieces from the "from" and "to" squares, then places the given piece at the
// destination square (the "to" square)
func (b *Board) MovePiece(piece Piece, from, to Square) {
	b.RemovePieceFromSquare(from)
	b.RemovePieceFromSquare(to)
	b.SetPieceAtPosition(piece, to)
}

func (b *Board) CastleMove(kingFrom, kingTo Square) {
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

func (b *Board) ReverseCastleMove(kingFrom, kingTo Square) {
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

	b.RemovePieceFromSquare(kingTo)
	b.RemovePieceFromSquare(rookTo)
	b.SetPieceAtPosition(king, kingFrom)
	b.SetPieceAtPosition(rook, rookFrom)
}

func (b Board) SquareIsUnderAttack(sq Square, activeSide Color) bool {
	enemy := activeSide.EnemyColor()
	occ := b.occupancy

	if PawnAttacks[activeSide][sq]&b.pieces[enemy][Pawn] != 0 {
		return true
	}
	if KnightMoves[sq]&b.pieces[enemy][Knight] != 0 {
		return true
	}
	if KingMoves[sq]&b.pieces[enemy][King] != 0 {
		return true
	}

	be := &BishopMagics[sq]
	diagAttacks := BishopAttacks[sq][((occ&be.mask)*Bitboard(be.magic))>>be.shift]
	if diagAttacks&(b.pieces[enemy][Bishop]|b.pieces[enemy][Queen]) != 0 {
		return true
	}

	re := &RookMagics[sq]
	orthoAttacks := RookAttacks[sq][((occ&re.mask)*Bitboard(re.magic))>>re.shift]
	return orthoAttacks&(b.pieces[enemy][Rook]|b.pieces[enemy][Queen]) != 0
}

func (b Board) KingIsUnderAttack(color Color) bool {
	return b.SquareIsUnderAttack(b.kingSq[color], color)
}

func (b Board) SquareIsUnderAttackByPawn(sq Square, activeSide Color) bool {
	enemy := activeSide.EnemyColor()

	pawnAttacks := (PawnAttacks[activeSide][sq] & b.getPiecesByColorAndType(enemy, Pawn)) != 0

	return pawnAttacks
}

func (b Board) generateSlidingMoves(square Square, occupied, allies, mask Bitboard) Bitboard {
	squareBoard := SquareMasks[square]

	bottom := ((occupied & mask) - (squareBoard << 1)) & mask
	top := ReverseBitboard(ReverseBitboard((occupied & mask)) - 2*ReverseBitboard(squareBoard))

	return (bottom ^ top) & mask & ^allies
}

func (b Board) GetBishopMoves(sq Square, occupied, allies Bitboard) Bitboard {
	e := &BishopMagics[sq]
	idx := (occupied & e.mask) * Bitboard(e.magic) >> e.shift
	return BishopAttacks[sq][idx] & ^allies
}

func (b Board) GetRookMoves(sq Square, occupied, allies Bitboard) Bitboard {
	e := &RookMagics[sq]
	idx := (occupied & e.mask) * Bitboard(e.magic) >> e.shift
	return RookAttacks[sq][idx] & ^allies
}

func (b Board) getAllPieces() Bitboard {
	return b.occupancy
}

func (b Board) GetAllPiecesByColor(color Color) Bitboard {
	return b.colorBB[color]
}

func (b Board) getPiecesByColorAndType(color Color, pieceType PieceType) Bitboard {
	return b.pieces[color][pieceType]
}

func (b Board) String() string {
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
