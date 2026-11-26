package action

import (
	"testing"

	"github.com/sedwards2009/smidgen/micro/buffer"
	"github.com/sedwards2009/smidgen/micro/clipboard"
)

func TestPasteEmptyClipboardKeepsSelection(t *testing.T) {
	previousMode := clipboard.Mode()
	t.Cleanup(func() {
		clipboard.SetMode(previousMode)
	})
	// Switching modes leaves the internal clipboard empty.
	clipboard.SetMode(clipboard.ModeTerminal)
	clipboard.SetMode(clipboard.ModeInternal)

	pastes := map[string]func(h *BufPane) bool{
		"Paste":        (*BufPane).Paste,
		"PastePrimary": (*BufPane).PastePrimary,
	}
	for name, paste := range pastes {
		h := newTestBufPane(t, "hello world")
		h.Cursor.SetSelectionStart(buffer.Loc{X: 0, Y: 0})
		h.Cursor.SetSelectionEnd(buffer.Loc{X: 5, Y: 0})

		paste(h)

		if got := string(h.Buf.Bytes()); got != "hello world" {
			t.Errorf("%s with empty clipboard: buffer = %q, want %q", name, got, "hello world")
		}
	}
}

func TestPasteReplacesSelection(t *testing.T) {
	previousMode := clipboard.Mode()
	t.Cleanup(func() {
		clipboard.SetMode(previousMode)
	})
	clipboard.SetMode(clipboard.ModeInternal)
	clipboard.ClipboardWriteAll("bye")

	h := newTestBufPane(t, "hello world")
	h.Cursor.SetSelectionStart(buffer.Loc{X: 0, Y: 0})
	h.Cursor.SetSelectionEnd(buffer.Loc{X: 5, Y: 0})

	h.Paste()

	if got := string(h.Buf.Bytes()); got != "bye world" {
		t.Errorf("buffer = %q, want %q", got, "bye world")
	}
}
