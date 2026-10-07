package memory

import (
	"context"
	"sync"
	"testing"

	"rest-api/internal/book"
)

func TestRepositoryCRUD(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	created, err := repo.Create(ctx, book.Input{Title: "Go in Action", Author: "K. Kennedy"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != 1 {
		t.Fatalf("expected id 1, got %d", created.ID)
	}

	found, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if found.Title != "Go in Action" {
		t.Fatalf("unexpected title: %q", found.Title)
	}

	updated, err := repo.Update(ctx, created.ID, book.Input{Title: "Go in Action 2nd", Author: "K. Kennedy"})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Title != "Go in Action 2nd" {
		t.Fatalf("unexpected title after update: %q", updated.Title)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := repo.Get(ctx, created.ID); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestRepositoryListSortedByID(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	for _, title := range []string{"C", "A", "B"} {
		if _, err := repo.Create(ctx, book.Input{Title: title, Author: "Author"}); err != nil {
			t.Fatalf("Create returned error: %v", err)
		}
	}

	books, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(books) != 3 {
		t.Fatalf("expected 3 books, got %d", len(books))
	}
	if books[0].ID != 1 || books[1].ID != 2 || books[2].ID != 3 {
		t.Fatalf("expected IDs [1 2 3], got [%d %d %d]", books[0].ID, books[1].ID, books[2].ID)
	}
}

func TestRepositoryConcurrentCreateUsesUniqueIDs(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	const workers = 50
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			_, _ = repo.Create(ctx, book.Input{
				Title:  "Book",
				Author: "Author",
			})
		}(i)
	}
	wg.Wait()

	books, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(books) != workers {
		t.Fatalf("expected %d books, got %d", workers, len(books))
	}

	seen := make(map[int64]struct{}, workers)
	for _, b := range books {
		if _, exists := seen[b.ID]; exists {
			t.Fatalf("duplicate ID detected: %d", b.ID)
		}
		seen[b.ID] = struct{}{}
	}
}
