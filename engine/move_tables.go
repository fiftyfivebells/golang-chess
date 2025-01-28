package engine

var KnightMoves = [64]Bitboard{}

const (
	notAFile    = FileA ^ FullBitboard
	notBFile    = FileB ^ FullBitboard
	notAOrBFile = notAFile & notBFile

	notGFile    = FileG ^ FullBitboard
	notHFile    = FileH ^ FullBitboard
	notHOrGFile = notHFile & notGFile
)

func InitializeMoveTables() {

	for square := A1; square <= H8; square++ {
		KnightMoves[square] = CreateKnightMovesForSquare(square)
	}
}

func CreateKnightMovesForSquare(square Square) Bitboard {

	startingSquare := SquareMasks[square]

	northNorthWest := startingSquare >> 15 & notHFile
	northNorthEast := startingSquare >> 17 & notAFile

	eastEastNorth := startingSquare >> 10 & notAOrBFile
	eastEastSouth := startingSquare << 6 & notAOrBFile

	westWestNorth := startingSquare >> 6 & notHOrGFile
	westWestSouth := startingSquare << 10 & notHOrGFile

	southSouthEast := startingSquare << 15 & notAFile
	southSouthWest := startingSquare << 17 & notHFile

	return northNorthWest | northNorthEast | eastEastNorth | westWestNorth | southSouthEast | eastEastSouth | westWestSouth | southSouthWest
}
