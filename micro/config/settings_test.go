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

func TestDefaultSettingsNonEmpty(t *testing.T) {
	if len(DefaultAllSettings()) == 0 {
		t.Fatal("DefaultAllSettings empty")
	}
	if len(DefaultCommonSettings()) == 0 {
		t.Fatal("DefaultCommonSettings empty")
	}
	if len(DefaultGlobalOnlySettings) == 0 {
		t.Fatal("DefaultGlobalOnlySettings empty")
	}
}
