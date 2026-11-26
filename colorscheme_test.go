package smidgen

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/gdamore/tcell/v2"
)

const userColorschemeText = `color-link default "#112233,#445566"
color-link comment "#778899"
`

// setUpUserColorschemes creates a user colorscheme directory with the given files.
func setUpUserColorschemes(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	SetUserColorschemeDir(dir)
	t.Cleanup(func() {
		SetUserColorschemeDir("")
	})
	return dir
}

func defaultBackground(colorscheme Colorscheme) tcell.Color {
	_, bg, _ := colorscheme.GetColor("default").Decompose()
	return bg
}

func TestListColorschemesIncludesUserColorschemes(t *testing.T) {
	builtIn := ListColorschemes()
	dir := setUpUserColorschemes(t, map[string]string{
		"mine.micro":    userColorschemeText,
		"monokai.micro": userColorschemeText,
		"notes.txt":     "not a colorscheme",
	})
	if err := os.Mkdir(filepath.Join(dir, "folder.micro"), 0755); err != nil {
		t.Fatal(err)
	}

	got := ListColorschemes()

	want := append(slices.Clone(builtIn), "mine")
	if !slices.Equal(got, want) {
		t.Errorf("ListColorschemes() = %v, want %v", got, want)
	}
}

func TestListColorschemesWithoutUserDir(t *testing.T) {
	builtIn := ListColorschemes()
	SetUserColorschemeDir(filepath.Join(t.TempDir(), "missing"))
	t.Cleanup(func() {
		SetUserColorschemeDir("")
	})

	if got := ListColorschemes(); !slices.Equal(got, builtIn) {
		t.Errorf("ListColorschemes() with a missing directory = %v, want the built-in list", got)
	}
}

func TestLoadColorscheme(t *testing.T) {
	builtInMonokai, _ := LoadInternalColorscheme("monokai")
	setUpUserColorschemes(t, map[string]string{
		"mine.micro":    userColorschemeText,
		"monokai.micro": userColorschemeText,
	})
	userBackground := tcell.GetColor("#445566")

	tests := []struct {
		name           string
		wantOK         bool
		wantBackground tcell.Color
	}{
		{name: "mine", wantOK: true, wantBackground: userBackground},
		{name: "monokai", wantOK: true, wantBackground: userBackground},
		{name: "default", wantOK: true},
		{name: "no-such-scheme", wantOK: false},
		{name: "../mine", wantOK: false},
		{name: "", wantOK: false},
	}
	for _, tt := range tests {
		colorscheme, ok := LoadColorscheme(tt.name)
		if ok != tt.wantOK {
			t.Errorf("LoadColorscheme(%q) ok = %v, want %v", tt.name, ok, tt.wantOK)
			continue
		}
		if tt.wantBackground != 0 && defaultBackground(colorscheme) != tt.wantBackground {
			t.Errorf("LoadColorscheme(%q) background = %v, want %v", tt.name, defaultBackground(colorscheme), tt.wantBackground)
		}
	}
	if defaultBackground(builtInMonokai) == userBackground {
		t.Fatal("test colours must differ from built-in monokai")
	}
}

func TestParseColorschemeInvalidLineDoesNotWriteToStdout(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	savedStdout := os.Stdout
	os.Stdout = writer
	log.SetOutput(io.Discard)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
	})
	ParseColorscheme("this is not a color-link line\n")
	os.Stdout = savedStdout
	writer.Close()

	buf := make([]byte, 1)
	if n, _ := reader.Read(buf); n > 0 {
		t.Error("ParseColorscheme wrote to stdout, which would draw over the editor")
	}
}
