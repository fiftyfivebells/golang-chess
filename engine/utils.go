package engine

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
