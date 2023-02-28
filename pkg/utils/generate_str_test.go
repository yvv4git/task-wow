package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yvv4git/task-wow/pkg/utils"
)

func TestGenerateStr(t *testing.T) {
	// Check size.
	result, err := utils.GenerateStr(5)
	require.NoError(t, err)
	assert.Equal(t, 5, len(result))

	// Not deterministic.
	result2, err := utils.GenerateStr(5)
	require.NoError(t, err)
	assert.NotEqual(t, result, result2)
}
