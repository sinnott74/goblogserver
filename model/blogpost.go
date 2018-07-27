package model

import (
	"context"
	"time"
)

type BlogPost struct {
	ID        int64
	Title     string
	Text      string
	UserID    int64
	CreatedOn time.Time
}

func (b *BlogPost) PreInsert(ctx context.Context) error {
	b.CreatedOn = time.Now().UTC()
	return nil
}
