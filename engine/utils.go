package engine

func CoordToBoardIndex(coord string) Square {
	if len(coord) != 2 {
		return NoSquare
	}

	file := 7 - (coord[0] - 'a')
	rank := byte(coord[1]-'0') - 1

	return Square(rank*8 + file)
}

func SquareToCoord(square Square) string {
	file := square % 8
	rank := square / 8

	return string(rune('h'-file)) + string(rune('0'+rank+1))
}

func SquareToFileRank(square Square) (int, int) {
	file := int(square % 8)
	rank := int(square / 8)

	return file, rank
}

func CharToColor(ch string) Color {
	var color Color
	switch ch {
	case "w":
		color = White
	case "b":
		color = Black
	default:
		color = Blank
	}

	return color
}

func CastlingRookPositions(kingFrom, kingTo Square) (rookFrom, rookTo Square) {
	switch kingTo {
	case G1:
		rookFrom, rookTo = H1, F1
	case C1:
		rookFrom, rookTo = A1, D1
	case G8:
		rookFrom, rookTo = H8, F8
	case C8:
		rookFrom, rookTo = A8, D8
	}

	return rookFrom, rookTo
}

func makePiece(pt PieceType, c Color) Piece {
	p := uint8(pt) + uint8(c*6)

	return Piece(p)
}

// TODO: I can probably make this an array instead of a map
var CharToPiece = map[byte]Piece{
	'P': WhitePawn,
	'p': BlackPawn,
	'N': WhiteKnight,
	'n': BlackKnight,
	'B': WhiteBishop,
	'b': BlackBishop,
	'R': WhiteRook,
	'r': BlackRook,
	'Q': WhiteQueen,
	'q': BlackQueen,
	'K': WhiteKing,
	'k': BlackKing,
}

var PieceToChar = [Piece]{
		'P',
		'N',
		'B',
		'R',
		'Q',
		'K',
		'p',
		'n',
		'b',
		'r',
		'q',
		'k',
		'.'
}
