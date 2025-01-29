package engine

var KnightMoves = [64]Bitboard{}
var KingMoves = [64]Bitboard{}

const (
	notFileA    = FileA ^ FullBitboard
	notFileB    = FileB ^ FullBitboard
	notFileAOrB = notFileA & notFileB

	notFileG    = FileG ^ FullBitboard
	notFileH    = FileH ^ FullBitboard
	notFileHOrG = notFileH & notFileG
)

func InitializeMoveTables() {

	for square := A1; square <= H8; square++ {
		KnightMoves[square] = CreateKnightMovesForSquare(square)
		KingMoves[square] = CreateKingMovesForSquare(square)
	}
}

func CreateKingMovesForSquare(square Square) Bitboard {

	startingSquare := SquareMasks[square]

	north := startingSquare >> 8
	northEast := startingSquare >> 9 & notFileA
	east := startingSquare >> 1 & notFileA
	southEast := startingSquare << 7 & notFileA
	south := startingSquare << 8
	southWest := startingSquare << 9 & notFileH
	west := startingSquare << 1 & notFileH
	northWest := startingSquare >> 7 & notFileH

	return north | northEast | east | southEast | south | southWest | west | northWest
}

func CreateKnightMovesForSquare(square Square) Bitboard {

	startingSquare := SquareMasks[square]

	northNorthWest := startingSquare >> 15 & notFileH
	northNorthEast := startingSquare >> 17 & notFileA

	eastEastNorth := startingSquare >> 10 & notFileAOrB
	eastEastSouth := startingSquare << 6 & notFileAOrB

	westWestNorth := startingSquare >> 6 & notFileHOrG
	westWestSouth := startingSquare << 10 & notFileHOrG

	southSouthEast := startingSquare << 15 & notFileA
	southSouthWest := startingSquare << 17 & notFileH

	return northNorthWest | northNorthEast | eastEastNorth | westWestNorth | southSouthEast | eastEastSouth | westWestSouth | southSouthWest
}
