package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("henri@henri.com", "henri", "password", "https://example.com/", "https://example.com/", "Hey, I'm Henri!")

	assert.Nil(t, err)
	fmt.Print("%+v\n", user)
}
