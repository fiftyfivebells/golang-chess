package engine

type MoveGen struct {
	board *Board
}

func NewMoveGen(board *Board) *MoveGen {
	return &MoveGen{
		board: board,
	}
}

func (bmg *MoveGen) GenerateMoves(buf *[256]Move, activeSide Color, enPassant Square, castleAvailability CastleAvailability) int {
	targets := bmg.board.GetAllPiecesByColor(activeSide.EnemyColor())
	occupied := bmg.board.getAllPieces()
	activePieces := bmg.board.GetAllPiecesByColor(activeSide)

	count := 0

	for pieceType := Knight; pieceType < None; pieceType++ {
		pieceBoard := bmg.board.getPiecesByColorAndType(activeSide, pieceType)

		for pieceBoard != 0 {
			square := pieceBoard.PopLSB()
			count = bmg.generateMovesByPiece(buf, count, pieceType, square, occupied, activePieces, targets)
		}
	}

	count = bmg.generatePawnMoves(buf, count, activeSide, enPassant, occupied)
	count = bmg.generateCastlingMoves(buf, count, activeSide, castleAvailability, occupied)

	return count
}

func (bmg *MoveGen) generateMovesByPiece(
	buf *[256]Move,
	count int,
	pieceType PieceType,
	from Square,
	occupied, allies, targets Bitboard,
) int {

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

	return bmg.createMovesFromBitboard(buf, count, from, moves, targets, pieceType)
}

func (bmg *MoveGen) generatePawnMoves(buf *[256]Move, count int, activeSide Color, enPassant Square, occupied Bitboard) int {
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
			buf[count] = (NewMove(Square(int(to)-8), to, Pawn, Quiet))
			count++
		}
		for bb := singles & Rank8; bb != 0; {
			to := bb.PopLSB()
			count = bmg.addPromotionMoves(buf, count, Square(int(to)-8), to, false)
		}
		for bb := doubles; bb != 0; {
			to := bb.PopLSB()
			buf[count] = (NewMove(Square(int(to)-16), to, Pawn, DoublePush))
			count++
		}

		// Left = toward A-file (<<7), Right = toward H-file (<<9)
		leftAttacks := (pawns << 7) & ^FileA & captureTargets
		rightAttacks := (pawns << 9) & ^FileH & captureTargets

		for bb := leftAttacks & ^Rank8; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) - 7)
			if to == enPassant {
				buf[count] = (NewMove(from, to, Pawn, EnPassant))
				count++
			} else {
				buf[count] = (NewMove(from, to, Pawn, Capture))
				count++
			}
		}
		for bb := leftAttacks & Rank8; bb != 0; {
			to := bb.PopLSB()
			count = bmg.addPromotionMoves(buf, count, Square(int(to)-7), to, true)
		}
		for bb := rightAttacks & ^Rank8; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) - 9)
			if to == enPassant {
				buf[count] = (NewMove(from, to, Pawn, EnPassant))
				count++
			} else {
				buf[count] = (NewMove(from, to, Pawn, Capture))
				count++
			}
		}
		for bb := rightAttacks & Rank8; bb != 0; {
			to := bb.PopLSB()
			count = bmg.addPromotionMoves(buf, count, Square(int(to)-9), to, true)
		}

	} else {
		singles := (pawns >> 8) & empty
		doubles := ((singles & Rank6) >> 8) & empty

		for bb := singles & ^Rank1; bb != 0; {
			to := bb.PopLSB()
			buf[count] = (NewMove(Square(int(to)+8), to, Pawn, Quiet))
			count++
		}
		for bb := singles & Rank1; bb != 0; {
			to := bb.PopLSB()
			count = bmg.addPromotionMoves(buf, count, Square(int(to)+8), to, false)
		}
		for bb := doubles; bb != 0; {
			to := bb.PopLSB()
			buf[count] = (NewMove(Square(int(to)+16), to, Pawn, DoublePush))
			count++
		}

		// Left = toward H-file (>>7), Right = toward A-file (>>9)
		leftAttacks := (pawns >> 7) & ^FileH & captureTargets
		rightAttacks := (pawns >> 9) & ^FileA & captureTargets

		for bb := leftAttacks & ^Rank1; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) + 7)
			if to == enPassant {
				buf[count] = (NewMove(from, to, Pawn, EnPassant))
				count++
			} else {
				buf[count] = (NewMove(from, to, Pawn, Capture))
				count++
			}
		}
		for bb := leftAttacks & Rank1; bb != 0; {
			to := bb.PopLSB()
			count = bmg.addPromotionMoves(buf, count, Square(int(to)+7), to, true)
		}
		for bb := rightAttacks & ^Rank1; bb != 0; {
			to := bb.PopLSB()
			from := Square(int(to) + 9)
			if to == enPassant {
				buf[count] = (NewMove(from, to, Pawn, EnPassant))
				count++
			} else {
				buf[count] = (NewMove(from, to, Pawn, Capture))
				count++
			}
		}
		for bb := rightAttacks & Rank1; bb != 0; {
			to := bb.PopLSB()
			count = bmg.addPromotionMoves(buf, count, Square(int(to)+9), to, true)
		}
	}

	return count
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

