package engine

var Rays = [8][64]Bitboard{}

var KnightMoves = [64]Bitboard{}
var KingMoves = [64]Bitboard{}
var BishopMoves = [64]Bitboard{}
var RookMoves = [64]Bitboard{}
var QueenMoves = [64]Bitboard{}

const (
	North int = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest

	notFileAOrB = ^(FileA | FileB)
	notFileHOrG = ^(FileH | FileG)
)

func InitializeMoveTables() {

	for square := H1; square <= A8; square++ {
		KnightMoves[square] = createKnightMovesForSquare(square)
		KingMoves[square] = createKingMovesForSquare(square)
		BishopMoves[square] = createBishopMovesForSquare(square)
		RookMoves[square] = createRookMovesForSquare(square)
		QueenMoves[square] = BishopMoves[square] | RookMoves[square]
	}
}

func createKingMovesForSquare(square Square) Bitboard {

	startingSquare := SquareMasks[square]

	north := startingSquare << 8
	northEast := startingSquare << 7 & ^FileA
	east := startingSquare >> 1 & ^FileA
	southEast := startingSquare >> 9 & ^FileA
	south := startingSquare >> 8
	southWest := startingSquare >> 7 & ^FileH
	west := startingSquare << 1 & ^FileH
	northWest := startingSquare << 9 & ^FileH

	return north | northEast | east | southEast | south | southWest | west | northWest
}

func createKnightMovesForSquare(square Square) Bitboard {
	startingSquare := SquareMasks[square]

	northNorthWest := startingSquare << 17 & ^FileH
	northNorthEast := startingSquare << 15 & ^FileA

	eastEastNorth := startingSquare << 6 & notFileAOrB
	eastEastSouth := startingSquare >> 10 & notFileAOrB

	westWestNorth := startingSquare << 10 & notFileHOrG
	westWestSouth := startingSquare >> 6 & notFileHOrG

	southSouthEast := startingSquare >> 17 & ^FileA
	southSouthWest := startingSquare >> 15 & ^FileH

	return northNorthWest | northNorthEast | eastEastNorth | westWestNorth | southSouthEast | eastEastSouth | westWestSouth | southSouthWest
}

// createDiagonalMask takes in a square and returns the bitboard that masks the diagonal lines
// (positive and negative) from that square. The formula for calculating this comes from this link:
// https://www.chessprogramming.org/On_an_empty_Board#By_Calculation_3
func createDiagonalMask(square Square) Bitboard {
	diagonal := 56 - 8*(int(square)&7) - (int(square) & 56)
	north := -diagonal & (diagonal >> 31)
	south := diagonal & (-diagonal >> 31)

	return Diagonal >> south << north
}

// createAntiDiagonalMask takes in a square and returns the bitboard that masks the antidiagonal lines
// (positive and megative) from that square. The formula for calculating this comes from the link above.
func createAntiDiagonalMask(square Square) Bitboard {
	antiDiagonal := 8*(int(square)&7) - (int(square) & 56)
	north := -antiDiagonal & (antiDiagonal >> 31)
	south := antiDiagonal & (-antiDiagonal >> 31)

	return AntiDiagonal >> south << north
}

func createBishopMovesForSquare(square Square) Bitboard {
	diagonal := createDiagonalMask(square)
	antiDiagonal := createAntiDiagonalMask(square)

	return (diagonal | antiDiagonal) ^ SquareMasks[square]
}

func createRookMovesForSquare(square Square) Bitboard {
	// the calculations here are taken from the same link as above
	rank := Bitboard(0xff) << (square & 56)
	file := Bitboard(0x0101010101010101) << (square & 7)

	return (rank | file) ^ SquareMasks[square]
}
