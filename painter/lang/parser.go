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
	if p.updateOperation != nil {
		p.updateOperation = nil
	}
	if p.moveOperations != nil {
		p.moveOperations = nil
	}
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

		X1, _ := strconv.ParseFloat(fields[1], 64)
		Y1, _ := strconv.ParseFloat(fields[2], 64)
		X2, _ := strconv.ParseFloat(fields[3], 64)
		Y2, _ := strconv.ParseFloat(fields[4], 64)

		p.currentRect = &painter.RectOperation{X1: X1, Y1: Y1, X2: X2, Y2: Y2}
	case "figure":
		if len(fields) != 3 {
			return fmt.Errorf("invalid figure command format")
		}

		X, _ := strconv.ParseFloat(fields[1], 64)
		Y, _ := strconv.ParseFloat(fields[2], 64)

		fig := &painter.FigureOperation{X: X, Y: Y}
		p.figureOperations = append(p.figureOperations, fig)
	case "move":
		if len(fields) != 3 {
			return fmt.Errorf("invalid move command format")
		}

		X, _ := strconv.ParseFloat(fields[1], 64)
		Y, _ := strconv.ParseFloat(fields[2], 64)

		moveOp := &painter.MoveFiguresOperation{X: X, Y: Y, Figures: &p.figureOperations}
		p.moveOperations = append(p.moveOperations, moveOp)
	case "reset":
		p.resetState()
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
