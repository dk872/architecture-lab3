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

// parseCoordinates отримує координати та перевіряє правильність їх введення: мають бути від 0 до 1
func parseCoordinates(raw string, name string) (float64, error) {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value: %v", name, err)
	}
	if value <= 0.0 || value >= 1.0 {
		return 0, fmt.Errorf("%s value %.2f out of range [0.0 - 1.0]", name, value)
	}
	return value, nil
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

		X1, err := parseCoordinates(fields[1], "X1")
		if err != nil {
			return err
		}
		Y1, err := parseCoordinates(fields[2], "Y1")
		if err != nil {
			return err
		}
		X2, err := parseCoordinates(fields[3], "X2")
		if err != nil {
			return err
		}
		Y2, err := parseCoordinates(fields[4], "Y2")
		if err != nil {
			return err
		}

		p.currentRect = &painter.RectOperation{X1: scale(X1), Y1: scale(Y1), X2: scale(X2), Y2: scale(Y2)}
	case "figure":
		if len(fields) != 3 {
			return fmt.Errorf("invalid figure command format")
		}

		X, err := parseCoordinates(fields[1], "X")
		if err != nil {
			return err
		}
		Y, err := parseCoordinates(fields[2], "Y")
		if err != nil {
			return err
		}

		fig := &painter.FigureOperation{X: scale(X), Y: scale(Y)}
		p.figureOperations = append(p.figureOperations, fig)
	case "move":
		if len(fields) != 3 {
			return fmt.Errorf("invalid move command format")
		}

		X, err := parseCoordinates(fields[1], "X")
		if err != nil {
			return err
		}
		Y, err := parseCoordinates(fields[2], "Y")
		if err != nil {
			return err
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
