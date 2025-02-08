package engine

var KnightMoves = [64]Bitboard{}
var KingMoves = [64]Bitboard{}

const (
	notFileAOrB = ^(FileA | FileB)
	notFileHOrG = ^(FileH | FileG)
)

func InitializeMoveTables() {

	for square := A1; square <= H8; square++ {
		KnightMoves[square] = CreateKnightMovesForSquare(square)
		KingMoves[square] = CreateKingMovesForSquare(square)
	}
}

func CreateKingMovesForSquare(square Square) Bitboard {

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

func CreateKnightMovesForSquare(square Square) Bitboard {
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
