package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/dk872/architecture-lab3/painter"
)

// Parser обробляє вхідні дані та генерує відповідні операції.
type Parser struct {
	currentBgColor   painter.Operation          // Поточний фон
	currentRect      *painter.RectOperation     // Поточний прямокутник
	updateOperation  painter.Operation          // Операція оновлення
	figureOperations []*painter.FigureOperation // Операції фігур
	moveOperations   []painter.Operation        // Операції руху
}

// clearOperations очищає всі зібрані операції
func (p *Parser) clearOperations() {
	p.updateOperation = nil
	p.moveOperations = nil
}

// getAllOperations повертає всі операції, зібрані парсером
func (p *Parser) getAllOperations() []painter.Operation {
	var res []painter.Operation

	// Додавання поточного фону, прямокутника, операцій руху, фігур та оновлення
	if p.currentBgColor != nil {
		res = append(res, p.currentBgColor)
	}
	if p.currentRect != nil {
		res = append(res, p.currentRect)
	}
	if len(p.moveOperations) != 0 {
		res = append(res, p.moveOperations...)
	}
	if len(p.figureOperations) != 0 {
		for _, figure := range p.figureOperations {
			res = append(res, figure)
		}
	}
	if p.updateOperation != nil {
		res = append(res, p.updateOperation)
	}
	return res
}

// Parse зчитує вхідні дані та обробляє команди, створюючи відповідні операції
func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.clearOperations()

	var res []painter.Operation

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text()) // Розділяє рядок на поля
		if len(fields) == 0 {
			continue
		}

		err := p.parse(fields) // Обробка кожної команди
		if err != nil {
			return nil, err
		}
	}

	res = append(res, p.getAllOperations()...)
	return res, nil
}

// scale змінює користувацькі координати на ті, з якими працює програма
func scale(value float64) float64 {
	const canvasSize = 800
	return value * canvasSize
}

// parse обробляє окремі команди з вхідних даних
func (p *Parser) parse(fields []string) error {
	command := fields[0]

	switch command {
	case "white":
		p.currentBgColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.currentBgColor = painter.OperationFunc(painter.GreenFill)
	case "update":
		p.updateOperation = painter.UpdateOp
	case "bgrect":
		if len(fields) != 5 {
			return fmt.Errorf("invalid bgrect command format")
		}

		X1, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return fmt.Errorf("invalid X1 value: %v", err)
		}
		Y1, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return fmt.Errorf("invalid Y1 value: %v", err)
		}
		X2, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return fmt.Errorf("invalid X2 value: %v", err)
		}
		Y2, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return fmt.Errorf("invalid Y2 value: %v", err)
		}

		p.currentRect = &painter.RectOperation{X1: scale(X1), Y1: scale(Y1), X2: scale(X2), Y2: scale(Y2)}
	case "figure":
		if len(fields) != 3 {
			return fmt.Errorf("invalid figure command format")
		}

		X, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return fmt.Errorf("invalid figure X value: %v", err)
		}
		Y, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return fmt.Errorf("invalid figure Y value: %v", err)
		}

		fig := &painter.FigureOperation{X: scale(X), Y: scale(Y)}
		p.figureOperations = append(p.figureOperations, fig)
	case "move":
		if len(fields) != 3 {
			return fmt.Errorf("invalid move command format")
		}

		X, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return fmt.Errorf("invalid move X value: %v", err)
		}
		Y, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return fmt.Errorf("invalid move Y value: %v", err)
		}

		moveOp := &painter.MoveFiguresOperation{X: scale(X), Y: scale(Y), Figures: &p.figureOperations}
		p.moveOperations = append(p.moveOperations, moveOp)
	case "reset":
		p.resetState()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

// resetState скидає всі зібрані операції та налаштовує початковий стан
func (p *Parser) resetState() {
	p.currentBgColor = painter.OperationFunc(painter.ResetOperation)
	p.currentRect = nil
	p.updateOperation = nil
	p.figureOperations = nil
	p.moveOperations = nil
}
