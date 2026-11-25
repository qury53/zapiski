package core_encrypt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

type Encryptor struct {
	secret string
}

func NewEncryptor(config *Config) *Encryptor {
	return &Encryptor{
		secret: config.EncryptionSecret,
	}
}

// GeneratePasswordHash генерит хэш пароля
//
// Returns: хэш пароля
func (e *Encryptor) GeneratePasswordHash(val string) (string, error) {
	hashedVal, err := bcrypt.GenerateFromPassword([]byte(val), 10)
	if err != nil {
		return "", fmt.Errorf("hash error: %w", err)
	}

	return string(hashedVal), nil
}

// CompareHashAndPassword сравнивает хэш с паролем
//
// Returns: ошибку wrong password если пароль не совпадает
func (e *Encryptor) CompareHashAndPassword(hash string, val string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(val),
	); err != nil {
		return fmt.Errorf("wrong password: %w", core_errors.ErrInvalidCredentials)
	}

	return nil
}

// GenerateVerificationCode отвечает за генерацию 6-значного криптостойкого кода
func (e *Encryptor) GenerateVerificationCode() (string, error) {
	max := big.NewInt(900000)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	code := n.Int64() + 100000

	return fmt.Sprintf("%d", code), nil
}

// HashToken хэширует токен алгоритов SHA256
//
// Returns: string

func (e *Encryptor) HashToken(token string) string {
	h := hmac.New(sha256.New, []byte(e.secret))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
