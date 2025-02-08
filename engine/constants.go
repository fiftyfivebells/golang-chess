package engine

type Square byte

const (
	InitialStateFenString = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	A1, B1, C1, D1, E1, F1, G1, H1 Square = 7, 6, 5, 4, 3, 2, 1, 0
	A2, B2, C2, D2, E2, F2, G2, H2 Square = 15, 14, 13, 12, 11, 10, 9, 8
	A3, B3, C3, D3, E3, F3, G3, H3 Square = 23, 22, 21, 20, 19, 18, 17, 16
	A4, B4, C4, D4, E4, F4, G4, H4 Square = 31, 30, 29, 28, 27, 26, 25, 24
	A5, B5, C5, D5, E5, F5, G5, H5 Square = 39, 38, 37, 36, 35, 34, 33, 32
	A6, B6, C6, D6, E6, F6, G6, H6 Square = 47, 46, 45, 44, 43, 42, 41, 40
	A7, B7, C7, D7, E7, F7, G7, H7 Square = 55, 54, 53, 52, 51, 50, 49, 48
	A8, B8, C8, D8, E8, F8, G8, H8 Square = 63, 62, 61, 60, 59, 58, 57, 56
	NoSquare                       Square = 64
)
