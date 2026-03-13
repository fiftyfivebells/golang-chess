package engine

import (
	"fmt"
	"strconv"
	"strings"
)

type GameState struct {
	Board        Board
	ActiveSide   Color
	CastleRights CastleAvailability
	EPSquare     Square
	HalfMove     uint16
	FullMove     byte

	LegalMovesBuffer [256]Move
	StatePly         uint16
	PreviousStates   [100]IrreversibleState

	moveGen MoveGen
}

type IrreversibleState struct {
	CastleRights CastleAvailability
	EPSquare     Square
	HalfMove     uint16
	Destination  Piece
}

func InitializeGameState(fen string) *GameState {
	gs := &GameState{FullMove: 1}
	gs.moveGen = *NewMoveGen(&gs.Board)
	gs.SetStateFromFENString(fen)
	return gs
}

func (gs *GameState) ClearGameState() {
	gs.Board.ClearBoard()
	gs.CastleRights = 0
	gs.EPSquare = NoSquare
	gs.HalfMove = 0
	gs.FullMove = 1

	gs.StatePly = 0
}

func (gs *GameState) GetPseudoLegalMovesForPosition() []Move {

	count := gs.moveGen.GenerateMoves(&gs.LegalMovesBuffer, gs.ActiveSide, gs.EPSquare, gs.CastleRights)

	return gs.LegalMovesBuffer[:count]
}

