//go:build unix

package clipboard

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	unixPasteCmd []string
	unixCopyCmd  []string
	unixTrimDos  bool

	unixDetectionDone bool
)

// findExecutable locates an executable by name using PATH and os.Stat, without
// going through os/exec.LookPath (which relies on faccessat2 and can trigger
// SIGSYS under seccomp filters such as Termux's).
func findExecutable(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	if strings.ContainsRune(name, os.PathSeparator) {
		info, err := os.Stat(name)
		if err != nil || info.IsDir() {
			return "", false
		}
		if runtime.GOOS == "windows" {
			return name, true
		}
		if info.Mode()&0o111 == 0 {
			return "", false
		}
		return name, true
	}

	dirs := filepath.SplitList(os.Getenv("PATH"))
	if len(dirs) == 0 {
		dirs = []string{""}
	}
	for _, dir := range dirs {
		if dir == "" {
			dir = "."
		}
		base := filepath.Join(dir, name)
		candidates := []string{base}
		if runtime.GOOS == "windows" {
			if filepath.Ext(base) == "" {
				for _, ext := range []string{".exe", ".bat", ".cmd", ".com"} {
					candidates = append(candidates, base+ext)
				}
			}
		}
		for _, c := range candidates {
			info, err := os.Stat(c)
			if err != nil || info.IsDir() {
				continue
			}
			if runtime.GOOS != "windows" && info.Mode()&0o111 == 0 {
				continue
			}
			return c, true
		}
	}
	return "", false
}

func detectUnixClipboard() {
	if unixDetectionDone {
		return
	}
	unixDetectionDone = true

	unixPasteCmd = nil
	unixCopyCmd = nil
	unixTrimDos = false

	// Wayland
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if p, ok := findExecutable("wl-paste"); ok {
			if c, ok2 := findExecutable("wl-copy"); ok2 {
				unixPasteCmd = []string{p, "--no-newline"}
				unixCopyCmd = []string{c}
				return
			}
		}
	}

	// xclip
	if p, ok := findExecutable("xclip"); ok {
		unixPasteCmd = []string{p, "-out", "-selection", "clipboard"}
		unixCopyCmd = []string{p, "-in", "-selection", "clipboard"}
		return
	}

	// xsel
	if p, ok := findExecutable("xsel"); ok {
		unixPasteCmd = []string{p, "--output", "--clipboard"}
		unixCopyCmd = []string{p, "--input", "--clipboard"}
		return
	}

	// Termux
	if p, ok := findExecutable("termux-clipboard-get"); ok {
		if c, ok2 := findExecutable("termux-clipboard-set"); ok2 {
			unixPasteCmd = []string{p}
			unixCopyCmd = []string{c}
			return
		}
	}

	// WSL
	if c, ok := findExecutable("clip.exe"); ok {
		if p, ok2 := findExecutable("powershell.exe"); ok2 {
			unixPasteCmd = []string{p, "Get-Clipboard"}
			unixCopyCmd = []string{c}
			unixTrimDos = true
			return
		}
	}
}

func nativeClipboardSupported() bool {
	detectUnixClipboard()
	return unixPasteCmd != nil && unixCopyCmd != nil
}

func readNativeClipboard() (string, error) {
	detectUnixClipboard()
	if unixPasteCmd == nil {
		return "", errors.New("No clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set.")
	}
	cmd := exec.Command(unixPasteCmd[0], unixPasteCmd[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	result := string(out)
	if unixTrimDos && len(result) > 1 {
		result = result[:len(result)-2]
	}
	return result, nil
}

func writeNativeClipboard(text string) error {
	detectUnixClipboard()
	if unixCopyCmd == nil {
		return errors.New("No clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set.")
	}
	cmd := exec.Command(unixCopyCmd[0], unixCopyCmd[1:]...)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		in.Close()
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
