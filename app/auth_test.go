package app

import (
	"golang.org/x/crypto/bcrypt"
	"net/http/httptest"
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

func TestExtractToken(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost:8989?token=1234567890", nil)
	TokenFromParam := ExtractToken(request)
	if TokenFromParam != "1234567890" {
		t.Error("Get token from param fail")
	}

	request = httptest.NewRequest("POST", "http://localhost:8989", nil)
	request.Header.Add("Authorization", "Bearer 0987654321")
	bearerToken := ExtractToken(request)
	if bearerToken != "0987654321" {
		t.Error("Get bearer token fail")
	}
}
