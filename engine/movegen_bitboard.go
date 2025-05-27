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
	bmg.moves = bmg.moves[:0]
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
		moves = bmg.board.GetRookMoves(from, activeSide)
	case Bishop:
		moves = bmg.board.GetBishopMoves(from, activeSide)
	case Queen:
		moves = bmg.generateQueenMoves(from, activeSide)
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

		for singleMove != 0 {
			to := singleMove.PopLSB()

			if isPromotion(to, activeSide) {
				bmg.addPromotionMoves(from, to, false)
				continue
			}

			move := NewMove(from, to, Pawn, Quiet)
			bmg.addMove(move)
		}

		for doubleMove != 0 {
			to := doubleMove.PopLSB()

			move := NewMove(from, to, Pawn, DoublePush)
			bmg.addMove(move)
		}

		targets := bmg.board.GetAllPiecesByColor(activeSide.EnemyColor()) | SquareMasks[enPassant]
		pawnAttacks := PawnAttacks[activeSide][from] & targets

		for pawnAttacks != 0 {
			var move Move
			to := pawnAttacks.PopLSB()

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
		bmg.addMove(NewPromotionMove(from, to, Pawn, Knight, CapturePromotion))
		bmg.addMove(NewPromotionMove(from, to, Pawn, Bishop, CapturePromotion))
		bmg.addMove(NewPromotionMove(from, to, Pawn, Rook, CapturePromotion))
		bmg.addMove(NewPromotionMove(from, to, Pawn, Queen, CapturePromotion))
	} else {
		bmg.addMove(NewPromotionMove(from, to, Pawn, Knight, Promotion))
		bmg.addMove(NewPromotionMove(from, to, Pawn, Bishop, Promotion))
		bmg.addMove(NewPromotionMove(from, to, Pawn, Rook, Promotion))
		bmg.addMove(NewPromotionMove(from, to, Pawn, Queen, Promotion))
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
		bmg.addMove(NewMove(E1, G1, King, CastleKingside))
	}

	if (castleAvailability&QueensideWhiteCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(B1, White) &&
		!bmg.board.SquareIsUnderAttack(C1, White) &&
		!bmg.board.SquareIsUnderAttack(D1, White) &&
		(occupied&B1C1D1Mask) == 0 {
		bmg.addMove(NewMove(E1, C1, King, CastleQueenside))
	}
}

func (bmg *BitboardMoveGenerator) generateBlackCastles(occupied Bitboard, castleAvailability CastleAvailability) {
	if (castleAvailability&KingsideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(F8, Black) &&
		!bmg.board.SquareIsUnderAttack(G8, Black) &&
		(occupied&F8G8Mask) == 0 {
		bmg.addMove(NewMove(E8, G8, King, CastleKingside))
	}

	if (castleAvailability&QueensideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(B8, Black) &&
		!bmg.board.SquareIsUnderAttack(C8, Black) &&
		!bmg.board.SquareIsUnderAttack(D8, Black) &&
		(occupied&B8C8D8Mask) == 0 {
		bmg.addMove(NewMove(E8, C8, King, CastleQueenside))
	}
}

func (bmg BitboardMoveGenerator) generateQueenMoves(from Square, activeSide Color) Bitboard {
	return bmg.board.GetBishopMoves(from, activeSide) | bmg.board.GetRookMoves(from, activeSide)
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
