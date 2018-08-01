package model

import (
	"context"

	"github.com/sinnott74/goblogserver/orm"
)

// Tag entity
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// PreSave checks if this tag already exists on the database
func (tag *Tag) PreSave(ctx context.Context) error {
	storedTag, err := findByTagName(ctx, tag.Name)
	if err != nil {
		return err
	}
	if storedTag != nil {
		tag.ID = storedTag.ID
	}
	return nil
}

// findByTagName searches for a single tag by its name
func findByTagName(ctx context.Context, tagName string) (*Tag, error) {
	tag := Tag{Name: tagName}
	tags := []Tag{}
	err := orm.SelectAll(ctx, &tags, &tag)
	if err != nil || len(tags) == 0 {
		return nil, err
	}
	return &tags[0], nil
}