func (gs *GameState) GetLegalMovesForPosition() []Move {
	pseudoLegalMoves := gs.GetPseudoLegalMovesForPosition()
	legalMoves := gs.LegalMovesBuffer[:0]
	for _, move := range pseudoLegalMoves {
		if gs.ApplyMove(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

func (gs *GameState) ApplyMove(move Move) bool {
	moveType := move.MoveType()
	from := move.FromSquare()
	to := move.ToSquare()

	color := gs.ActiveSide
	pt := move.PieceType()
	movingPiece := makePiece(pt, color)

	previous := IrreversibleState{
		CastleRights: gs.CastleRights,
		EPSquare:     gs.EPSquare,
		HalfMove:     gs.HalfMove,
		Destination:  NoPiece,
	}

	gs.HalfMove++
	gs.EPSquare = NoSquare

	switch moveType {
	case Quiet:
		mask := SquareMasks[from] | SquareMasks[to]
		gs.Board.pieces[color][pt] ^= mask
		gs.Board.colorBB[color] ^= mask
		gs.Board.occupancy ^= mask
		gs.Board.squares[to] = movingPiece
		gs.Board.squares[from] = NoPiece
		if pt == King {
			gs.Board.kingSq[color] = to
		}

	case Capture:
		previous.Destination = gs.Board.squares[to]
		captured := previous.Destination
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]
		gs.Board.pieces[color][pt] ^= fromBB | toBB
		gs.Board.colorBB[color] ^= fromBB | toBB
		gs.Board.pieces[captured.Color()][captured.Type()] &^= toBB
		gs.Board.colorBB[captured.Color()] &^= toBB
		gs.Board.occupancy &^= fromBB // to stays occupied, now has moving piece
		gs.Board.squares[to] = movingPiece
		gs.Board.squares[from] = NoPiece
		if pt == King {
			gs.Board.kingSq[color] = to
		}

	case DoublePush:
		mask := SquareMasks[from] | SquareMasks[to]
		gs.Board.pieces[color][pt] ^= mask
		gs.Board.colorBB[color] ^= mask
		gs.Board.occupancy ^= mask
		gs.Board.squares[to] = movingPiece
		gs.Board.squares[from] = NoPiece
		gs.EPSquare = Square(int(from) + pawnDirection[gs.ActiveSide])

	case CapturePromotion:
		previous.Destination = gs.Board.squares[to]
		captured := previous.Destination
		promoPiece := makePiece(move.PromotionPieceType(), color)
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]

		// Remove captured piece from to
		gs.Board.pieces[captured.Color()][captured.Type()] &^= toBB
		gs.Board.colorBB[captured.Color()] &^= toBB

		// Remove pawn from from
		gs.Board.pieces[color][Pawn] &^= fromBB

		// Place promotion piece at to
		gs.Board.pieces[color][promoPiece.Type()] |= toBB

		// Color: from cleared, to already had enemy color cleared above, now set ours
		gs.Board.colorBB[color] ^= fromBB | toBB

		// Occupancy: from cleared (pawn left), to stays occupied (captured removed, promo placed)
		gs.Board.occupancy &^= fromBB

		gs.Board.squares[from] = NoPiece
		gs.Board.squares[to] = promoPiece

	case Promotion:
		promoPiece := makePiece(move.PromotionPieceType(), color)
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]

		// Remove pawn from from
		gs.Board.pieces[color][Pawn] &^= fromBB

		// Place promotion piece at to
		gs.Board.pieces[color][promoPiece.Type()] |= toBB

		// pawn leaves from, promo arrives at to
		gs.Board.colorBB[color] ^= fromBB | toBB

		// clear from, set to
		gs.Board.occupancy ^= fromBB | toBB

		gs.Board.squares[from] = NoPiece
		gs.Board.squares[to] = promoPiece

	case CastleKingside, CastleQueenside:
		rookFrom, rookTo := CastlingRookPositions(from, to)

		kingMask := SquareMasks[from] | SquareMasks[to]
		rookMask := SquareMasks[rookFrom] | SquareMasks[rookTo]

		// Move king
		gs.Board.pieces[color][King] ^= kingMask
		// Move rook
		gs.Board.pieces[color][Rook] ^= rookMask

		gs.Board.colorBB[color] ^= kingMask | rookMask
		gs.Board.occupancy ^= kingMask | rookMask

		// Squares array
		rook := makePiece(Rook, color)
		gs.Board.squares[to] = movingPiece   // king destination
		gs.Board.squares[from] = NoPiece     // king origin
		gs.Board.squares[rookTo] = rook      // rook destination
		gs.Board.squares[rookFrom] = NoPiece // rook origin

		gs.Board.kingSq[color] = to

	case EnPassant:
		dir := pawnDirection[gs.ActiveSide]
		capturedSq := Square(int(to) - dir)

		// save the captured pawn before it gets removed
		previous.Destination = gs.Board.squares[capturedSq]

		moveMask := SquareMasks[from] | SquareMasks[to]
		capMask := SquareMasks[capturedSq]
		enemyColor := gs.ActiveSide.EnemyColor()

		// move the pawn
		gs.Board.pieces[color][Pawn] ^= moveMask
		gs.Board.colorBB[color] ^= moveMask

		// remove the captured pawn
		gs.Board.pieces[enemyColor][Pawn] &^= capMask
		gs.Board.colorBB[enemyColor] &^= capMask

		// from is cleared, to is set, capturedSq is cleared
		gs.Board.occupancy ^= moveMask
		gs.Board.occupancy &^= capMask

		gs.Board.squares[to] = movingPiece
		gs.Board.squares[from] = NoPiece
		gs.Board.squares[capturedSq] = NoPiece

	}

	// Update castle rights
	if pt == King || pt == Rook || previous.Destination.Type() == Rook {
		gs.UpdateCastleRights(movingPiece, previous.Destination, move)
	}

	// The halfmove clock gets reset if the move was a capture or if the moved piece was a pawn
	if IsAttackMove(moveType) || pt == Pawn {
		gs.HalfMove = 0
	}

	gs.PreviousStates[gs.StatePly] = previous
	gs.StatePly++

	// The fullmove number is incremented only after the black side has moved
	if gs.ActiveSide == Black {
		gs.FullMove++
	}

	gs.ActiveSide = gs.ActiveSide.EnemyColor()

	return !gs.Board.KingIsUnderAttack(gs.ActiveSide.EnemyColor())
}

