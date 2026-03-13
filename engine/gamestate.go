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

	movingPiece := Piece{
		Color:     gs.ActiveSide,
		PieceType: move.PieceType(),
	}

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
		gs.Board.pieces[movingPiece.Color][movingPiece.PieceType] ^= mask
		gs.Board.colorBB[movingPiece.Color] ^= mask
		gs.Board.occupancy ^= mask
		gs.Board.squares[to] = movingPiece
		gs.Board.squares[from] = NoPiece
		if movingPiece.PieceType == King {
			gs.Board.kingSq[movingPiece.Color] = to
		}

	case Capture:
		previous.Destination = gs.Board.squares[to]
		captured := previous.Destination
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]
		gs.Board.pieces[movingPiece.Color][movingPiece.PieceType] ^= fromBB | toBB
		gs.Board.colorBB[movingPiece.Color] ^= fromBB | toBB
		gs.Board.pieces[captured.Color][captured.PieceType] &^= toBB
		gs.Board.colorBB[captured.Color] &^= toBB
		gs.Board.occupancy &^= fromBB // to stays occupied, now has moving piece
		gs.Board.squares[to] = movingPiece
		gs.Board.squares[from] = NoPiece
		if movingPiece.PieceType == King {
			gs.Board.kingSq[movingPiece.Color] = to
		}

	case DoublePush:
		mask := SquareMasks[from] | SquareMasks[to]
		gs.Board.pieces[movingPiece.Color][movingPiece.PieceType] ^= mask
		gs.Board.colorBB[movingPiece.Color] ^= mask
		gs.Board.occupancy ^= mask
		gs.Board.squares[to] = movingPiece
		gs.Board.squares[from] = NoPiece
		gs.EPSquare = Square(int(from) + pawnDirection[gs.ActiveSide])

	case CapturePromotion:
		previous.Destination = gs.Board.squares[to]
		promotionPiece := Piece{
			Color:     gs.ActiveSide,
			PieceType: move.PromotionPieceType(),
		}
		gs.Board.MovePiece(promotionPiece, from, to)

	case Promotion:
		promotionPiece := Piece{
			Color:     gs.ActiveSide,
			PieceType: move.PromotionPieceType(),
		}
		gs.Board.MovePiece(promotionPiece, from, to)
	case CastleKingside, CastleQueenside:
		gs.Board.CastleMove(from, to)
	case EnPassant:
		pawnDirection := pawnDirection[gs.ActiveSide]

		capturedPawn := Square(int(to) - pawnDirection)
		previous.Destination = gs.Board.squares[capturedPawn]

		gs.Board.MovePiece(movingPiece, from, to)
		gs.Board.RemovePieceFromSquare(capturedPawn)
	}

	// Update castle rights
	if movingPiece.PieceType == King || movingPiece.PieceType == Rook || previous.Destination.PieceType == Rook {
		gs.UpdateCastleRights(movingPiece, previous.Destination, move)
	}

	// The halfmove clock gets reset if the move was a capture or if the moved piece was a pawn
	if IsAttackMove(moveType) || movingPiece.PieceType == Pawn {
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

	movingPiece := Piece{
		Color:     gs.ActiveSide,
		PieceType: move.PieceType(),
	}
	capturedPiece := previous.Destination

	switch moveType {
	case Quiet:
		mask := SquareMasks[from] | SquareMasks[to]
		gs.Board.pieces[movingPiece.Color][movingPiece.PieceType] ^= mask
		gs.Board.colorBB[movingPiece.Color] ^= mask
		gs.Board.occupancy ^= mask
		gs.Board.squares[from] = movingPiece
		gs.Board.squares[to] = NoPiece
		if movingPiece.PieceType == King {
			gs.Board.kingSq[movingPiece.Color] = from
		}

	case Capture:
		fromBB := SquareMasks[from]
		toBB := SquareMasks[to]
		// Move piece back
		gs.Board.pieces[movingPiece.Color][movingPiece.PieceType] ^= fromBB | toBB
		gs.Board.colorBB[movingPiece.Color] ^= fromBB | toBB
		gs.Board.occupancy ^= fromBB // from now occupied, to stays occupied
		gs.Board.squares[from] = movingPiece
		// Restore captured piece
		gs.Board.pieces[capturedPiece.Color][capturedPiece.PieceType] |= toBB
		gs.Board.colorBB[capturedPiece.Color] |= toBB
		gs.Board.squares[to] = capturedPiece
		if movingPiece.PieceType == King {
			gs.Board.kingSq[movingPiece.Color] = from
		}

	case DoublePush:
		gs.Board.RemovePieceFromSquare(to)
		gs.Board.SetPieceAtPosition(movingPiece, from)

	case EnPassant:
		direction := pawnDirection[gs.ActiveSide]
		capturedSquare := Square(int(to) - direction)

		gs.Board.RemovePieceFromSquare(to)
		gs.Board.SetPieceAtPosition(movingPiece, from)
		gs.Board.SetPieceAtPosition(previous.Destination, capturedSquare)

	case Promotion:
		gs.Board.RemovePieceFromSquare(to)
		gs.Board.SetPieceAtPosition(movingPiece, from)

	case CapturePromotion:
		gs.Board.RemovePieceFromSquare(to)
		gs.Board.SetPieceAtPosition(capturedPiece, to)
		gs.Board.SetPieceAtPosition(movingPiece, from)

	case CastleKingside, CastleQueenside:
		gs.Board.ReverseCastleMove(from, to)
	}
}

func (gs *GameState) UpdateCastleRights(moved Piece, captured Piece, move Move) {
	if moved.PieceType == King {
		gs.CastleRights.RemoveAllRights(gs.ActiveSide)
	} else if moved.PieceType == Rook {
		gs.UpdateRookRights(gs.ActiveSide, move.FromSquare())
	} else if captured.PieceType == Rook {
		gs.UpdateRookRights(gs.ActiveSide.EnemyColor(), move.ToSquare())
	}
}

func (gs *GameState) UpdateRookRights(color Color, square Square) {
	kingside, queenside := rookSquares(color)

	if square == kingside {
		gs.CastleRights.Remove(color, "kingside")
	} else if square == queenside {
		gs.CastleRights.Remove(color, "queenside")
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
