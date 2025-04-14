package lang

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dk872/architecture-lab3/painter"
)

func TestParser(t *testing.T) {
	input := `white
bgrect 0.1 0.2 0.3 0.4
figure 0.5 0.5
move 0.1 0.1
update`

	parser := &Parser{}
	ops, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Помилка: %v", err)
	}

	if len(ops) != 5 {
		t.Fatalf("Очікувалося 5 операцій, але отримано %d", len(ops))
	}

	t.Run("check types and values", func(t *testing.T) {
		// #0: WhiteFill
		if _, ok := ops[0].(painter.OperationFunc); !ok {
			t.Errorf("Очікувалась OperationFunc (white), отримано %T", ops[0])
		}

		// #1: RectOperation
		if rect, ok := ops[1].(*painter.RectOperation); ok {
			want := &painter.RectOperation{X1: 80, Y1: 160, X2: 240, Y2: 320}
			if rect.X1 != want.X1 || rect.Y1 != want.Y1 || rect.X2 != want.X2 || rect.Y2 != want.Y2 {
				t.Errorf("RectOperation має неправильні координати: %+v", rect)
			}
		} else {
			t.Errorf("Очікувався RectOperation, отримано %T", ops[1])
		}

		// #2: MoveFiguresOperation
		if move, ok := ops[2].(*painter.MoveFiguresOperation); ok {
			if move.X != 80 || move.Y != 80 {
				t.Errorf("MoveFiguresOperation має неправильні координати: %+v", move)
			}
			if move.Figures == nil || len(*move.Figures) != 1 {
				t.Errorf("MoveFiguresOperation має некоректний список фігур")
			}
		} else {
			t.Errorf("Очікувався MoveFiguresOperation, отримано %T", ops[2])
		}

		// #3: FigureOperation
		if fig, ok := ops[3].(*painter.FigureOperation); ok {
			if fig.X != 400 || fig.Y != 400 {
				t.Errorf("FigureOperation має неправильні координати: %+v", fig)
			}
		} else {
			t.Errorf("Очікувався FigureOperation, отримано %T", ops[3])
		}

		// #4: Update
		if reflect.TypeOf(ops[4]).String() != "painter.updateOp" {
			t.Errorf("Очікувався updateOp, отримано %T", ops[4])
		}
	})
}