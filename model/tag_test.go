package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sinnott74/goblogserver/orm"
)

func TestFindByTagName(t *testing.T) {

	transaction, ctx := orm.NewTransaction(context.Background())
	defer transaction.Rollback()

	tagName := "TestTag"
	tag := &Tag{Name: tagName}
	err := orm.Insert(ctx, tag)
	assert.NoError(t, err, "Error during Tag insert")

	foundTag, err := findByTagName(ctx, tagName)
	assert.NoError(t, err, "Error during Tag findByTagName")
	assert.Equal(t, tag, foundTag, "Insert tag is not the same as tag found by nam")
}