func (bmg *MoveGen) addPromotionMoves(buf *[256]Move, count int, from, to Square, isCapture bool) int {
	if isCapture {
		buf[count] = (NewPromotionMove(from, to, Pawn, Knight, CapturePromotion))
		count++
		buf[count] = (NewPromotionMove(from, to, Pawn, Bishop, CapturePromotion))
		count++
		buf[count] = (NewPromotionMove(from, to, Pawn, Rook, CapturePromotion))
		count++
		buf[count] = (NewPromotionMove(from, to, Pawn, Queen, CapturePromotion))
		count++
	} else {
		buf[count] = (NewPromotionMove(from, to, Pawn, Knight, Promotion))
		count++
		buf[count] = (NewPromotionMove(from, to, Pawn, Bishop, Promotion))
		count++
		buf[count] = (NewPromotionMove(from, to, Pawn, Rook, Promotion))
		count++
		buf[count] = (NewPromotionMove(from, to, Pawn, Queen, Promotion))
		count++
	}

	return count
}

func (bmg *MoveGen) generateCastlingMoves(
	buf *[256]Move,
	count int,
	activeSide Color,
	castleAvailability CastleAvailability,
	occupied Bitboard,
) int {
	if !bmg.board.KingIsUnderAttack(activeSide) {
		if activeSide == White {
			count = bmg.generateWhiteCastles(buf, count, occupied, castleAvailability)
		} else if activeSide == Black {
			count = bmg.generateBlackCastles(buf, count, occupied, castleAvailability)
		}
	}

	return count
}

func (bmg *MoveGen) generateWhiteCastles(buf *[256]Move, count int, occupied Bitboard, castleAvailability CastleAvailability) int {

	if (castleAvailability&KingsideWhiteCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(F1, White) &&
		!bmg.board.SquareIsUnderAttack(G1, White) &&
		(occupied&F1G1Mask) == 0 {
		buf[count] = (NewMove(E1, G1, King, CastleKingside))
		count++
	}

	if (castleAvailability&QueensideWhiteCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(C1, White) &&
		!bmg.board.SquareIsUnderAttack(D1, White) &&
		(occupied&B1C1D1Mask) == 0 {
		buf[count] = (NewMove(E1, C1, King, CastleQueenside))
		count++
	}

	return count
}

func (bmg *MoveGen) generateBlackCastles(buf *[256]Move, count int, occupied Bitboard, castleAvailability CastleAvailability) int {

	if (castleAvailability&KingsideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(F8, Black) &&
		!bmg.board.SquareIsUnderAttack(G8, Black) &&
		(occupied&F8G8Mask) == 0 {
		buf[count] = (NewMove(E8, G8, King, CastleKingside))
		count++
	}

	if (castleAvailability&QueensideBlackCastle) != 0 &&
		!bmg.board.SquareIsUnderAttack(C8, Black) &&
		!bmg.board.SquareIsUnderAttack(D8, Black) &&
		(occupied&B8C8D8Mask) == 0 {
		buf[count] = (NewMove(E8, C8, King, CastleQueenside))
		count++
	}

	return count
}

func (bmg MoveGen) generateQueenMoves(from Square, occupied, allies Bitboard) Bitboard {
	return bmg.board.GetBishopMoves(from, occupied, allies) | bmg.board.GetRookMoves(from, occupied, allies)
}

func (bmg *MoveGen) createMovesFromBitboard(
	buf *[256]Move,
	count int,
	from Square,
	moves, targets Bitboard,
	pieceType PieceType,
) int {

	for moves != 0 {
		to := moves.PopLSB()
		toBoard := SquareMasks[to]

		moveType := Quiet

		if (toBoard & targets) != 0 {
			moveType = Capture
		}

		move := NewMove(from, to, pieceType, moveType)
		buf[count] = move
		count++
	}

	return count
}
