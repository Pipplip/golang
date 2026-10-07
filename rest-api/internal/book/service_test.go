package book

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	createInput Input
	updateInput Input
	updateErr   error
}

func (f *fakeRepo) List(_ context.Context) ([]Book, error) {
	return nil, nil
}

func (f *fakeRepo) Get(_ context.Context, _ int64) (Book, error) {
	return Book{}, nil
}

func (f *fakeRepo) Create(_ context.Context, input Input) (Book, error) {
	f.createInput = input
	return Book{ID: 1, Title: input.Title, Author: input.Author, PublishedYear: input.PublishedYear}, nil
}

func (f *fakeRepo) Update(_ context.Context, _ int64, input Input) (Book, error) {
	f.updateInput = input
	if f.updateErr != nil {
		return Book{}, f.updateErr
	}
	return Book{ID: 7, Title: input.Title, Author: input.Author, PublishedYear: input.PublishedYear}, nil
}

func (f *fakeRepo) Delete(_ context.Context, _ int64) error {
	return nil
}

func TestServiceCreateNormalizesInput(t *testing.T) {
	repo := &fakeRepo{}
	service := NewService(repo)
	year := 2025

	created, err := service.Create(context.Background(), Input{
		Title:         "  Domain-Driven Design  ",
		Author:        "  Eric Evans  ",
		PublishedYear: &year,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if created.Title != "Domain-Driven Design" {
		t.Fatalf("unexpected title: %q", created.Title)
	}
	if created.Author != "Eric Evans" {
		t.Fatalf("unexpected author: %q", created.Author)
	}
	if repo.createInput.Title != "Domain-Driven Design" {
		t.Fatalf("repo received non-normalized title: %q", repo.createInput.Title)
	}
}

func TestServiceCreateRejectsInvalidInput(t *testing.T) {
	service := NewService(&fakeRepo{})

	_, err := service.Create(context.Background(), Input{Title: " ", Author: "Someone"})
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestServiceUpdatePropagatesNotFound(t *testing.T) {
	repo := &fakeRepo{updateErr: ErrNotFound}
	service := NewService(repo)

	_, err := service.Update(context.Background(), 99, Input{Title: "Go", Author: "Team"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
