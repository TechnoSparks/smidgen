package config

import "testing"

func TestLockBindingsDefaultExists(t *testing.T) {
	v, ok := DefaultGlobalOnlySettings["lockbindings"]
	if !ok {
		t.Fatal("lockbindings missing from DefaultGlobalOnlySettings")
	}
	b, ok := v.(bool)
	if !ok || b {
		t.Fatalf("lockbindings default should be false, got %v", v)
	}
}
