package auth

import (
    "testing"
)

func TestHashPasswordArgon2id(t *testing.T) {
    password := "test_password_123"
    hash, err := HashPasswordArgon2id(password)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }

    // Verify hash format
    if !IsArgon2idHash(hash) {
        t.Errorf("Hash is not in argon2id format")
    }

    // Verify password
    err = VerifyPasswordArgon2id(password, hash)
    if err != nil {
        t.Errorf("Failed to verify password: %v", err)
    }

    // Verify wrong password
    err = VerifyPasswordArgon2id("wrong_password", hash)
    if err == nil {
        t.Errorf("Should fail with wrong password")
    }
}

func TestIsArgon2idHash(t *testing.T) {
    tests := []struct {
        hash     string
        expected bool
    }{
        {"$argon2id$v=19$m=65536,t=3,p=4$abc$xyz", true},
        {"$bcrypt$12$abc$xyz", false},
        {"invalid_hash", false},
    }

    for _, tt := range tests {
        result := IsArgon2idHash(tt.hash)
        if result != tt.expected {
            t.Errorf("IsArgon2idHash(%s) = %v, want %v", tt.hash, result, tt.expected)
        }
    }
}
