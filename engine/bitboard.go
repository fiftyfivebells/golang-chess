package engine

type Bitboard uint64

const (
	InitialStateFenString = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

type BitboardBoard struct {
	pieces  [2][6]Bitboard
	squares [64]Piece
}

func (b *BitboardBoard) SetBoardFromFEN(fen string) {
	for i := range b.squares {
		b.squares[i] = Piece{None, Blank}
	}

	for i, square := 0, byte(56); i < len(fen); i++ {
		char := fen[i]
		switch char {
		case 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k':
			piece := CharToPiece[char]

			pieceMask := 1 << square
			b.pieces[piece.Color][piece.PieceType] |= Bitboard(pieceMask)

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
	return "fen rep"
}

func (b *BitboardBoard) SetPieceAtPosition(p Piece, coord string) {
	boardIndex := CoordToBoardIndex(coord)
	b.squares[boardIndex] = p

	pieceMask := 1 << boardIndex

	color := p.Color
	pieceType := p.PieceType

	b.pieces[color][pieceType] |= Bitboard(pieceMask)
}

func (b *BitboardBoard) GetPieceAtSquare(sq byte) Piece {
	return b.squares[sq]
}
