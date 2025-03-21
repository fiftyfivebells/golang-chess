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

var CharToPiece = map[byte]Piece{
	'P': {Pawn, White},
	'p': {Pawn, Black},
	'N': {Knight, White},
	'n': {Knight, Black},
	'B': {Bishop, White},
	'b': {Bishop, Black},
	'R': {Rook, White},
	'r': {Rook, Black},
	'Q': {Queen, White},
	'q': {Queen, Black},
	'K': {King, White},
	'k': {King, Black},
}

var PieceToChar = map[Piece]byte{
	{Pawn, White}:   'P',
	{Pawn, Black}:   'p',
	{Knight, White}: 'N',
	{Knight, Black}: 'n',
	{Bishop, White}: 'B',
	{Bishop, Black}: 'b',
	{Rook, White}:   'R',
	{Rook, Black}:   'r',
	{Queen, White}:  'Q',
	{Queen, Black}:  'q',
	{King, White}:   'K',
	{King, Black}:   'k',
	{None, Blank}:   '.',
}
