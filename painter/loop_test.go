package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(tx screen.Texture) {
		t.Log(msg)
		op(tx)
	}
}

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr
	l.stop = make(chan struct{})

	var testOps []string

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.Post(logOp(t, "do green fill", GreenFill))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", GreenFill))
	}

	l.Post(OperationFunc(func(screen.Texture) {
		testOps = append(testOps, "op 1")
		l.Post(OperationFunc(func(screen.Texture) {
			testOps = append(testOps, "op 2")
		}))
	}))
	l.Post(OperationFunc(func(screen.Texture) {
		testOps = append(testOps, "op 3")
	}))

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	if mt.Colors[0] != color.White {
		t.Error("First color is not white:", mt.Colors)
	}
	if len(mt.Colors) != 2 {
		t.Error("Unexpected size of colors:", mt.Colors)
	}

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 3", "op 2"}) {
		t.Error("Bad order:", testOps)
	}
}

func TestLoop_PostEmptyOperation(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr
	l.stop = make(chan struct{})

	l.Start(mockScreen{})

	l.Post(nil)

	l.StopAndWait()

	if tr.lastTexture != nil {
		t.Error("Texture should not have been updated")
	}
}

func TestMessageQueue_PushAndPullOrder(t *testing.T) {
	var mq messageQueue
	var result []string

	mq.push(OperationFunc(func(screen.Texture) { result = append(result, "1") }))
	mq.push(OperationFunc(func(screen.Texture) { result = append(result, "2") }))
	mq.pushFront(OperationFunc(func(screen.Texture) { result = append(result, "0") }))

	for i := 0; i < 3; i++ {
		op := mq.pull()
		if op == nil {
			t.Fatalf("Expected operation %d, got nil", i)
		}
		op.Do(nil)
	}

	expected := []string{"0", "1", "2"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected execution order: got %v, want %v", result, expected)
	}
}

func TestMessageQueue_EmptyBlocksUntilPush(t *testing.T) {
	var mq messageQueue
	done := make(chan struct{})

	go func() {
		op := mq.pull()
		if op == nil {
			t.Error("Expected operation, got nil")
		}
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	mq.push(OperationFunc(func(screen.Texture) {}))

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("Pull did not unblock after push")
	}
}

func TestMessageQueue_PushFront(t *testing.T) {
	var mq messageQueue
	var result []string

	mq.push(OperationFunc(func(screen.Texture) { result = append(result, "tail") }))
	mq.pushFront(OperationFunc(func(screen.Texture) { result = append(result, "head") }))

	for i := 0; i < 2; i++ {
		op := mq.pull()
		op.Do(nil)
	}

	expected := []string{"head", "tail"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Wrong order after pushFront: got %v, want %v", result, expected)
	}
}

func TestMessageQueue_EmptyReturnsTrue(t *testing.T) {
	var mq messageQueue
	if !mq.empty() {
		t.Error("Expected new queue to be empty")
	}
	mq.push(OperationFunc(func(screen.Texture) {}))
	if mq.empty() {
		t.Error("Expected non-empty queue after push")
	}
}

func TestMessageQueue_EmptyAfterStop(t *testing.T) {
	var l Loop
	l.stop = make(chan struct{})

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.StopAndWait()

	if !l.mq.empty() {
		t.Error("Expected queue to be empty after StopAndWait")
	}
}
