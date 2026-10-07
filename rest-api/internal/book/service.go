package book

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrNotFound = errors.New("book not found")

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]Book, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int64) (Book, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, input Input) (Book, error) {
	normalized, err := normalizeAndValidate(input)
	if err != nil {
		return Book{}, err
	}

	return s.repo.Create(ctx, normalized)
}

func (s *Service) Update(ctx context.Context, id int64, input Input) (Book, error) {
	normalized, err := normalizeAndValidate(input)
	if err != nil {
		return Book{}, err
	}

	return s.repo.Update(ctx, id, normalized)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func normalizeAndValidate(input Input) (Input, error) {
	// Normalisieren vor der Validierung, damit "  Titel  " als gueltiger Inhalt zaehlt.
	input.Title = strings.TrimSpace(input.Title)
	input.Author = strings.TrimSpace(input.Author)

	if input.Title == "" {
		return Input{}, &ValidationError{Message: "title must not be empty"}
	}
	if input.Author == "" {
		return Input{}, &ValidationError{Message: "author must not be empty"}
	}
	if input.PublishedYear != nil && *input.PublishedYear <= 0 {
		return Input{}, &ValidationError{
			Message: fmt.Sprintf("published_year must be greater than 0, got %d", *input.PublishedYear),
		}
	}

	return input, nil
}
