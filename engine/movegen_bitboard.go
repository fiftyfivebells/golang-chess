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

func (bmg *BitboardMoveGenerator) GenerateMoves(activeSide Color, enPassant Square) {
	activePieces := bmg.board.getAllPiecesByColor(activeSide)
	targets := bmg.board.getAllPiecesByColor(activeSide.EnemyColor())

	for pieceType := Knight; pieceType < None; pieceType++ {
		pieceBoard := bmg.board.getPiecesByColorAndType(activeSide, pieceType)

		for pieceBoard != 0 {
			square := pieceBoard.PopLSB()
			bmg.generateMovesByPiece(pieceType, square, activePieces, targets)
		}
	}

	bmg.generatePawnMoves(activeSide, enPassant)
	bmg.generateCastlingMoves(activeSide)
}

func (bmg *BitboardMoveGenerator) generateMovesByPiece(pieceType PieceType, from Square, activePieces, targets Bitboard) {
	var moves Bitboard
	switch pieceType {
	case Knight:
		moves = (KnightMoves[from] & ^activePieces)
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

		targets := bmg.board.getAllPiecesByColor(activeSide.EnemyColor()) | SquareMasks[enPassant]
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

func (bmg *BitboardMoveGenerator) generateCastlingMoves(activeSide Color) {}

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
