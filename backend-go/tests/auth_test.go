package tests

import (
	"testing"

	"careerflow/backend/internal/auth"
)

func TestPasswordHashing(t *testing.T) {
	hash, err := auth.HashPassword("supersecret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := auth.CheckPassword(hash, "supersecret"); err != nil {
		t.Fatalf("expected password to match, got %v", err)
	}
}
