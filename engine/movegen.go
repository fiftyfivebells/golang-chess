package engine

type MoveGen struct {
	board *Board
	moves []Move
}

func NewMoveGen(board *Board) *MoveGen {
	return &MoveGen{
		board: board,
		moves: []Move{},
	}
}

func (bmg *MoveGen) addMove(move Move) {
	bmg.moves = append(bmg.moves, move)
}

func (bmg MoveGen) GetMoves() []Move {
	return bmg.moves
}

func (bmg *MoveGen) GenerateMoves(activeSide Color, enPassant Square, castleAvailability CastleAvailability) {
	bmg.moves = bmg.moves[:0]
	targets := bmg.board.GetAllPiecesByColor(activeSide.EnemyColor())
	occupied := bmg.board.getAllPieces()
	activePieces := bmg.board.GetAllPiecesByColor(activeSide)

	for pieceType := Knight; pieceType < None; pieceType++ {
		pieceBoard := bmg.board.getPiecesByColorAndType(activeSide, pieceType)

		for pieceBoard != 0 {
			square := pieceBoard.PopLSB()
			bmg.generateMovesByPiece(pieceType, square, occupied, activePieces, targets)
		}
	}

	bmg.generatePawnMoves(activeSide, enPassant, occupied)
	bmg.generateCastlingMoves(activeSide, castleAvailability, occupied)
}

func (bmg *MoveGen) generateMovesByPiece(pieceType PieceType, from Square, occupied, allies, targets Bitboard) {

	var moves Bitboard
	switch pieceType {
	case Knight:
		moves = (KnightMoves[from] & ^allies)
	case Rook:
		moves = bmg.board.GetRookMoves(from, occupied, allies)
	case Bishop:
		moves = bmg.board.GetBishopMoves(from, occupied, allies)
	case Queen:
		moves = bmg.generateQueenMoves(from, occupied, allies)
	case King:
		moves = (KingMoves[from] & ^allies)
	}

	bmg.createMovesFromBitboard(from, moves, targets, pieceType)
}

