package engine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
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

func (u UCI) goResponse() {
	moves := u.state.GetMovesForPosition()

	// TODO: implement real searching, instead of just choosing a random move from all possible moves
	moveIndex := rand.Intn(len(moves))
	bestMove := moves[moveIndex]

	fmt.Printf("bestmove %s\n", bestMove)
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
