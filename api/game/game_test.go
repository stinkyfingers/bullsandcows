package game

import "testing"

func TestRandomCode(t *testing.T) {
	code := randomCode()
	if len(code) != 4 {
		t.Errorf("expected 4, got %d", len(code))
	}
}

func TestNewGame(t *testing.T) {
	game := NewGame(nil, "123")
	if len(game.Code) != 4 {
		t.Errorf("expected 4, got %d", len(game.Code))
	}
	if game.Rounds == nil {
		t.Error("expected game, got nil")
	}
}

func TestPlayRound(t *testing.T) {
	g := &Game{
		Code: [4]Peg{Red, Yellow, Blue, Purple},
	}

	tests := []struct {
		guess         Sequence
		expectedBulls int
		expectedCows  int
	}{
		{
			guess:         Sequence{Yellow, Red, Blue, Green},
			expectedBulls: 1,
			expectedCows:  2,
		}, {
			guess:         Sequence{Green, Green, Green, Green},
			expectedBulls: 0,
			expectedCows:  0,
		}, {
			guess:         Sequence{Purple, Red, Yellow, Blue},
			expectedBulls: 0,
			expectedCows:  4,
		}, {
			guess:         Sequence{Red, Yellow, Blue, Purple},
			expectedBulls: 4,
			expectedCows:  0,
		},
	}

	for i, test := range tests {
		err := g.PlayRound(test.guess)
		if err != nil {
			t.Error(err)
		}
		if len(g.Rounds) != i+1 {
			t.Errorf("expected %d rounds, got %d", i+1, len(g.Rounds))
		}
		if g.Rounds[i].Bulls != test.expectedBulls {
			t.Errorf("expected %d bulls, got %d", test.expectedBulls, g.Rounds[i].Bulls)
		}
		if g.Rounds[i].Cows != test.expectedCows {
			t.Errorf("expected %d bulls, got %d", test.expectedCows, g.Rounds[i].Cows)
		}
	}
}
