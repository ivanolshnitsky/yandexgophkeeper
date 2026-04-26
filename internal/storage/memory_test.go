package storage

import "testing"

func TestMemoryStorage_CRUD(t *testing.T) {
	s := NewMemory()

	d := Data{
		ID:   "1",
		User: "user",
		Type: TypeText,
	}

	// Save
	s.Save("user", d)

	if len(s.List("user")) != 1 {
		t.Fatal("expected 1 item after save")
	}

	// GetByID
	got, ok := s.GetByID("user", "1")
	if !ok || got.ID != "1" {
		t.Fatal("failed to get by id")
	}

	// Update
	d.Meta = "updated"
	ok = s.Update("user", d)
	if !ok {
		t.Fatal("update failed")
	}

	got, _ = s.GetByID("user", "1")
	if got.Meta != "updated" {
		t.Fatal("update not applied")
	}

	// Delete
	ok = s.Delete("user", "1")
	if !ok {
		t.Fatal("delete failed")
	}

	if len(s.List("user")) != 0 {
		t.Fatal("expected empty after delete")
	}
}
