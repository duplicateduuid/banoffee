package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	user, err := NewUser("henri@henri.com", "henri", "password", nil, nil, nil)

	assert.Nil(t, err)
	assert.Nil(t, user.AvatarUrl)
	assert.Nil(t, user.HeaderUrl)
	assert.Nil(t, user.Bio)
}
