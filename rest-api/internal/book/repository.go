package book

import "context"

type Repository interface {
	List(ctx context.Context) ([]Book, error)
	Get(ctx context.Context, id int64) (Book, error)
	Create(ctx context.Context, input Input) (Book, error)
	Update(ctx context.Context, id int64, input Input) (Book, error)
	Delete(ctx context.Context, id int64) error
}