func (gs *GameState) UnapplyMove(move Move) {
	gs.StatePly--
	previous := gs.PreviousStates[gs.StatePly]

	gs.CastleRights = previous.CastleRights
	gs.EPSquare = previous.EPSquare
	gs.HalfMove = previous.HalfMove

	gs.ActiveSide = gs.ActiveSide.EnemyColor()
	if gs.ActiveSide == Black {
		gs.FullMove--
	}

	from := move.FromSquare()
	to := move.ToSquare()
	moveType := move.MoveType()

	color := gs.ActiveSide
	pt := move.PieceType()

	movingPiece := makePiece(pt, color)
	capturedPiece := previous.Destination
	capturedColor := capturedPiece.Color()

	switch moveType {
	case Quiet, DoublePush:
		// Identical undo: move piece from to back to from
		mask := SquareMasks[from] | SquareMasks[to]
		gs.Board.pieces[color][pt] ^= mask
		gs.Board.colorBB[color] ^= mask
		gs.Board.occupancy ^= mask
		gs.Board.squares[from] = movingPiece
		gs.Board.squares[to] = NoPiece
		if pt == King {
			gs.Board.kingSq[color] = from
		}

	case Capture:
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]

		// Move our piece back: to → from
		gs.Board.pieces[color][pt] ^= fromBB | toBB
		gs.Board.colorBB[color] ^= fromBB | toBB
		gs.Board.occupancy ^= fromBB // from was empty, now occupied; to stays occupied

		// Restore captured piece at to
		gs.Board.pieces[capturedColor][capturedPiece.Type()] |= toBB
		gs.Board.colorBB[capturedColor] |= toBB

		gs.Board.squares[from] = movingPiece
		gs.Board.squares[to] = capturedPiece
		if pt == King {
			gs.Board.kingSq[color] = from
		}

	case EnPassant:
		dir := pawnDirection[gs.ActiveSide]
		capturedSq := Square(int(to) - dir)

		moveMask := SquareMasks[from] | SquareMasks[to]
		capMask := SquareMasks[capturedSq]
		enemyColor := gs.ActiveSide.EnemyColor()

		// Move our pawn back: to → from
		gs.Board.pieces[color][Pawn] ^= moveMask
		gs.Board.colorBB[color] ^= moveMask
		gs.Board.occupancy ^= moveMask // toggles from (on) and to (off)

		// Restore captured pawn at capturedSq
		gs.Board.pieces[enemyColor][Pawn] |= capMask
		gs.Board.colorBB[enemyColor] |= capMask
		gs.Board.occupancy |= capMask

		gs.Board.squares[from] = movingPiece
		gs.Board.squares[to] = NoPiece
		gs.Board.squares[capturedSq] = capturedPiece

	case Promotion:
		promoPieceType := move.PromotionPieceType()
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]

		// Remove promotion piece from to
		gs.Board.pieces[color][promoPieceType] &^= toBB

		// Restore pawn at from
		gs.Board.pieces[color][Pawn] |= fromBB

		// Color: to cleared, from set
		gs.Board.colorBB[color] ^= fromBB | toBB

		// Occupancy: to cleared, from set
		gs.Board.occupancy ^= fromBB | toBB

		gs.Board.squares[from] = movingPiece
		gs.Board.squares[to] = NoPiece

	case CapturePromotion:
		promoPieceType := move.PromotionPieceType()
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]

		// Remove promotion piece from to
		gs.Board.pieces[color][promoPieceType] &^= toBB
		gs.Board.colorBB[color] &^= toBB

		// Restore pawn at from
		gs.Board.pieces[color][Pawn] |= fromBB
		gs.Board.colorBB[color] |= fromBB

		// Restore captured piece at to
		gs.Board.pieces[capturedColor][capturedPiece.Type()] |= toBB
		gs.Board.colorBB[capturedColor] |= toBB

		// Occupancy: from becomes occupied, to stays occupied
		gs.Board.occupancy |= fromBB

		gs.Board.squares[from] = movingPiece
		gs.Board.squares[to] = capturedPiece

	case CastleKingside, CastleQueenside:
		rookFrom, rookTo := CastlingRookPositions(from, to)

		kingMask := SquareMasks[from] | SquareMasks[to]
		rookMask := SquareMasks[rookFrom] | SquareMasks[rookTo]

		// Move king back: to → from
		gs.Board.pieces[color][King] ^= kingMask
		// Move rook back: rookTo → rookFrom
		gs.Board.pieces[color][Rook] ^= rookMask

		gs.Board.colorBB[color] ^= kingMask | rookMask
		gs.Board.occupancy ^= kingMask | rookMask

		rook := makePiece(Rook, color)
		gs.Board.squares[from] = movingPiece // king back to origin
		gs.Board.squares[to] = NoPiece       // king's destination cleared
		gs.Board.squares[rookFrom] = rook    // rook back to origin
		gs.Board.squares[rookTo] = NoPiece   // rook's destination cleared

		gs.Board.kingSq[color] = from
	}
}

