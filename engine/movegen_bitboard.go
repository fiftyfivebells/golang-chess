package engine

type BitboardMoveGenerator struct {
	board *BitboardBoard
	moves []Move
}

func NewBitboardMoveGenerator(board *BitboardBoard) MoveGenerator {
	return &BitboardMoveGenerator{
		board: board,
		moves: []Move{},
	}
}

func (bmg *BitboardMoveGenerator) addMove(move Move) {
	bmg.moves = append(bmg.moves, move)
}

func (bmg BitboardMoveGenerator) GetMoves() []Move {
	return bmg.moves
}

func (bmg *BitboardMoveGenerator) GenerateMoves(activeSide Color, enPassant Square, castleAvailability CastleAvailability) {
	targets := bmg.board.GetAllPiecesByColor(activeSide.EnemyColor())

	for pieceType := Knight; pieceType < None; pieceType++ {
		pieceBoard := bmg.board.getPiecesByColorAndType(activeSide, pieceType)

		for pieceBoard != 0 {
			square := pieceBoard.PopLSB()
			bmg.generateMovesByPiece(pieceType, square, activeSide, targets)
		}
	}

	bmg.generatePawnMoves(activeSide, enPassant)
	bmg.generateCastlingMoves(activeSide, castleAvailability)
}

func (bmg *BitboardMoveGenerator) generateMovesByPiece(pieceType PieceType, from Square, activeSide Color, targets Bitboard) {
	activePieces := bmg.board.GetAllPiecesByColor(activeSide)

	var moves Bitboard
	switch pieceType {
	case Knight:
		moves = (KnightMoves[from] & ^activePieces)
	case Rook:
		moves = bmg.generateRookMoves(from, activeSide)
	case Bishop:
		moves = bmg.generateBishopMoves(from, activeSide)
	case King:
		moves = (KingMoves[from] & ^activePieces)
	}

	bmg.createMovesFromBitboard(from, moves, targets, pieceType)
}

func (bmg *BitboardMoveGenerator) generatePawnMoves(activeSide Color, enPassant Square) {
	allPieces := bmg.board.getAllPieces()
	pawns := bmg.board.getPiecesByColorAndType(activeSide, Pawn)

	for pawns != 0 {
		from := pawns.PopLSB()
		singleMove := PawnPushes[activeSide][from] & ^allPieces
		doubleMove := ((singleMove & Rank3) << North) & ^allPieces
		if activeSide == Black {
			doubleMove = ((singleMove & Rank6) >> South) & ^allPieces
		}

		moves := singleMove | doubleMove
		for moves != 0 {
			to := moves.PopLSB()

			if isPromotion(to, activeSide) {
				bmg.addPromotionMoves(from, to, false)
				continue
			}

			move := NewMove(from, to, Pawn, Quiet)
			bmg.addMove(move)
		}

		targets := bmg.board.GetAllPiecesByColor(activeSide.EnemyColor()) | SquareMasks[enPassant]
		pawnAttacks := PawnAttacks[activeSide][from] & targets

		for pawnAttacks != 0 {
			var move Move
			to := moves.PopLSB()

			if isPromotion(to, activeSide) {
				bmg.addPromotionMoves(from, to, true)
				continue
			}

			if to == enPassant {
				move = NewMove(from, to, Pawn, EnPassant)
			} else {
				move = NewMove(from, to, Pawn, Capture)
			}

			bmg.addMove(move)
		}
	}
}

func isPromotion(to Square, color Color) bool {
	switch color {
	case White:
		return to >= H8 && to <= A8
	case Black:
		return to >= H1 && to <= A1
	default:
		return false
	}
}

func (bmg *BitboardMoveGenerator) addPromotionMoves(from, to Square, isCapture bool) {

	if isCapture {
		bmg.addMove(NewMove(from, to, Pawn, CapturePromotionKnight))
		bmg.addMove(NewMove(from, to, Pawn, CapturePromotionBishop))
		bmg.addMove(NewMove(from, to, Pawn, CapturePromotionRook))
		bmg.addMove(NewMove(from, to, Pawn, CapturePromotionQueen))
	} else {
		bmg.addMove(NewMove(from, to, Pawn, PromotionKnight))
		bmg.addMove(NewMove(from, to, Pawn, PromotionBishop))
		bmg.addMove(NewMove(from, to, Pawn, PromotionRook))
		bmg.addMove(NewMove(from, to, Pawn, PromotionQueen))
	}
}

