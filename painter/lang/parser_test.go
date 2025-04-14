package lang

import (
	"strings"
	"testing"
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
}
