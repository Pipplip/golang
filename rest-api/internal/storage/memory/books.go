package memory

import (
	"context"
	"sort"
	"sync"

	"rest-api/internal/book"
)

type Repository struct {
	// RWMutex erlaubt parallele Leser und exklusiven Schreibzugriff.
	mu     sync.RWMutex
	nextID int64
	books  map[int64]book.Book // enthält die Bücher die verwaltet werden in memory, wobei die ID als Schlüssel dient.
}

func NewRepository() *Repository {
	return &Repository{
		nextID: 1,
		books:  make(map[int64]book.Book),
	}
}

func (r *Repository) List(ctx context.Context) ([]book.Book, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]book.Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, cloneBook(b))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func (r *Repository) Get(ctx context.Context, id int64) (book.Book, error) {
	if err := ctx.Err(); err != nil {
		return book.Book{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	b, ok := r.books[id]
	if !ok {
		return book.Book{}, book.ErrNotFound
	}

	return cloneBook(b), nil
}

func (r *Repository) Create(ctx context.Context, input book.Input) (book.Book, error) {
	if err := ctx.Err(); err != nil {
		return book.Book{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Monoton steigende IDs verhindern Kollisionen nach Delete-Operationen.
	id := r.nextID
	r.nextID++

	b := book.Book{
		ID:            id,
		Title:         input.Title,
		Author:        input.Author,
		PublishedYear: cloneYear(input.PublishedYear),
	}
	r.books[id] = b

	return cloneBook(b), nil
}

func (r *Repository) Update(ctx context.Context, id int64, input book.Input) (book.Book, error) {
	if err := ctx.Err(); err != nil {
		return book.Book{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return book.Book{}, book.ErrNotFound
	}

	updated := book.Book{
		ID:            id,
		Title:         input.Title,
		Author:        input.Author,
		PublishedYear: cloneYear(input.PublishedYear),
	}
	r.books[id] = updated

	return cloneBook(updated), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return book.ErrNotFound
	}

	delete(r.books, id)
	return nil
}

func cloneBook(b book.Book) book.Book {
	// Defensive Kopie: Externe Aufrufer erhalten keine intern geteilten Referenzen.
	return book.Book{
		ID:            b.ID,
		Title:         b.Title,
		Author:        b.Author,
		PublishedYear: cloneYear(b.PublishedYear),
	}
}

func cloneYear(year *int) *int {
	if year == nil {
		return nil
	}
	value := *year
	return &value
}
