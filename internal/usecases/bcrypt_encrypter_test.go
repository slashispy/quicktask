package usecases

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptEncrypter_Encrypt(t *testing.T) {
	enc := &BcryptEncrypter{}

	hash, err := enc.Encrypt("supersecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected a non-empty hash")
	}
	if hash == "supersecret" {
		t.Fatal("hash must not equal the plaintext password")
	}

	// The produced hash must validate against the original password.
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("supersecret")); err != nil {
		t.Errorf("hash does not verify against original password: %v", err)
	}
	// And must fail for a wrong password.
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("wrong")); err == nil {
		t.Error("hash verified against a wrong password")
	}
}

func TestBcryptEncrypter_ProducesDifferentHashes(t *testing.T) {
	enc := &BcryptEncrypter{}

	h1, err := enc.Encrypt("samepass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, err := enc.Encrypt("samepass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// bcrypt uses a random salt, so identical inputs yield different hashes.
	if h1 == h2 {
		t.Error("expected different hashes for the same password (salted)")
	}
}
