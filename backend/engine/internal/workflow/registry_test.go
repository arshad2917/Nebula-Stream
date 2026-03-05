package workflow

import "testing"

func TestRegistryUpsertAndActive(t *testing.T) {
	r := NewRegistry(Definition{Name: "hello", Version: "v1", Triggers: []Trigger{{Type: "manual"}}, Steps: []Step{{ID: "s1", Type: "builtin.log"}}})

	if _, ok := r.Active(); !ok {
		t.Fatal("expected active workflow")
	}

	r.Upsert(Definition{Name: "other", Version: "v1", Triggers: []Trigger{{Type: "manual"}}, Steps: []Step{{ID: "s1", Type: "builtin.log"}}})
	if err := r.SetActive("other"); err != nil {
		t.Fatalf("set active: %v", err)
	}

	active, ok := r.Active()
	if !ok || active.Name != "other" {
		t.Fatalf("unexpected active workflow: %+v", active)
	}
}
