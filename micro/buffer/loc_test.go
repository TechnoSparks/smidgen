package buffer

import "testing"

func TestLocLessThan(t *testing.T) {
	a := Loc{0, 0}
	b := Loc{1, 0}
	c := Loc{0, 1}
	if !a.LessThan(b) {
		t.Error("LessThan x")
	}
	if !a.LessThan(c) {
		t.Error("LessThan y")
	}
	if b.LessThan(a) {
		t.Error("not less")
	}
	if a.LessThan(a) {
		t.Error("equal not less")
	}
}

func TestLocGreaterThan(t *testing.T) {
	a := Loc{0, 0}
	b := Loc{1, 0}
	c := Loc{0, 1}
	if !b.GreaterThan(a) {
		t.Error("GreaterThan x")
	}
	if !c.GreaterThan(a) {
		t.Error("GreaterThan y")
	}
	if a.GreaterThan(b) {
		t.Error("not greater")
	}
}

func TestLocClamp(t *testing.T) {
	start := Loc{2, 0}
	end := Loc{8, 0}
	low := Loc{0, 0}
	high := Loc{10, 0}
	mid := Loc{5, 0}
	if got := low.Clamp(start, end); got != start {
		t.Errorf("clamp low = %+v, want %+v", got, start)
	}
	if got := high.Clamp(start, end); got != end {
		t.Errorf("clamp high = %+v, want %+v", got, end)
	}
	if got := mid.Clamp(start, end); got != mid {
		t.Errorf("clamp mid = %+v", got)
	}
}

func TestLocGreaterEqualLessEqual(t *testing.T) {
	a := Loc{1, 0}
	b := Loc{2, 0}
	c := Loc{0, 0}
	if !a.GreaterEqual(a) {
		t.Error("GreaterEqual equal")
	}
	if !b.GreaterEqual(a) {
		t.Error("GreaterEqual greater")
	}
	if c.GreaterEqual(a) {
		t.Error("GreaterEqual false")
	}
	if !a.LessEqual(a) {
		t.Error("LessEqual equal")
	}
	if !c.LessEqual(a) {
		t.Error("LessEqual less")
	}
	if b.LessEqual(a) {
		t.Error("LessEqual false")
	}
}
