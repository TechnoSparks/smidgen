package action

import (
	"testing"

	"github.com/sedwards2009/smidgen/micro/buffer"
	"github.com/sedwards2009/smidgen/micro/config"
	"github.com/sedwards2009/smidgen/micro/display"
)

func newTestBufPane(t *testing.T, content string) *BufPane {
	t.Helper()
	config.InitRuntimeFiles()
	buf := buffer.NewBufferFromString(content, "")
	win := display.NewBufWindow(0, 0, 80, 24, buf)
	return NewBufPane(buf, win)
}

func TestSearchForwardSelectionAdjustment(t *testing.T) {
	h := newTestBufPane(t, "foo bar foo")
	h.Cursor.Loc = buffer.Loc{X: 3, Y: 0}
	h.Cursor.CurSelection[0] = buffer.Loc{X: 0, Y: 0}
	h.Cursor.CurSelection[1] = buffer.Loc{X: 3, Y: 0}
	h.Cursor.HasSelection() // ensure selection state is populated

	ok, err := h.Search("foo", false, false, true)
	if err != nil || !ok {
		t.Fatalf("forward search failed: ok=%v err=%v", ok, err)
	}
	if h.Cursor.CurSelection[0] != (buffer.Loc{X: 8, Y: 0}) {
		t.Fatalf("forward search did not skip current selection: got %v", h.Cursor.CurSelection[0])
	}
}

func TestSearchBackwardSelectionAdjustment(t *testing.T) {
	h := newTestBufPane(t, "foo bar foo")
	// Start with the second match selected.
	h.Cursor.Loc = buffer.Loc{X: 11, Y: 0}
	h.Cursor.CurSelection[0] = buffer.Loc{X: 8, Y: 0}
	h.Cursor.CurSelection[1] = buffer.Loc{X: 11, Y: 0}

	ok, err := h.Search("foo", false, false, false)
	if err != nil || !ok {
		t.Fatalf("backward search failed: ok=%v err=%v", ok, err)
	}
	if h.Cursor.CurSelection[0] != (buffer.Loc{X: 0, Y: 0}) {
		t.Fatalf("backward search did not move to previous match: got %v", h.Cursor.CurSelection[0])
	}
}

func TestSearchBackwardNoMatch(t *testing.T) {
	h := newTestBufPane(t, "foo bar foo")
	h.Cursor.Loc = buffer.Loc{X: 3, Y: 0}
	h.Cursor.CurSelection[0] = buffer.Loc{X: 0, Y: 0}
	h.Cursor.CurSelection[1] = buffer.Loc{X: 3, Y: 0}

	ok, err := h.Search("zzz", false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected no match")
	}
}