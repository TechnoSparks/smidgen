package clipboard

import (
	"os"
	"runtime"

	"github.com/atotto/clipboard"
)

var internalClipboard string
var useInternalClipboard bool

func ClipboardReadAll() (string, error) {
	if !useInternalClipboard {
		if clip, err := clipboard.ReadAll(); err == nil {
			return clip, nil
		}
	}
	return internalClipboard, nil
}

func ClipboardWriteAll(s string) error {
	if !useInternalClipboard && clipboard.WriteAll(s) == nil {
		return nil
	}
	internalClipboard = s
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
