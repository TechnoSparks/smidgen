package config

import "testing"

func TestLockBindingsDefaultExists(t *testing.T) {
	v, ok := defaultCommonSettings["lockbindings"]
	if !ok {
		t.Fatal("lockbindings missing from defaultCommonSettings")
	}
	b, ok := v.(bool)
	if !ok || b {
		t.Fatalf("lockbindings default should be false, got %v", v)
	}
}

func TestDefaultSettingsNonEmpty(t *testing.T) {
	if len(DefaultCommonSettings()) == 0 {
		t.Fatal("DefaultCommonSettings empty")
	}
	if len(defaultCommonSettings) == 0 {
		t.Fatal("defaultCommonSettings empty")
	}
}
