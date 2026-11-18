package core_encrypt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

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
func GeneratePasswordHash(val string) (string, error) {
	hashedVal, err := bcrypt.GenerateFromPassword([]byte(val), 8)
	if err != nil {
		return "", fmt.Errorf("hash error")
	}

	return string(hashedVal), nil
}

// CompareHashAndPassword сравнивает хэш с паролем
//
// Returns: ошибку wrong password если пароль не совпадает
func CompareHashAndPassword(hash string, val string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(val),
	); err != nil {
		return fmt.Errorf("wrong password: %w", core_errors.ErrInvalidCredentials)
	}

	return nil
}

// GenerateVerificationCode генерирует криптографически безопасный N-значный цифровой код
func GenerateVerificationCode(length int) (string, error) {
	const digits = "0123456789"
	var code strings.Builder

	for i := 0; i < length; i++ {
		// Используем crypto/rand для максимальной безопасности (вместо math/rand)
		maxIdx := big.NewInt(int64(len(digits)))
		idx, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return "", err
		}
		code.WriteByte(digits[idx.Int64()])
	}

	return code.String(), nil
}

// HashToken хэширует токен алгоритов SHA256
//
// Returns: string

func (e *Encryptor) HashToken(token string) string {
	h := hmac.New(sha256.New, []byte(e.secret))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