func (gs *GameState) UpdateCastleRights(moved Piece, captured Piece, move Move) {
	if moved.Type() == King {
		gs.CastleRights.RemoveAllRights(gs.ActiveSide)
	} else if moved.Type() == Rook {
		gs.UpdateRookRights(gs.ActiveSide, move.FromSquare())
	} else if captured.Type() == Rook {
		gs.UpdateRookRights(gs.ActiveSide.EnemyColor(), move.ToSquare())
	}
}

func (gs *GameState) UpdateRookRights(color Color, square Square) {
	kingside, queenside := rookSquares(color)

	if square == kingside {
		gs.CastleRights.Remove(color, Kingside)
	} else if square == queenside {
		gs.CastleRights.Remove(color, Queenside)
	}
}

func rookSquares(color Color) (Square, Square) {
	if color == White {
		return H1, A1
	} else {
		return H8, A8
	}
}

func (gs GameState) GetGameStateFENString() string {
	var fenString strings.Builder

	fenString.WriteString(gs.Board.GetFENRepresentation())

	fenString.WriteString(" ")
	fenString.WriteString(gs.ActiveSide.String())

	fenString.WriteString(" ")
	fenString.WriteString(gs.CastleRights.String())

	fenString.WriteString(" ")
	if gs.EPSquare == NoSquare {
		fenString.WriteString("-")
	} else {
		fenString.WriteString(SquareToCoord(gs.EPSquare))
	}

	fenString.WriteString(" ")
	fenString.WriteString(strconv.Itoa(int(gs.HalfMove)))

	fenString.WriteString(" ")
	fenString.WriteString(strconv.Itoa(int(gs.FullMove)))

	return fenString.String()
}

func (gs *GameState) SetStateFromFENString(fenString string) {
	fenValues := strings.Fields(fenString)
	// if the fen string is not valid (doesn't have 6 fields, anyway), just set to initial state
	if len(fenValues) != 6 {
		fenValues = strings.Fields(InitialStateFenString)
	}

	pieces := fenValues[0]
	activeSide := fenValues[1]
	castleAvailability := fenValues[2]
	enPassantSquare := fenValues[3]
	halfMove, _ := strconv.Atoi(fenValues[4])
	fullMove, _ := strconv.Atoi(fenValues[5])

	gs.Board.SetBoardFromFEN(pieces)
	// if gs.Board != nil {
	// 	gs.Board.SetBoardFromFEN(pieces)
	// } else {
	// 	gs.Board = NewBoard(pieces)
	// }

	gs.ActiveSide = CharToColor(activeSide)
	gs.HalfMove = uint16(halfMove)
	gs.FullMove = byte(fullMove)
	gs.EPSquare = CoordToBoardIndex(enPassantSquare)

	for _, availability := range castleAvailability {
		switch availability {
		case 'K':
			gs.CastleRights |= KingsideWhiteCastle
		case 'Q':
			gs.CastleRights |= QueensideWhiteCastle
		case 'k':
			gs.CastleRights |= KingsideBlackCastle
		case 'q':
			gs.CastleRights |= QueensideBlackCastle
		}
	}
}

func (gs GameState) String() string {
	gameStateString := fmt.Sprintf("\n%s\n", gs.Board)

	gameStateString += "side to move: "
	if gs.ActiveSide == White {
		gameStateString += "white"
	} else if gs.ActiveSide == Black {
		gameStateString += "black"
	}
	gameStateString += "\n"

	gameStateString += "castle availability: "
	gameStateString += gs.CastleRights.String()
	gameStateString += "\n"

	gameStateString += "en passant square: "
	if gs.EPSquare == NoSquare {
		gameStateString += "-"
	} else {
		gameStateString += SquareToCoord(gs.EPSquare)
	}
	gameStateString += "\n\n"

	gameStateString += fmt.Sprintf("half move clock: %d\n", gs.HalfMove)
	gameStateString += fmt.Sprintf("full move clock: %d\n\n", gs.FullMove)
	return gameStateString
}

var pawnDirection = [2]int{8, -8}
