package gitter

import (
	"testing"
)

var (
	token = "4b409f3d662592192095055ac603eaf106b0b92b"
)

func TestNewGitter(t *testing.T) {
	gitter, err := NewGitter(token)
	if err != nil {
		t.Error(err)
	}

	gitter.Initialize()
}
