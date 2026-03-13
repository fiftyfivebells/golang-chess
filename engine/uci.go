package engine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type UCI struct {
	state GameState
}

func (u UCI) uciResponse() {
	fmt.Printf("\nid name %s\n", EngineName)
	fmt.Printf("id author %s\n", EngineAuthor)

	// TODO: add in explanation of different options

	fmt.Printf("uciok\n\n")
}

func (u UCI) newgameResponse() {
	u.state.ClearGameState()
}

func (u *UCI) positionResponse(command string) {
	args := strings.TrimPrefix(command, "position ")

	var fen string
	if strings.HasPrefix(args, "startpos") {
		args = strings.TrimPrefix(args, "startpos ")
		fen = InitialStateFenString
	} else if strings.HasPrefix(args, "fen") {
		args = strings.TrimPrefix(args, "fen ")
		fenArgs := strings.Fields(args)
		fen = strings.Join(fenArgs[0:6], " ")
		args = strings.Join(fenArgs[6:], " ")
	}

	u.state = InitializeGameState(fen)
}

func (u UCI) perftResponse(command string) {
	args := strings.TrimPrefix(command, "go perft ")
	depth, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil {
		fmt.Printf("invalid depth: %s\n", args)
		return
	}

	perftState := PerftState{}
	perftState.gameState = &u.state

	PerftDivide(&perftState, depth)
}

func (u UCI) goResponse() {
	legalMoves := u.state.GetPseudoLegalMovesForPosition()

	if len(legalMoves) == 0 {
		return
	}

	// TODO: implement real searching, instead of just choosing a random move from all possible moves
	fmt.Printf("bestmove %s\n", legalMoves[rand.Intn(len(legalMoves))])
}

func (u UCI) Loop() {

	reader := bufio.NewReader(os.Stdin)

	for {
		command, _ := reader.ReadString('\n')
		command = strings.Replace(command, "\r\n", "\n", -1)

		if commandContains(command, "uci") {
			u.uciResponse()
		} else if commandContains(command, "isready") {
			fmt.Println("readyok")
		} else if commandContains(command, "ucinewgame") {
			u.newgameResponse()
		} else if commandContains(command, "position") {
			u.positionResponse(command)
		} else if commandContains(command, "go perft") {
			u.perftResponse(command)
		} else if commandContains(command, "go") {
			u.goResponse()
		} else if commandContains(command, "printposition") {
			fmt.Println(u.state)
		} else if commandContains(command, "quit") {
			break
		}
	}
}

func commandContains(command, expected string) bool {
	return strings.HasPrefix(command, expected)
}
