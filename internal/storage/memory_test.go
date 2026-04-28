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
	if err := s.Save("user", d); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	list, err := s.List("user")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatal("expected 1 item after save")
	}

	// GetByID
	got, ok, err := s.GetByID("user", "1")
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if !ok || got.ID != "1" {
		t.Fatal("failed to get by id")
	}

	// Update
	d.Meta = "updated"

	ok, err = s.Update("user", d)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if !ok {
		t.Fatal("update returned false")
	}

	got, ok, err = s.GetByID("user", "1")
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}
	if !ok || got.Meta != "updated" {
		t.Fatal("update not applied")
	}

	// Delete
	ok, err = s.Delete("user", "1")
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if !ok {
		t.Fatal("delete returned false")
	}

	list, err = s.List("user")
	if err != nil {
		t.Fatalf("list after delete failed: %v", err)
	}

	if len(list) != 0 {
		t.Fatal("expected empty after delete")
	}
}
