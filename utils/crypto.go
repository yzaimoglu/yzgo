package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// HashPassword returns a Argon2id hash of a password with the given parameters.
func HashPassword(password string) (encodedHash string, err error) {
	p := &Argon2Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}

	salt, err := GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// generateRandomBytes returns securely generated random bytes. It will return
// an error if the system's secure random number generator fails to function
// correctly, in which case the caller should not continue.
func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// CheckPassword checks that the provided password matches the hashed password.
func CheckPassword(password string, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := DecodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// decodeHash extracts the parameters, salt and derived key from an encoded
// argon2 password hash.
func DecodeHash(encodedHash string) (p *Argon2Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// Hash with SHA512
func HashSHA512(input string) string {
	// Get the bytes of the input
	inputBytes := []byte(input)
	// Create the SHA512 Hasher
	sha512obj := sha512.New()

	// Write input bytes to the hasher
	sha512obj.Write(inputBytes)

	// Get the hashed bytes
	hashedInputBytes := sha512obj.Sum(nil)

	// Convert hashed bytes to hex string
	hashedInputHex := hex.EncodeToString(hashedInputBytes)

	return hashedInputHex
}

// Hash with SHA512
func HashBytesSHA512(inputBytes []byte) string {
	// Create the SHA512 Hasher
	sha512obj := sha512.New()

	// Write input bytes to the hasher
	sha512obj.Write(inputBytes)

	// Get the hashed bytes
	hashedInputBytes := sha512obj.Sum(nil)

	// Convert hashed bytes to hex string
	hashedInputHex := hex.EncodeToString(hashedInputBytes)

	return hashedInputHex
}

func HashMD5(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// Generate a random UUID for the session
func GenerateRandomSession() []byte {
	uuid := uuid.New()
	return []byte(uuid.String())
}

// Generate a random alphanumeric string
func RandomAlphanumeric(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[r.Int64()]
	}
	return string(result), nil
}
