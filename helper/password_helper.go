package helper

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"math/big"
	"strings"
)

type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// GeneratePassword is used to generate a new password hash for storing and
// comparing at a later date.
func GeneratePassword(c *PasswordConfig, password string) (map[string]string, error) {

	// Generate a Salt
	var keySize, _ = rand.Int(rand.Reader, big.NewInt(1000))
	salt := make([]byte, keySize.Int64())
	var full = make(map[string]string)
	if _, err := rand.Read(salt); err != nil {
		return full, err
	}

	hash := argon2.IDKey([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	//format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full["password"] = b64Hash
	full["salt"] = b64Salt

	return full, nil
}

// ComparePassword is used to compare a user-inputted password to a hash to see
// if the password matches or not.
func ComparePassword(password, hash string)(bool, error) {

	parts := strings.Split(hash, "$")

	c := &PasswordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Time, &c.Threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	c.KeyLen = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen)
	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
