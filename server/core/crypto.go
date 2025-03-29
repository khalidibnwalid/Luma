package core

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Cryptographically secure random bytes generator
func GenerateRandomBytes(length uint32) ([]byte, error) {
	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return nil, err
	}
	return buffer, nil
}

type hashingOptions struct {
	iterations  uint32
	memory      uint32
	parallelism uint8
	keyLength   uint32
	saltLength  uint32
}

var defaultHashingOptions = &hashingOptions{
	memory:      4 * 1024,
	iterations:  3,
	parallelism: 1,
	keyLength:   32,
	saltLength:  16,
}

// Hashing with Argon2id
func CreateHashWithSalt(value string, providedSalt ...[]byte) ([]byte, []byte, error) {
	var (
		err  error
		salt []byte
		opt  = defaultHashingOptions
	)

	if len(providedSalt) > 0 {
		salt = providedSalt[0]
	} else {
		if salt, err = GenerateRandomBytes(opt.saltLength); err != nil {
			return nil, nil, err
		}
	}

	hash := argon2.IDKey([]byte(value), salt, opt.iterations, opt.memory, opt.parallelism, opt.keyLength)
	return hash, salt, nil
}

// Serialize the hash and salt into a string format
func SerializeHashWithSalt(hash []byte, salt []byte) string {
	opt := defaultHashingOptions
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	serializedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, opt.memory, opt.iterations, opt.parallelism, b64Salt, b64Hash)
	return serializedHash
}

var (
	ErrInvalidHashFormat      = errors.New("invalid hash format")
	ErrHashVerificationFailed = errors.New("hash verification failed")
)

func DeserializeHash(serializedHash string) ([]byte, []byte, error) {
	parts := strings.Split(serializedHash, "$")
	if len(parts) != 6 {
		return nil, nil, ErrInvalidHashFormat
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, err
	}

	return hash, salt, nil
}

// Verify a row unhashed value against the stored hash.
func VerifyHashWithSalt(raw, dataHash string) error {
	storedHash, salt, err := DeserializeHash(dataHash)
	if err != nil {
		return err
	}

	toHash, _, err := CreateHashWithSalt(raw, salt)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(storedHash, []byte(toHash)) == 1 {
		return nil
	}
	return ErrHashVerificationFailed
}
