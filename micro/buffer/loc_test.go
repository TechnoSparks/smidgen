package buffer

import "testing"

func TestLocLessThan(t *testing.T) {
	if !(Loc{0, 0}.LessThan(Loc{1, 0})) {
		t.Error("LessThan x")
	}
	if !(Loc{0, 0}.LessThan(Loc{0, 1})) {
		t.Error("LessThan y")
	}
	if Loc{1, 0}.LessThan(Loc{0, 0}) {
		t.Error("not less")
	}
	if Loc{0, 0}.LessThan(Loc{0, 0}) {
		t.Error("equal not less")
	}
}

func TestLocGreaterThan(t *testing.T) {
	if !(Loc{1, 0}.GreaterThan(Loc{0, 0})) {
		t.Error("GreaterThan x")
	}
	if !(Loc{0, 1}.GreaterThan(Loc{0, 0})) {
		t.Error("GreaterThan y")
	}
	if Loc{0, 0}.GreaterThan(Loc{1, 0}) {
		t.Error("not greater")
	}
}

func TestLocClamp(t *testing.T) {
	start := Loc{2, 0}
	end := Loc{8, 0}
	if got := Loc{0, 0}.Clamp(start, end); got != start {
		t.Errorf("clamp low = %+v, want %+v", got, start)
	}
	if got := Loc{10, 0}.Clamp(start, end); got != end {
		t.Errorf("clamp high = %+v, want %+v", got, end)
	}
	if got := Loc{5, 0}.Clamp(start, end); got != (Loc{5, 0}) {
		t.Errorf("clamp mid = %+v", got)
	}
}

func TestLocGreaterEqualLessEqual(t *testing.T) {
	if !(Loc{1, 0}.GreaterEqual(Loc{1, 0})) {
		t.Error("GreaterEqual equal")
	}
	if !(Loc{2, 0}.GreaterEqual(Loc{1, 0})) {
		t.Error("GreaterEqual greater")
	}
	if Loc{0, 0}.GreaterEqual(Loc{1, 0}) {
		t.Error("GreaterEqual false")
	}
	if !(Loc{1, 0}.LessEqual(Loc{1, 0})) {
		t.Error("LessEqual equal")
	}
	if !(Loc{0, 0}.LessEqual(Loc{1, 0})) {
		t.Error("LessEqual less")
	}
	if Loc{2, 0}.LessEqual(Loc{1, 0}) {
		t.Error("LessEqual false")
	}
}
