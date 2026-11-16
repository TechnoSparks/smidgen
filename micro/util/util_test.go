package util

import "testing"

func TestAbs(t *testing.T) {
	cases := []struct{ n, want int }{
		{0, 0}, {5, 5}, {-5, 5}, {-1, 1},
	}
	for _, c := range cases {
		if got := Abs(c.n); got != c.want {
			t.Errorf("Abs(%d) = %d, want %d", c.n, got, c.want)
		}
	}
}

func TestMin(t *testing.T) {
	if Min(1, 2) != 1 {
		t.Error("Min(1,2)")
	}
	if Min(2, 1) != 1 {
		t.Error("Min(2,1)")
	}
	if Min(3, 3) != 3 {
		t.Error("Min(3,3)")
	}
}

func TestMax(t *testing.T) {
	if Max(1, 2) != 2 {
		t.Error("Max(1,2)")
	}
	if Max(2, 1) != 2 {
		t.Error("Max(2,1)")
	}
	if Max(3, 3) != 3 {
		t.Error("Max(3,3)")
	}
}

func TestClamp(t *testing.T) {
	cases := []struct{ val, min, max, want int }{
		{5, 0, 10, 5},
		{-1, 0, 10, 0},
		{11, 0, 10, 10},
		{0, 0, 10, 0},
	}
	for _, c := range cases {
		if got := Clamp(c.val, c.min, c.max); got != c.want {
			t.Errorf("Clamp(%d,%d,%d) = %d, want %d", c.val, c.min, c.max, got, c.want)
		}
	}
}

func TestIsWhitespace(t *testing.T) {
	cases := []struct {
		r    rune
		want bool
	}{
		{' ', true}, {'\t', true}, {'\n', true}, {'\r', true},
		{'a', false}, {'0', false},
	}
	for _, c := range cases {
		if got := IsWhitespace(c.r); got != c.want {
			t.Errorf("IsWhitespace(%q) = %v, want %v", c.r, got, c.want)
		}
	}
}

func TestIsWordChar(t *testing.T) {
	if !IsWordChar('a') {
		t.Error("a should be word char")
	}
	if !IsWordChar('_') {
		t.Error("_ should be word char")
	}
	if IsWordChar(' ') {
		t.Error("space should not be word char")
	}
	if IsWordChar('!') {
		t.Error("! should not be word char")
	}
}

func TestIsAlphanumeric(t *testing.T) {
	if !IsAlphanumeric('a') || !IsAlphanumeric('Z') || !IsAlphanumeric('0') {
		t.Error("alphanumeric check failed")
	}
	if IsAlphanumeric(' ') || IsAlphanumeric('!') {
		t.Error("non-alphanumeric accepted")
	}
}

func TestIntOpt(t *testing.T) {
	if IntOpt(float64(5)) != 5 {
		t.Error("IntOpt(float64)")
	}
	if IntOpt(float64(0)) != 0 {
		t.Error("IntOpt(float64 zero)")
	}
}

func TestParseBool(t *testing.T) {
	cases := []struct {
		s    string
		want bool
		ok   bool
	}{
		{"true", true, true},
		{"false", false, true},
		{"1", true, true},
		{"0", false, true},
		{"t", true, true},
		{"f", false, true},
		{"maybe", false, false},
	}
	for _, c := range cases {
		got, err := ParseBool(c.s)
		if c.ok && err != nil {
			t.Errorf("ParseBool(%q) unexpected err: %v", c.s, err)
			continue
		}
		if !c.ok && err == nil {
			t.Errorf("ParseBool(%q) expected err", c.s)
			continue
		}
		if c.ok && got != c.want {
			t.Errorf("ParseBool(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}

func TestCharacterCount(t *testing.T) {
	if CharacterCount([]byte("hello")) != 5 {
		t.Error("CharacterCount hello")
	}
	if CharacterCount([]byte("")) != 0 {
		t.Error("CharacterCount empty")
	}
	if CharacterCountInString("hello") != 5 {
		t.Error("CharacterCountInString hello")
	}
}

func TestHasTrailingWhitespace(t *testing.T) {
	if !HasTrailingWhitespace([]byte("hi  ")) {
		t.Error("expected trailing whitespace")
	}
	if HasTrailingWhitespace([]byte("hi")) {
		t.Error("unexpected trailing whitespace")
	}
}
