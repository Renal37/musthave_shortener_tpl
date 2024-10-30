package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	userID := "test-id"
	isNew := true

	user := NewUser(userID, isNew)

	assert.Equal(t, userID, user.ID, "User ID должен совпадать")
	assert.Equal(t, isNew, user.New, "New флаг должен совпадать")
}
