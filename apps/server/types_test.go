package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("henri@henri.com", "henri", "password", nil, nil, nil)

	assert.Nil(t, err)
	assert.Nil(user.AvatarUrl)
	assert.Nil(user.HeaderUrl)
	assert.Nil(user.Bio)
	fmt.Print("%+v\n", user)
}
