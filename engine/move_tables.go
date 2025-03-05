package engine

var PawnPushes = [2][64]Bitboard{}
var PawnAttacks = [2][64]Bitboard{}

var KnightMoves = [64]Bitboard{}
var KingMoves = [64]Bitboard{}
var BishopMoves = [64]Bitboard{}
var RookMoves = [64]Bitboard{}
var QueenMoves = [64]Bitboard{}

var DiagonalMasks = [64]Bitboard{}
var AntiDiagonalMasks = [64]Bitboard{}
var HorizontalMasks = [64]Bitboard{}
var VerticalMasks = [64]Bitboard{}

const (
	North     = 8
	NorthEast = 7
	East      = 1
	SouthEast = 9
	South     = 8
	SouthWest = 7
	West      = 1
	NorthWest = 9

	notFileAOrB = ^(FileA | FileB)
	notFileHOrG = ^(FileH | FileG)
)

func InitializeMoveTables() {

	for square := H1; square <= A8; square++ {
		board := SquareMasks[square]
		PawnPushes[White][square] = board << North
		PawnPushes[Black][square] = board >> South
		PawnAttacks[White][square] = createWhitePawnAttacksForSquare(square)
		PawnAttacks[Black][square] = createBlackPawnAttacksForSquare(square)

		KnightMoves[square] = createKnightMovesForSquare(square)
		KingMoves[square] = createKingMovesForSquare(square)

		HorizontalMasks[square] = createHorizontalMasks(square)
		VerticalMasks[square] = createVerticalMasks(square)
		DiagonalMasks[square] = createDiagonalMask(square)
		AntiDiagonalMasks[square] = createAntiDiagonalMask(square)
	}
}

func createWhitePawnAttacksForSquare(square Square) Bitboard {
	board := SquareMasks[square]
	rightAttack := (board << NorthEast) & ^FileA
	leftAttack := (board << NorthWest) & ^FileH
	return rightAttack | leftAttack
}

func createBlackPawnAttacksForSquare(square Square) Bitboard {
	board := SquareMasks[square]
	rightAttack := (board >> SouthEast) & ^FileA
	leftAttack := (board >> SouthWest) & ^FileH
	return rightAttack | leftAttack
}

func createKingMovesForSquare(square Square) Bitboard {

	startingSquare := SquareMasks[square]

	north := startingSquare << North
	northEast := startingSquare << NorthEast & ^FileA
	east := startingSquare >> East & ^FileA
	southEast := startingSquare >> SouthEast & ^FileA
	south := startingSquare >> South
	southWest := startingSquare >> SouthWest & ^FileH
	west := startingSquare << West & ^FileH
	northWest := startingSquare << NorthWest & ^FileH

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

	return (Diagonal >> south << north) ^ SquareMasks[square]
}

// createAntiDiagonalMask takes in a square and returns the bitboard that masks the antidiagonal lines
// (positive and megative) from that square. The formula for calculating this comes from the link above.
func createAntiDiagonalMask(square Square) Bitboard {
	antiDiagonal := 8*(int(square)&7) - (int(square) & 56)
	north := -antiDiagonal & (antiDiagonal >> 31)
	south := antiDiagonal & (-antiDiagonal >> 31)

	return (AntiDiagonal >> south << north) ^ SquareMasks[square]
}

func createHorizontalMasks(square Square) Bitboard {
	return (Bitboard(0xff) << (square & 56)) ^ SquareMasks[square]
}

func createVerticalMasks(square Square) Bitboard {
	return (Bitboard(0x0101010101010101) << (square & 7)) ^ SquareMasks[square]
}
