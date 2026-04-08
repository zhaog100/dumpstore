package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters following OWASP recommendations
const (
	argon2idTime    = 3
	argon2idMemory  = 64 * 1024 // 64MB
	argon2idThreads = 4
	argon2idKeyLen  = 32
	argon2idSaltLen = 16
)

// ErrInvalidHash is returned when the hash format is invalid
var ErrInvalidHash = errors.New("invalid hash format")

// ErrIncompatibleVersion is returned when the hash version is incompatible
var ErrIncompatibleVersion = errors.New("incompatible hash version")

// HashPasswordArgon2id generates an argon2id hash of the password
// Returns a PHC string format: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
func HashPasswordArgon2id(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, argon2idSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	// Generate hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2idTime,
		argon2idMemory,
		argon2idThreads,
		argon2idKeyLen,
	)

	// Encode as PHC string
	// Format: $argon2id$v=19$m=65536,t=3,p=4$<base64(salt)>$<base64(hash)>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2idMemory,
		argon2idTime,
		argon2idThreads,
		b64Salt,
		b64Hash,
	), nil
}

// VerifyPasswordArgon2id verifies a password against an argon2id hash
func VerifyPasswordArgon2id(password, encodedHash string) error {
	// Parse PHC string
	salt, hash, params, err := decodeArgon2idHash(encodedHash)
	if err != nil {
		return err
	}

	// Generate hash from password
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.time,
		params.memory,
		params.threads,
		params.keyLen,
	)

	// Compare hashes (constant-time comparison)
	if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
		return errors.New("invalid password")
	}

	return nil
}

// IsArgon2idHash checks if a hash is in argon2id format
func IsArgon2idHash(hash string) bool {
	return strings.HasPrefix(hash, "$argon2id$")
}

// argon2idParams holds the parameters extracted from a PHC string
type argon2idParams struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// decodeArgon2idHash parses a PHC format argon2id hash
func decodeArgon2idHash(encodedHash string) (salt, hash []byte, params *argon2idParams, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidHash
	}

	if parts[2] != "v=19" {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	// Parse parameters: m=65536,t=3,p=4
	params = &argon2idParams{}
	paramParts := strings.Split(parts[3], ",")
	if len(paramParts) != 3 {
		return nil, nil, nil, ErrInvalidHash
	}

	// Parse m (memory)
	if _, err := fmt.Sscanf(paramParts[0], "m=%d", &params.memory); err != nil {
		return nil, nil, nil, fmt.Errorf("parse memory: %w", err)
	}

	// Parse t (time)
	if _, err := fmt.Sscanf(paramParts[1], "t=%d", &params.time); err != nil {
		return nil, nil, nil, fmt.Errorf("parse time: %w", err)
	}

	// Parse p (threads)
	var threads uint32
	if _, err := fmt.Sscanf(paramParts[2], "p=%d", &threads); err != nil {
		return nil, nil, nil, fmt.Errorf("parse threads: %w", err)
	}
	params.threads = uint8(threads)

	// Decode salt
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decode salt: %w", err)
	}

	// Decode hash
	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decode hash: %w", err)
	}

	params.keyLen = uint32(len(hash))

	return salt, hash, params, nil
}