func (bmg *MoveGen) generatePawnMoves(activeSide Color, enPassant Square, occupied Bitboard) {
	pawns := bmg.board.getPiecesByColorAndType(activeSide, Pawn)
	empty := ^occupied
	enemies := bmg.board.GetAllPiecesByColor(activeSide.EnemyColor())
	epMask := SquareMasks[enPassant] // 0 when enPassant == NoSquare
	captureTargets := enemies | epMask

	if activeSide == White {
		singles := (pawns << 8) & empty
		doubles := ((singles & Rank3) << 8) & empty

		for bb := singles & ^Rank8; bb != 0; {
			to := bb.PopLSB()
			bmg.addMove(NewMove(Square(int(to)-8), to, Pawn, Quiet))
		}
		for bb := singles & Rank8; bb != 0; {
			to := bb.PopLSB()
			bmg.addPromotionMoves(Square(int(to)-8), to, false)
		}
		for bb := doubles; bb != 0; {
			to := bb.PopLSB()
			bmg.addMove(NewMove(Square(int(to)-16), to, Pawn, DoublePush))
		}

		// Left = toward A-file (<<7), Right = toward H-file (<<9)
		leftAttacks := (pawns << 7) & ^FileA & captureTargets
		rightAttacks := (pawns << 9) & ^FileH & captureTargets

		for bb := leftAttacks & ^Rank8; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) - 7)
			if to == enPassant {
				bmg.addMove(NewMove(from, to, Pawn, EnPassant))
			} else {
				bmg.addMove(NewMove(from, to, Pawn, Capture))
			}
		}
		for bb := leftAttacks & Rank8; bb != 0; {
			to := bb.PopLSB()
			bmg.addPromotionMoves(Square(int(to)-7), to, true)
		}
		for bb := rightAttacks & ^Rank8; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) - 9)
			if to == enPassant {
				bmg.addMove(NewMove(from, to, Pawn, EnPassant))
			} else {
				bmg.addMove(NewMove(from, to, Pawn, Capture))
			}
		}
		for bb := rightAttacks & Rank8; bb != 0; {
			to := bb.PopLSB()
			bmg.addPromotionMoves(Square(int(to)-9), to, true)
		}

	} else {
		singles := (pawns >> 8) & empty
		doubles := ((singles & Rank6) >> 8) & empty

		for bb := singles & ^Rank1; bb != 0; {
			to := bb.PopLSB()
			bmg.addMove(NewMove(Square(int(to)+8), to, Pawn, Quiet))
		}
		for bb := singles & Rank1; bb != 0; {
			to := bb.PopLSB()
			bmg.addPromotionMoves(Square(int(to)+8), to, false)
		}
		for bb := doubles; bb != 0; {
			to := bb.PopLSB()
			bmg.addMove(NewMove(Square(int(to)+16), to, Pawn, DoublePush))
		}

		// Left = toward H-file (>>7), Right = toward A-file (>>9)
		leftAttacks := (pawns >> 7) & ^FileH & captureTargets
		rightAttacks := (pawns >> 9) & ^FileA & captureTargets

		for bb := leftAttacks & ^Rank1; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) + 7)
			if to == enPassant {
				bmg.addMove(NewMove(from, to, Pawn, EnPassant))
			} else {
				bmg.addMove(NewMove(from, to, Pawn, Capture))
			}
		}
		for bb := leftAttacks & Rank1; bb != 0; {
			to := bb.PopLSB()
			bmg.addPromotionMoves(Square(int(to)+7), to, true)
		}
		for bb := rightAttacks & ^Rank1; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) + 9)
			if to == enPassant {
				bmg.addMove(NewMove(from, to, Pawn, EnPassant))
			} else {
				bmg.addMove(NewMove(from, to, Pawn, Capture))
			}
		}
		for bb := rightAttacks & Rank1; bb != 0; {
			to := bb.PopLSB()
			bmg.addPromotionMoves(Square(int(to)+9), to, true)
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

func (bmg *MoveGen) addPromotionMoves(from, to Square, isCapture bool) {

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

func (bmg *MoveGen) generateCastlingMoves(activeSide Color, castleAvailability CastleAvailability, occupied Bitboard) {
	if !bmg.board.KingIsUnderAttack(activeSide) {
		if activeSide == White {
			bmg.generateWhiteCastles(occupied, castleAvailability)
		} else if activeSide == Black {
			bmg.generateBlackCastles(occupied, castleAvailability)
		}
	}

}

func (bmg *MoveGen) generateWhiteCastles(occupied Bitboard, castleAvailability CastleAvailability) {
	if (castleAvailability&KingsideWhiteCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(F1, White) &&
		!bmg.board.SquareIsUnderAttack(G1, White) &&
		(occupied&F1G1Mask) == 0 {
		bmg.addMove(NewMove(E1, G1, King, CastleKingside))
	}

	if (castleAvailability&QueensideWhiteCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(C1, White) &&
		!bmg.board.SquareIsUnderAttack(D1, White) &&
		(occupied&B1C1D1Mask) == 0 {
		bmg.addMove(NewMove(E1, C1, King, CastleQueenside))
	}
}

func (bmg *MoveGen) generateBlackCastles(occupied Bitboard, castleAvailability CastleAvailability) {
	if (castleAvailability&KingsideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(F8, Black) &&
		!bmg.board.SquareIsUnderAttack(G8, Black) &&
		(occupied&F8G8Mask) == 0 {
		bmg.addMove(NewMove(E8, G8, King, CastleKingside))
	}

	if (castleAvailability&QueensideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(C8, Black) &&
		!bmg.board.SquareIsUnderAttack(D8, Black) &&
		(occupied&B8C8D8Mask) == 0 {
		bmg.addMove(NewMove(E8, C8, King, CastleQueenside))
	}
}

func (bmg MoveGen) generateQueenMoves(from Square, occupied, allies Bitboard) Bitboard {
	return bmg.board.GetBishopMoves(from, occupied, allies) | bmg.board.GetRookMoves(from, occupied, allies)
}

func (bmg *MoveGen) createMovesFromBitboard(from Square, moves, targets Bitboard, pieceType PieceType) {

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
