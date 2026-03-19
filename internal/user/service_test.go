package user

import (
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestServiceCreateAndGet(t *testing.T) {
	service := newTestService(t)

	created, err := service.Create("Alice", 30)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if created.ID == 0 {
		t.Fatal("Create returned user without ID")
	}

	if created.Name != "Alice" {
		t.Fatalf("Create returned wrong name: got %q", created.Name)
	}

	if created.Age != 30 {
		t.Fatalf("Create returned wrong age: got %d", created.Age)
	}

	fetched, err := service.Get(created.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if fetched.ID != created.ID {
		t.Fatalf("Get returned wrong ID: got %d want %d", fetched.ID, created.ID)
	}

	if fetched.Name != "Alice" {
		t.Fatalf("Get returned wrong name: got %q", fetched.Name)
	}

	if fetched.Age != 30 {
		t.Fatalf("Get returned wrong age: got %d", fetched.Age)
	}
}

func TestServiceUpdate(t *testing.T) {
	service := newTestService(t)

	created, err := service.Create("Alice", 30)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	updated, err := service.Update(created.ID, "Bob", 41)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if updated.ID != created.ID {
		t.Fatalf("Update changed ID: got %d want %d", updated.ID, created.ID)
	}

	if updated.Name != "Bob" {
		t.Fatalf("Update returned wrong name: got %q", updated.Name)
	}

	if updated.Age != 41 {
		t.Fatalf("Update returned wrong age: got %d", updated.Age)
	}

	fetched, err := service.Get(created.ID)
	if err != nil {
		t.Fatalf("Get after Update returned error: %v", err)
	}

	if fetched.Name != "Bob" || fetched.Age != 41 {
		t.Fatalf("Get after Update returned wrong user: %+v", fetched)
	}
}

func TestServiceGetMissingUser(t *testing.T) {
	service := newTestService(t)

	_, err := service.Get(999)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Get missing user error = %v, want %v", err, gorm.ErrRecordNotFound)
	}
}

func TestServiceUpdateMissingUser(t *testing.T) {
	service := newTestService(t)

	_, err := service.Update(999, "Bob", 41)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Update missing user error = %v, want %v", err, gorm.ErrRecordNotFound)
	}
}

func newTestService(t *testing.T) *Service {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	return NewService(db)
}
