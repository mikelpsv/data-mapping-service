package app

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHash(t *testing.T) {
	var test_password = "test_password"
	hash, err := Hash(test_password)
	if err != nil {
		t.Error(err != nil)
	}
	if string(hash) == "" {
		t.Error("Hash incorrect")
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(test_password))
	if err != nil {
		t.Error("Hash incorrect")
	}
}
