package clipboard

import (
	"os"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

var internalClipboard string
var useInternalClipboard bool
var terminalScreen tcell.Screen

// SetScreen sets the screen used to copy to the terminal clipboard (OSC 52)
// when the native clipboard is not available, e.g. over SSH.
func SetScreen(s tcell.Screen) {
	terminalScreen = s
}

func ClipboardReadAll() (string, error) {
	if !useInternalClipboard {
		if clip, err := clipboard.ReadAll(); err == nil {
			return clip, nil
		}
	}
	return internalClipboard, nil
}

func ClipboardWriteAll(s string) error {
	internalClipboard = s
	if !useInternalClipboard && clipboard.WriteAll(s) == nil {
		return nil
	}
	if terminalScreen != nil {
		terminalScreen.SetClipboard([]byte(s))
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
