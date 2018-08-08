package model

import (
	"context"
	"time"
)

// Blogpost entity
type Blogpost struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	UserID    int64     `json:"user_id"`
	CreatedOn time.Time `json:"created_on"`
}

// PreInsert hook which sets the CreatedOn date
func (b *Blogpost) PreInsert(ctx context.Context) error {
	b.CreatedOn = time.Now().UTC()
	return nil
}
