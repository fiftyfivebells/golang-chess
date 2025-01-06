package engine

func CoordToBoardIndex(coord string) byte {
	file := coord[0] - 'a'
	rank := int(coord[1]-'0') - 1

	return byte(rank*8) + file
}

var CharToPiece = map[byte]Piece{
	'P': Piece{Pawn, White},
	'p': Piece{Pawn, Black},
	'N': Piece{Knight, White},
	'n': Piece{Knight, Black},
	'B': Piece{Bishop, White},
	'b': Piece{Bishop, Black},
	'R': Piece{Rook, White},
	'r': Piece{Rook, Black},
	'Q': Piece{Queen, White},
	'q': Piece{Queen, Black},
	'K': Piece{King, White},
	'k': Piece{King, Black},
}
