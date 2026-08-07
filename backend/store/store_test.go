package store

import (
	"path/filepath"
	"testing"
)

func TestSaveRecordPersistsAndUpdatesByID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "qsl-mail-data.json")
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	first, err := database.SaveRecord("tasks", map[string]any{"id": "task-1", "title": "first"})
	if err != nil {
		t.Fatal(err)
	}
	if first["id"] != "task-1" {
		t.Fatalf("unexpected saved id: %v", first["id"])
	}

	if _, err := database.SaveRecord("tasks", map[string]any{"id": "task-1", "title": "updated"}); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	items, err := reopened.Collection("tasks")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0]["title"] != "updated" {
		t.Fatalf("unexpected persisted tasks: %#v", items)
	}
}

func TestUnknownCollectionReturnsError(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.Collection("unknown"); err == nil {
		t.Fatal("expected an error for an unknown collection")
	}
}
