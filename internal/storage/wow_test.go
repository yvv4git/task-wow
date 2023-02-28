package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yvv4git/task-wow/internal/storage"
)

func TestWOW_LoadPhrase(t *testing.T) {
	wowStorage := storage.NewWOW()

	phrase1 := wowStorage.LoadPhrase()
	phrase2 := wowStorage.LoadPhrase()

	assert.NotEqual(t, phrase1, phrase2)
}
