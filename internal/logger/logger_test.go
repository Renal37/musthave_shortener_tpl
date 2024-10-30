package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	// Проверяем успешную инициализацию с уровнем "debug"
	err := Initialize("debug")
	assert.NoError(t, err)
	assert.NotNil(t, Log)
}
