package clipboard

import (
	"encoding/base64"
	"os"
	"runtime"
	"slices"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

// Clipboard modes, as in micro's "clipboard" option.
const (
	// ModeTerminal copies to the terminal clipboard with OSC 52. Paste reads the internal clipboard.
	ModeTerminal = "terminal"
	// ModeExternal uses the system clipboard (pbcopy, xclip, xsel, wl-copy).
	ModeExternal = "external"
	// ModeInternal keeps copied text inside the program only.
	ModeInternal = "internal"
)

var modes = []string{ModeTerminal, ModeExternal, ModeInternal}

var mode = ModeExternal
var internalClipboard string
var useInternalClipboard bool
var terminalScreen tcell.Screen

// The system clipboard functions. Tests replace them.
var nativeReadAll = clipboard.ReadAll
var nativeWriteAll = clipboard.WriteAll

// SetMode sets the clipboard mode to one of ModeTerminal, ModeExternal or ModeInternal.
// Unknown modes are ignored. Changing the mode clears the internal clipboard.
func SetMode(newMode string) {
	if newMode == mode || !slices.Contains(modes, newMode) {
		return
	}
	mode = newMode
	internalClipboard = ""
}

// SetScreen sets the screen used to copy to the terminal clipboard (OSC 52)
// in terminal mode.
func SetScreen(s tcell.Screen) {
	terminalScreen = s
}

// writeTerminalClipboard sends s to the terminal clipboard with OSC 52.
// It writes to the tty directly because tcell only sends OSC 52 for
// xterm-like TERM values, which leaves out tmux, screen and kitty.
func writeTerminalClipboard(s string) {
	if terminalScreen == nil {
		return
	}
	tty, ok := terminalScreen.Tty()
	if !ok {
		return
	}
	tty.Write([]byte("\x1b]52;c;" + base64.StdEncoding.EncodeToString([]byte(s)) + "\a"))
}

func ClipboardReadAll() (string, error) {
	if mode != ModeExternal || useInternalClipboard {
		return internalClipboard, nil
	}
	clip, err := nativeReadAll()
	// Fall back only to text that the system clipboard did not take.
	if err != nil && internalClipboard != "" {
		return internalClipboard, nil
	}
	return clip, err
}

func ClipboardWriteAll(s string) error {
	switch mode {
	case ModeTerminal:
		internalClipboard = s
		writeTerminalClipboard(s)
	case ModeInternal:
		internalClipboard = s
	case ModeExternal:
		if !useInternalClipboard && nativeWriteAll(s) == nil {
			internalClipboard = ""
			return nil
		}
		internalClipboard = s
	}
	return nil
}

func isUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "openbsd" || runtime.GOOS == "netbsd"
}

func isX11Present() bool {
	return os.Getenv("DISPLAY") != ""
}

func isWaylandPresent() bool {
	return os.Getenv("WAYLAND_DISPLAY") != ""
}

func init() {
	useInternalClipboard = clipboard.Unsupported

	if !clipboard.Unsupported && isUnix() && !(isWaylandPresent() || isX11Present()) {
		useInternalClipboard = true
	}
}
