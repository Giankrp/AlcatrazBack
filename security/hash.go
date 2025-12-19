package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParams = ArgonParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func generatedRandomBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)

	_, err := rand.Read(bytes)

	return bytes, err
}

func HashPassword(password string) (string, error) {
	params := DefaultParams
	salt, err := generatedRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	final := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		params.Memory, params.Iterations, params.Parallelism, encodedSalt, encodedHash)
	return final, nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid hash format")
	}
	// partes: ["", "argon2id", "v=19", "m=..,t=..,p=..", salt, hash]
	paramsStr := parts[3]
	saltStr := parts[4]
	hashStr := parts[5]

	var m, t uint32
	var p uint8
	_, err := fmt.Sscanf(paramsStr, "m=%d,t=%d,p=%d", &m, &t, &p)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, err
	}
	expected, err := base64.RawStdEncoding.DecodeString(hashStr)
	if err != nil {
		return false, err
	}
	computed := argon2.IDKey([]byte(password), salt, t, m, p, uint32(len(expected)))
	return subtleConstantTimeCompare(expected, computed), nil
}

func NeedsRehash(encoded string, current ArgonParams) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid hash format")
	}
	var m, t uint32
	var p uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p)
	if err != nil {
		return false, err
	}
	return m != current.Memory || t != current.Iterations || p != current.Parallelism, nil
}

func subtleConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var diff byte
	for i := range a {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}
