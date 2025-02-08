package engine

import "testing"

func TestMove(t *testing.T) {
	type MoveTestCase struct {
		from      Square
		to        Square
		pieceType PieceType
		moveType  MoveType
	}

	t.Run("initializes a move and gets the right pieces", func(t *testing.T) {
		testCases := []MoveTestCase{
			{
				E2,
				E4,
				Pawn,
				Quiet,
			},
			{
				H8,
				H1,
				Knight,
				Quiet,
			},
		}

		for _, testCase := range testCases {
			move := NewMove(testCase.from, testCase.to, testCase.pieceType, testCase.moveType)

			testExpectation(move.FromSquare(), testCase.from, t)
			testExpectation(move.ToSquare(), testCase.to, t)
			testExpectation(move.PieceType(), testCase.pieceType, t)
		}
	})
}

func testExpectation(actual, expected interface{}, t *testing.T) {
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
