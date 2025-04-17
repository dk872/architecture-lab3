package lang

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dk872/architecture-lab3/painter"
)

func TestParser_ParseMultipleCommands(t *testing.T) {
	input := `white
bgrect 0.1 0.2 0.3 0.4
figure 0.5 0.5
move 0.1 0.1
update`

	parser := &Parser{}
	ops, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(ops) != 5 {
		t.Fatalf("Expected 5 operations, but got %d", len(ops))
	}

	t.Run("check types and values", func(t *testing.T) {
		if _, ok := ops[0].(painter.OperationFunc); !ok {
			t.Errorf("Expected OperationFunc (white), got %T", ops[0])
		}

		if rect, ok := ops[1].(*painter.RectOperation); ok {
			want := &painter.RectOperation{
				X1: scale(0.1),
				Y1: scale(0.2),
				X2: scale(0.3),
				Y2: scale(0.4),
			}
			if rect.X1 != want.X1 || rect.Y1 != want.Y1 || rect.X2 != want.X2 || rect.Y2 != want.Y2 {
				t.Errorf("RectOperation has incorrect coordinates: %+v", rect)
			}
		} else {
			t.Errorf("Expected RectOperation, got %T", ops[1])
		}

		if move, ok := ops[2].(*painter.MoveFiguresOperation); ok {
			if move.X != scale(0.1) || move.Y != scale(0.1) {
				t.Errorf("MoveFiguresOperation has incorrect coordinates: %+v", move)
			}
			if move.Figures == nil || len(*move.Figures) != 1 {
				t.Errorf("MoveFiguresOperation has invalid figures list")
			}
		} else {
			t.Errorf("Expected MoveFiguresOperation, got %T", ops[2])
		}

		if fig, ok := ops[3].(*painter.FigureOperation); ok {
			if fig.X != scale(0.5) || fig.Y != scale(0.5) {
				t.Errorf("FigureOperation has incorrect coordinates: %+v", fig)
			}
		} else {
			t.Errorf("Expected FigureOperation, got %T", ops[3])
		}

		if reflect.TypeOf(ops[4]).String() != "painter.updateOp" {
			t.Errorf("Expected updateOp, got %T", ops[4])
		}
	})
}

func TestParser_ParseOtherCommands(t *testing.T) {
	input := `green
figure 0.2 0.3
update`

	parser := &Parser{}
	ops, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(ops) != 3 {
		t.Fatalf("Expected 3 operations, but got %d", len(ops))
	}
	if _, ok := ops[0].(painter.OperationFunc); !ok {
		t.Errorf("Expected OperationFunc (green), got %T", ops[0])
	}
	if fig, ok := ops[1].(*painter.FigureOperation); ok {
		wantX, wantY := scale(0.2), scale(0.3)
		if fig.X != wantX || fig.Y != wantY {
			t.Errorf("FigureOperation has incorrect coordinates: %+v", fig)
		}
	} else {
		t.Errorf("Expected FigureOperation, got %T", ops[1])
	}
	if reflect.TypeOf(ops[2]).String() != "painter.updateOp" {
		t.Errorf("Expected updateOp, got %T", ops[2])
	}
}

func TestParser_ParseReset(t *testing.T) {
	input := `white
bgrect 0.1 0.1 0.2 0.2
reset
update`

	parser := &Parser{}
	ops, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(ops) != 2 {
		t.Fatalf("Expected 2 operations after reset (reset, update), but got %d", len(ops))
	}
	if _, ok := ops[0].(painter.OperationFunc); !ok {
		t.Errorf("Expected OperationFunc (reset), got %T", ops[0])
	}
	if reflect.TypeOf(ops[1]).String() != "painter.updateOp" {
		t.Errorf("Expected updateOp, got %T", ops[1])
	}
}

func TestParser_UnknownAndInvalidCommands(t *testing.T) {
	parser := &Parser{}

	t.Run("unknown command", func(t *testing.T) {
		_, err := parser.Parse(strings.NewReader("foo"))
		if err == nil || !strings.Contains(err.Error(), "unknown command") {
			t.Errorf("Expected unknown command error, got %v", err)
		}
	})

	t.Run("invalid bgrect format", func(t *testing.T) {
		_, err := parser.Parse(strings.NewReader("bgrect 0.1 0.2"))
		if err == nil || !strings.Contains(err.Error(), "invalid bgrect command format") {
			t.Errorf("Expected invalid bgrect command format error, got %v", err)
		}
	})

	t.Run("invalid bgrect value: not a number", func(t *testing.T) {
		_, err := parser.Parse(strings.NewReader("bgrect 0.6 b 0.8 0.9"))
		if err == nil || !strings.Contains(err.Error(), "invalid") {
			t.Errorf("Expected invalid Y1 value error, got %v", err)
		}
	})

	t.Run("invalid bgrect value: out of range low", func(t *testing.T) {
		_, err := parser.Parse(strings.NewReader("bgrect -0.1 0.2 0.3 0.4"))
		if err == nil || !strings.Contains(err.Error(), "X1 value -0.10 out of range") {
			t.Errorf("Expected out of range X1 error, got %v", err)
		}
	})

	t.Run("invalid bgrect value: out of range high", func(t *testing.T) {
		_, err := parser.Parse(strings.NewReader("bgrect 0.1 0.2 1.2 0.4"))
		if err == nil || !strings.Contains(err.Error(), "X2 value 1.20 out of range") {
			t.Errorf("Expected out of range X2 error, got %v", err)
		}
	})

	t.Run("invalid figure value", func(t *testing.T) {
		_, err := parser.Parse(strings.NewReader("figure 0.5 a"))
		if err == nil || !strings.Contains(err.Error(), "invalid Y value") {
			t.Errorf("Expected invalid figure Y value error, got %v", err)
		}
	})

	t.Run("invalid move value", func(t *testing.T) {
		_, err := parser.Parse(strings.NewReader("move 5 0.5"))
		if err == nil || !strings.Contains(err.Error(), "X value 5.00 out of range") {
			t.Errorf("Expected out of range X error in move, got %v", err)
		}
	})
}