func (bmg *BitboardMoveGenerator) generateCastlingMoves(activeSide Color, castleAvailability CastleAvailability) {
	occupied := bmg.board.getAllPieces()

	if activeSide == White {
		bmg.generateWhiteCastles(occupied, castleAvailability)
	} else if activeSide == Black {
		bmg.generateBlackCastles(occupied, castleAvailability)
	}
}

func (bmg *BitboardMoveGenerator) generateWhiteCastles(occupied Bitboard, castleAvailability CastleAvailability) {
	if (castleAvailability&KingsideWhiteCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(F1, White) &&
		!bmg.board.SquareIsUnderAttack(G1, White) &&
		(occupied&F1G1Mask) == 0 {
		bmg.addMove(NewMove(E1, H1, King, CastleKingside))
	}

	if (castleAvailability&QueensideWhiteCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(B1, White) &&
		!bmg.board.SquareIsUnderAttack(C1, White) &&
		!bmg.board.SquareIsUnderAttack(D1, White) &&
		(occupied&B1C1D1Mask) == 0 {
		bmg.addMove(NewMove(E1, A1, King, CastleQueenside))
	}
}

func (bmg *BitboardMoveGenerator) generateBlackCastles(occupied Bitboard, castleAvailability CastleAvailability) {
	if (castleAvailability&KingsideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(F8, Black) &&
		!bmg.board.SquareIsUnderAttack(G8, Black) &&
		(occupied&F8G8Mask) == 0 {
		bmg.addMove(NewMove(E8, H8, King, CastleKingside))
	}

	if (castleAvailability&QueensideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(B8, White) &&
		!bmg.board.SquareIsUnderAttack(C8, White) &&
		!bmg.board.SquareIsUnderAttack(D8, White) &&
		(occupied&B8C8D8Mask) == 0 {
		bmg.addMove(NewMove(E8, A8, King, CastleQueenside))
	}
}

func (bmg BitboardMoveGenerator) generateBishopMoves(from Square, activeSide Color) Bitboard {
	allies := bmg.board.GetAllPiecesByColor(activeSide)

	diagonalMask := DiagonalMasks[from]
	antiDiagonalMask := AntiDiagonalMasks[from]

	diagonal := bmg.generateSlidingMoves(activeSide, from, diagonalMask)
	antiDiagonal := bmg.generateSlidingMoves(activeSide, from, antiDiagonalMask)

	return (diagonal | antiDiagonal) & ^allies
}

func (bmg BitboardMoveGenerator) generateRookMoves(from Square, activeSide Color) Bitboard {
	allies := bmg.board.GetAllPiecesByColor(activeSide)

	rank := RankMaskForSquare(from)
	file := FileMaskForSquare(from)

	horizontal := bmg.generateSlidingMoves(activeSide, from, rank)
	vertical := bmg.generateSlidingMoves(activeSide, from, file)

	return (horizontal | vertical) & ^allies
}

func (bmg BitboardMoveGenerator) generateSlidingMoves(activeSide Color, square Square, mask Bitboard) Bitboard {
	squareBoard := SquareMasks[square]
	occupied := bmg.board.getAllPieces()

	bottom := ((occupied & mask) - (squareBoard << 1)) & mask
	top := ReverseBitboard(ReverseBitboard((occupied & mask)) - 2*ReverseBitboard(squareBoard))

	allies := bmg.board.GetAllPiecesByColor(activeSide)

	return (bottom ^ top) & mask & ^allies
}

func (bmg *BitboardMoveGenerator) createMovesFromBitboard(from Square, moves, targets Bitboard, pieceType PieceType) {

	for moves != 0 {
		to := moves.PopLSB()
		toBoard := SquareMasks[to]

		moveType := Quiet

		if (toBoard & targets) != 0 {
			moveType = Capture
		}

		move := NewMove(from, to, pieceType, moveType)
		bmg.addMove(move)
	}
}
