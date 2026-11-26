package clipboard

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
)

type fakeTty struct {
	tcell.Tty
	written bytes.Buffer
}

func (f *fakeTty) Write(p []byte) (int, error) {
	return f.written.Write(p)
}

type fakeScreen struct {
	tcell.Screen
	tty *fakeTty
}

func (f fakeScreen) Tty() (tcell.Tty, bool) {
	return f.tty, true
}

type fakeNative struct {
	text     string
	readErr  error
	writeErr error
}

// setUp resets the package state and installs a fake system clipboard and terminal.
func setUp(t *testing.T, startMode string, nativeAvailable bool) (*fakeNative, *fakeTty) {
	t.Helper()
	native := &fakeNative{}
	tty := &fakeTty{}

	savedRead, savedWrite := nativeReadAll, nativeWriteAll
	savedMode, savedUseInternal, savedScreen := mode, useInternalClipboard, terminalScreen
	t.Cleanup(func() {
		nativeReadAll, nativeWriteAll = savedRead, savedWrite
		mode, useInternalClipboard, terminalScreen = savedMode, savedUseInternal, savedScreen
		internalClipboard = ""
	})

	nativeReadAll = func() (string, error) {
		if native.readErr != nil {
			return "", native.readErr
		}
		return native.text, nil
	}
	nativeWriteAll = func(s string) error {
		if native.writeErr != nil {
			return native.writeErr
		}
		native.text = s
		return nil
	}
	mode = startMode
	useInternalClipboard = !nativeAvailable
	internalClipboard = ""
	SetScreen(fakeScreen{tty: tty})
	return native, tty
}

func readAll(t *testing.T) string {
	t.Helper()
	clip, err := ClipboardReadAll()
	if err != nil {
		t.Fatalf("ClipboardReadAll error: %v", err)
	}
	return clip
}

func TestTerminalModeSendsOSC52(t *testing.T) {
	native, tty := setUp(t, ModeTerminal, true)

	ClipboardWriteAll("héllo\nworld")

	want := "\x1b]52;c;aMOpbGxvCndvcmxk\a"
	if got := tty.written.String(); got != want {
		t.Errorf("tty got %q, want %q", got, want)
	}
	if native.text != "" {
		t.Errorf("system clipboard = %q, want it untouched", native.text)
	}
	if got := readAll(t); got != "héllo\nworld" {
		t.Errorf("paste = %q, want the last copy", got)
	}
}

func TestInternalModeStaysInside(t *testing.T) {
	native, tty := setUp(t, ModeInternal, true)

	ClipboardWriteAll("abc")

	if tty.written.Len() != 0 {
		t.Errorf("tty got %q, want nothing", tty.written.String())
	}
	if native.text != "" {
		t.Errorf("system clipboard = %q, want it untouched", native.text)
	}
	if got := readAll(t); got != "abc" {
		t.Errorf("paste = %q, want %q", got, "abc")
	}
}

func TestExternalModeUsesSystemClipboard(t *testing.T) {
	native, tty := setUp(t, ModeExternal, true)

	ClipboardWriteAll("abc")
	if native.text != "abc" {
		t.Errorf("system clipboard = %q, want %q", native.text, "abc")
	}
	if tty.written.Len() != 0 {
		t.Errorf("tty got %q, want nothing", tty.written.String())
	}

	native.text = "from another app"
	if got := readAll(t); got != "from another app" {
		t.Errorf("paste = %q, want the system clipboard", got)
	}
}

func TestExternalModeFailedReadDoesNotPasteOldCopy(t *testing.T) {
	native, _ := setUp(t, ModeExternal, true)

	ClipboardWriteAll("old copy")
	// e.g. xclip fails because the clipboard now holds an image.
	native.readErr = errors.New("no text in clipboard")

	clip, err := ClipboardReadAll()
	if err == nil || clip != "" {
		t.Errorf("paste = %q, %v; want no text and the read error", clip, err)
	}
}

func TestExternalModeFailedWriteFallsBackToInternal(t *testing.T) {
	native, _ := setUp(t, ModeExternal, true)
	native.writeErr = errors.New("X server gone")
	native.readErr = errors.New("X server gone")

	ClipboardWriteAll("abc")

	if got := readAll(t); got != "abc" {
		t.Errorf("paste = %q, want %q", got, "abc")
	}
}

func TestExternalModeWithoutDisplayUsesInternal(t *testing.T) {
	native, tty := setUp(t, ModeExternal, false)

	ClipboardWriteAll("abc")

	if native.text != "" || tty.written.Len() != 0 {
		t.Errorf("system clipboard = %q, tty = %q; want both untouched", native.text, tty.written.String())
	}
	if got := readAll(t); got != "abc" {
		t.Errorf("paste = %q, want %q", got, "abc")
	}
}

func TestSetMode(t *testing.T) {
	setUp(t, ModeInternal, true)
	ClipboardWriteAll("abc")

	SetMode("banana")
	if mode != ModeInternal || readAll(t) != "abc" {
		t.Errorf("unknown mode changed state: mode = %q", mode)
	}

	SetMode(ModeInternal)
	if readAll(t) != "abc" {
		t.Error("setting the same mode cleared the internal clipboard")
	}

	SetMode(ModeTerminal)
	if mode != ModeTerminal {
		t.Errorf("mode = %q, want %q", mode, ModeTerminal)
	}
	if got := readAll(t); got != "" {
		t.Errorf("paste after changing mode = %q, want empty", got)
	}
}
