package state

import "testing"

func TestMemoryStoreSaveLoad(t *testing.T) {
	store := NewMemoryStore()
	if err := store.Save("execution:1", []byte("ok")); err != nil {
		t.Fatalf("save: %v", err)
	}

	raw, err := store.Load("execution:1")
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if string(raw) != "ok" {
		t.Fatalf("unexpected value: %s", string(raw))
	}
}
