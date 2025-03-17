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

	PreviousStates []IrreversibleState

	moveGen MoveGenerator
}

type IrreversibleState struct {
	CastleRights CastleAvailability
	EPSquare     Square
	HalfMove     uint16
	Moved        Piece
	Destination  Piece
}

func InitializeGameState(board Board, moveGen MoveGenerator) GameState {
	return GameState{
		FullMove: 1,
		Board:    board,
		moveGen:  moveGen,
	}
}

func (gs *GameState) GetMovesForPosition() []Move {
	gs.moveGen.GenerateMoves(gs.ActiveSide, gs.EPSquare, gs.CastleRights)
	return gs.moveGen.GetMoves()
}

func (gs *GameState) ApplyMove(move Move) {
	previous := IrreversibleState{
		CastleRights: gs.CastleRights,
		EPSquare:     gs.EPSquare,
		HalfMove:     gs.HalfMove,
		Moved:        gs.Board.GetPieceAtSquare(move.FromSquare()),
		Destination:  gs.Board.GetPieceAtSquare(move.ToSquare()),
	}

	gs.HalfMove++
	gs.EPSquare = NoSquare

	movingPiece := Piece{
		Color:     gs.ActiveSide,
		PieceType: move.PieceType(),
	}

	moveType := move.MoveType()
	from := move.FromSquare()
	to := move.ToSquare()

	switch moveType {
	case Quiet, Capture:
		gs.Board.MovePiece(movingPiece, from, to)
	case Promotion, CapturePromotion:
		promotionPiece := Piece{
			Color:     gs.ActiveSide,
			PieceType: move.PromotionPieceType(),
		}
		gs.Board.MovePiece(promotionPiece, from, to)
	case CastleKingside, CastleQueenside:
		gs.Board.CastleMove(from, to)
	case EnPassant:
		pawnDirection := gs.getPawnDirection()

		capturedPawn := Square(int(to) - pawnDirection)
		previous.Destination = gs.Board.GetPieceAtSquare(capturedPawn)

		gs.Board.MovePiece(movingPiece, from, to)
		gs.Board.RemovePieceFromSquare(capturedPawn)
	}

	// The halfmove clock gets reset if the move was a capture or if the moved piece was a pawn
	if IsAttackMove(moveType) && previous.Moved.PieceType == Pawn {
		gs.HalfMove = 0
	}

	if previous.Moved.PieceType == Pawn && isDoublePawnPush(from, to) {
		pawnDirection := gs.getPawnDirection()
		epSquare := Square(int(from) - pawnDirection)

		if gs.Board.SquareIsUnderAttackByPawn(epSquare, gs.ActiveSide) {
			gs.EPSquare = epSquare
		}
	}

	gs.PreviousStates = append(gs.PreviousStates, previous)

	// The fullmove number is incremented only after the black side has moved
	if gs.ActiveSide == Black {
		gs.FullMove++
	}

	gs.ActiveSide = gs.ActiveSide.EnemyColor()
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

	if gs.Board != nil {
		gs.Board.SetBoardFromFEN(pieces)
	} else {
		gs.Board = NewBitboardBoard(pieces)
	}

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

func (gs GameState) getPawnDirection() int {
	pawnDirection := North
	if gs.ActiveSide == Black {
		pawnDirection = South
	}

	return pawnDirection
}

// isDoublePawnPush checks to see if the move was a pawn move that went two spaces up
// If the difference in from and to is 16, that means the pawn moved two square forward
func isDoublePawnPush(from, to Square) bool {
	fmt.Println("checking double pawn push")
	fmt.Printf("from: %d  to: %d\n", from, to)
	diff := from - to

	if diff < 0 {
		diff = -diff
	}

	return diff == 16
}
