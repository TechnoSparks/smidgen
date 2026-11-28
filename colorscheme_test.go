package smidgen

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
	defer reader.Close()
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

func TestLoadColorschemeBrokenUserFile(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
	})
	setUpUserColorschemes(t, map[string]string{
		"empty.micro":    "",
		"garbage.micro":  "this is not a colorscheme\n\x00\xff random bytes\n",
		"partial.micro":  "nonsense\ncolor-link default \"#112233,#445566\"\nmore nonsense\n",
		"badcolor.micro": "color-link default \"not-a-colour,also-not\"\n",
	})

	// Broken files still load. They give the terminal default, not a crash.
	plainBackground := tcell.ColorDefault
	for _, name := range []string{"empty", "garbage", "badcolor"} {
		colorscheme, ok := LoadColorscheme(name)
		if !ok {
			t.Errorf("LoadColorscheme(%q) ok = false, want true", name)
			continue
		}
		if got := defaultBackground(colorscheme); got != plainBackground {
			t.Errorf("LoadColorscheme(%q) background = %v, want %v", name, got, plainBackground)
		}
	}

	// Valid lines in a file with bad lines still apply.
	colorscheme, ok := LoadColorscheme("partial")
	if !ok {
		t.Fatal("LoadColorscheme(\"partial\") ok = false, want true")
	}
	if got, want := defaultBackground(colorscheme), tcell.GetColor("#445566"); got != want {
		t.Errorf("partial file background = %v, want %v", got, want)
	}
}

func TestLoadColorschemeWithProblems(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
	})
	setUpUserColorschemes(t, map[string]string{
		"valid.micro": "# comment\n\ncolor-link default \"#112233,#445566\"\n",
		"mixed.micro": "# comment\ncolor-link default \"#112233,#445566\"\n\nnonsense\ncolor-link comment \"#ff0000\"\nmore nonsense\r\ncolor-link\n",
		"empty.micro": "",
	})

	tests := []struct {
		name         string
		wantOK       bool
		wantProblems []string
	}{
		{name: "valid", wantOK: true, wantProblems: nil},
		{name: "empty", wantOK: true, wantProblems: nil},
		{name: "mixed", wantOK: true, wantProblems: []string{
			"line 4: not a valid color-link statement",
			"line 6: not a valid color-link statement",
			"line 7: not a valid color-link statement",
		}},
		{name: "monokai", wantOK: true, wantProblems: nil},
		{name: "no-such-scheme", wantOK: false, wantProblems: nil},
	}
	for _, tt := range tests {
		colorscheme, problems, ok := LoadColorschemeWithProblems(tt.name)
		if ok != tt.wantOK {
			t.Errorf("LoadColorschemeWithProblems(%q) ok = %v, want %v", tt.name, ok, tt.wantOK)
		}
		if !slices.Equal(problems, tt.wantProblems) {
			t.Errorf("LoadColorschemeWithProblems(%q) problems = %q, want %q", tt.name, problems, tt.wantProblems)
		}
		if tt.name == "mixed" {
			// The valid lines around the bad ones still apply.
			if got, want := defaultBackground(colorscheme), tcell.GetColor("#445566"); got != want {
				t.Errorf("mixed file background = %v, want %v", got, want)
			}
			if fg, _, _ := colorscheme.GetColor("comment").Decompose(); fg != tcell.GetColor("#ff0000") {
				t.Errorf("mixed file comment colour = %v, want #ff0000", fg)
			}
		}
	}
}

func TestBuiltInColorschemesHaveNoProblems(t *testing.T) {
	for _, name := range ListColorschemes() {
		_, problems, ok := LoadColorschemeWithProblems(name)
		if !ok {
			t.Errorf("built-in colorscheme %q failed to load", name)
		}
		if len(problems) > 0 {
			t.Errorf("built-in colorscheme %q has problems: %q", name, problems)
		}
	}
}

func TestLoadColorschemeWithProblemsUnreadableUserFile(t *testing.T) {
	dir := setUpUserColorschemes(t, map[string]string{})
	// A folder named like a colourscheme cannot be read as a file.
	if err := os.Mkdir(filepath.Join(dir, "broken.micro"), 0755); err != nil {
		t.Fatal(err)
	}

	colorscheme, problems, ok := LoadColorschemeWithProblems("broken")
	if len(problems) != 1 || !strings.HasPrefix(problems[0], "cannot read file: ") {
		t.Errorf("problems = %q, want one \"cannot read file: ...\" problem", problems)
	}
	if ok || colorscheme != nil {
		t.Errorf("ok = %v, colorscheme nil = %v; want no colourscheme, as no built-in has this name", ok, colorscheme == nil)
	}
}

func TestLoadColorschemeWithProblemsUnreadableUserFileUsesBuiltIn(t *testing.T) {
	dir := setUpUserColorschemes(t, map[string]string{})
	if err := os.Mkdir(filepath.Join(dir, "monokai.micro"), 0755); err != nil {
		t.Fatal(err)
	}

	colorscheme, problems, ok := LoadColorschemeWithProblems("monokai")
	if !ok || colorscheme == nil {
		t.Fatal("LoadColorschemeWithProblems(\"monokai\") did not fall back to the built-in colourscheme")
	}
	if len(problems) != 1 {
		t.Errorf("problems = %q, want one problem for the unreadable file", problems)
	}
}

func TestLoadColorschemeWithProblemsMissingUserFileIsSilent(t *testing.T) {
	setUpUserColorschemes(t, map[string]string{})

	_, problems, ok := LoadColorschemeWithProblems("monokai")
	if !ok || len(problems) != 0 {
		t.Errorf("ok = %v, problems = %q; want the built-in colourscheme with no problems", ok, problems)
	}
}
