package scaffolder

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// randomToken returns n random bytes, base64url-encoded without padding, so the
// value needs no quoting in .env, the Makefile or docker-compose.yml.
func randomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b) // crypto/rand.Read never fails (it crashes the program instead)
	return base64.RawURLEncoding.EncodeToString(b)
}

// hashPassword encodes password exactly the way the generated auth adapter
// verifies AUTH_PASSWORD_HASH (internal/adapters/handlers/auth): PBKDF2-HMAC-SHA256,
// 600,000 iterations, a 16-byte salt and a 32-byte key, written as
// "pbkdf2-sha256:<iterations>:<salt>:<key>" in base64url. Keep the two in step.
func hashPassword(password string) (string, error) {
	const iterations = 600000
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	key, err := pbkdf2.Key(sha256.New, password, salt, iterations, 32)
	if err != nil {
		return "", err
	}
	enc := base64.RawURLEncoding
	return fmt.Sprintf("pbkdf2-sha256:%d:%s:%s", iterations, enc.EncodeToString(salt), enc.EncodeToString(key)), nil
}
