package buffer

import "testing"

func TestFindNextLiteral(t *testing.T) {
	b := NewBufferFromString("hello world\nfoo bar", "")
	locs, found, err := b.FindNext("world", Loc{0, 0}, Loc{1000, 1000}, Loc{0, 0}, true, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("expected to find 'world'")
	}
	if locs[0] != (Loc{6, 0}) || locs[1] != (Loc{11, 0}) {
		t.Fatalf("unexpected match location: %+v", locs)
	}
}

func TestFindNextRegex(t *testing.T) {
	b := NewBufferFromString("abc123def\nxyz", "")
	locs, found, err := b.FindNext(`\d+`, Loc{0, 0}, Loc{1000, 1000}, Loc{0, 0}, true, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("expected regex match")
	}
	if locs[0] != (Loc{3, 0}) || locs[1] != (Loc{6, 0}) {
		t.Fatalf("unexpected match location: %+v", locs)
	}
}

func TestFindNextNotFound(t *testing.T) {
	b := NewBufferFromString("hello world", "")
	_, found, err := b.FindNext("zzz", Loc{0, 0}, Loc{1000, 1000}, Loc{0, 0}, true, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("expected no match")
	}
}
