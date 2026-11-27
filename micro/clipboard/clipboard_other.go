//go:build !unix

package clipboard

import "github.com/atotto/clipboard"

func nativeClipboardSupported() bool {
	return !clipboard.Unsupported
}

func readNativeClipboard() (string, error) {
	return clipboard.ReadAll()
}

func writeNativeClipboard(text string) error {
	return clipboard.WriteAll(text)
}