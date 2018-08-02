package model

import (
	"os"
	"testing"

	"github.com/sinnott74/goblogserver/database"
)

func TestMain(m *testing.M) {
	database.Init()
	result := m.Run()
	os.Exit(result)
}
