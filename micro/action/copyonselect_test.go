package action

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/smidgen/micro/buffer"
	"github.com/sedwards2009/smidgen/micro/clipboard"
)

// useInternalClipboard switches to an empty internal clipboard and restores the settings afterwards.
func useInternalClipboard(t *testing.T, copyOnSelect bool) {
	t.Helper()
	previousMode := clipboard.Mode()
	previousCopyOnSelect := clipboard.CopyOnSelect()
	t.Cleanup(func() {
		clipboard.SetMode(previousMode)
		clipboard.SetCopyOnSelect(previousCopyOnSelect)
	})
	// Switching modes leaves the internal clipboard empty.
	clipboard.SetMode(clipboard.ModeTerminal)
	clipboard.SetMode(clipboard.ModeInternal)
	clipboard.SetCopyOnSelect(copyOnSelect)
}

func readClipboard(t *testing.T) string {
	t.Helper()
	clip, err := clipboard.ClipboardReadAll()
	if err != nil {
		t.Fatalf("ClipboardReadAll error: %v", err)
	}
	return clip
}

func TestMouseReleaseCopiesSelectionOnlyWhenEnabled(t *testing.T) {
	copyOnSelectValues := []bool{true, false}
	for _, copyOnSelect := range copyOnSelectValues {
		useInternalClipboard(t, copyOnSelect)
		h := newTestBufPane(t, "hello world")
		h.Cursor.SetSelectionStart(buffer.Loc{X: 0, Y: 0})
		h.Cursor.SetSelectionEnd(buffer.Loc{X: 5, Y: 0})

		h.MouseRelease(tcell.NewEventMouse(0, 0, tcell.ButtonNone, tcell.ModNone), func() {})

		want := ""
		if copyOnSelect {
			want = "hello"
		}
		if got := readClipboard(t); got != want {
			t.Errorf("copy on select %v: clipboard = %q, want %q", copyOnSelect, got, want)
		}
	}
}

func TestDoubleClickCopiesWordOnlyWhenEnabled(t *testing.T) {
	copyOnSelectValues := []bool{true, false}
	for _, copyOnSelect := range copyOnSelectValues {
		useInternalClipboard(t, copyOnSelect)
		h := newTestBufPane(t, "hello world")
		h.Buf.Settings["ruler"] = false
		click := tcell.NewEventMouse(1, 0, tcell.Button1, tcell.ModNone)

		h.MousePress(click, func() {})
		h.MousePress(click, func() {})

		if !h.Cursor.HasSelection() {
			t.Fatalf("copy on select %v: double click did not select a word", copyOnSelect)
		}
		want := ""
		if copyOnSelect {
			want = "hello"
		}
		if got := readClipboard(t); got != want {
			t.Errorf("copy on select %v: clipboard = %q, want %q", copyOnSelect, got, want)
		}
	}
}

func TestCopyIgnoresCopyOnSelect(t *testing.T) {
	useInternalClipboard(t, false)
	h := newTestBufPane(t, "hello world")
	h.Cursor.SetSelectionStart(buffer.Loc{X: 0, Y: 0})
	h.Cursor.SetSelectionEnd(buffer.Loc{X: 5, Y: 0})

	h.Copy()

	if got := readClipboard(t); got != "hello" {
		t.Errorf("clipboard after Copy = %q, want %q", got, "hello")
	}
}
