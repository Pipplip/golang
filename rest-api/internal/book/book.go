package book

type Book struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear *int   `json:"published_year,omitempty"`
}

type Input struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear *int   `json:"published_year,omitempty"`
}
