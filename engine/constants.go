package engine

type Square byte

const (
	InitialStateFenString = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	A1, A2, A3, A4, A5, A6, A7, A8 Square = 0, 8, 16, 24, 32, 40, 48, 56
	B1, B2, B3, B4, B5, B6, B7, B8 Square = 1, 9, 17, 25, 33, 41, 49, 57
	C1, C2, C3, C4, C5, C6, C7, C8 Square = 2, 10, 18, 26, 34, 42, 50, 58
	D1, D2, D3, D4, D5, D6, D7, D8 Square = 3, 11, 19, 27, 35, 43, 51, 59
	E1, E2, E3, E4, E5, E6, E7, E8 Square = 4, 12, 20, 28, 36, 44, 52, 60
	F1, F2, F3, F4, F5, F6, F7, F8 Square = 5, 13, 21, 29, 37, 45, 53, 61
	G1, G2, G3, G4, G5, G6, G7, G8 Square = 6, 14, 22, 30, 38, 46, 54, 62
	H1, H2, H3, H4, H5, H6, H7, H8 Square = 7, 15, 23, 31, 39, 47, 55, 63
	NoSquare                       Square = 64
)
