package board

import (
	"nsdb-go-edition/engine"
	"testing"
)

func TestBoardFromFEN(t *testing.T) {
	type BoardFromFenTestCase struct {
		position string
		expected engine.Piece
	}

	board := BitboardBoard{}

	t.Run("initializes board from starting fen string", func(t *testing.T) {
		board.SetBoardFromFEN(InitialStateFenString)

		testCases := []BoardFromFenTestCase{
			{
				"a1",
				engine.Piece{
					PieceType: engine.Rook,
					Color:     engine.White,
				},
			},
			{
				"e1",
				engine.Piece{
					PieceType: engine.King,
					Color:     engine.White,
				},
			},
			{
				"e8",
				engine.Piece{
					PieceType: engine.King,
					Color:     engine.Black,
				},
			},
			{
				"e4",
				engine.Piece{
					PieceType: engine.None,
					Color:     engine.Blank,
				},
			},
		}

		for _, testCase := range testCases {

			index := CoordToBoardIndex(testCase.position)
			actual := board.squares[index]
			if actual != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, actual)
			}

			piece := testCase.expected
			// I only want to check the bitboards if the piece is a real piece
			if piece.PieceType != engine.None {
				bit := board.pieces[piece.Color][piece.PieceType] & (1 << index)

				if bit == 0 {
					t.Errorf("expected bit at index %d to be set, but it was not", index)
				}

			}
		}
	})

	t.Run("correctly sets board from non-starting fen string", func(t *testing.T) {

	})
}
