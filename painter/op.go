package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує текстуру у білий колір. Може бути використана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує текстуру у зелений колір. Може бути використана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

// RectOperation визначає координати прямокутника та малює його
type RectOperation struct {
	X1, Y1, X2, Y2 float64
}

func (op RectOperation) Do(t screen.Texture) bool {
	rect := image.Rect(int(op.X1), int(op.Y1), int(op.X2), int(op.Y2))
	t.Fill(rect, color.Black, screen.Src)
	return false
}

// FigureOperation визначає координати центру фігури та виконує малювання
type FigureOperation struct {
	X, Y float64
}

func (op FigureOperation) Do(t screen.Texture) bool {
	// Розміри фігури
	tWidth, tHeight := 400, 300
	centerX, centerY := op.X, op.Y

	// Горизонтальний прямокутник
	horRect := image.Rect(int(centerX)-tWidth/2, int(centerY), int(centerX)+tWidth/2, int(centerY)-tHeight/2)
	// Вертикальний прямокутник
	verRect := image.Rect(int(centerX)-tWidth/5, int(centerY)+tHeight/2, int(centerX)+tWidth/5, int(centerY))

	// Заповнення прямокутників жовтим кольором
	t.Fill(horRect, color.RGBA{R: 255, G: 255, B: 0, A: 255}, screen.Src)
	t.Fill(verRect, color.RGBA{R: 255, G: 255, B: 0, A: 255}, screen.Src)

	return false
}

// MoveFiguresOperation переміщує всі фігури в нові координати
type MoveFiguresOperation struct {
	X, Y    float64
	Figures *[]*FigureOperation
}

func (op MoveFiguresOperation) Do(t screen.Texture) bool {
	for i := range *op.Figures {
		(*op.Figures)[i].X = op.X
		(*op.Figures)[i].Y = op.Y
	}
	return false
}

// ResetOperation очищає текстуру і зафарбовує її чорним
func ResetOperation(t screen.Texture) {
	t.Fill(t.Bounds(), color.Black, screen.Src)
}
