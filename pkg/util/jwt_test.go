package util

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	_, err := GenerateToken("xxrl", true)
	assert.Equal(t,err , nil )
}