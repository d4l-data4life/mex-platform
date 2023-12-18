package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

func redisAuthHashName(challenge string) string {
	return fmt.Sprintf("oauth:challenge:%s", challenge)
}

func redisAuthRefreshTokenHashName(refreshToken string) string {
	h := sha256.New()
	h.Write([]byte(refreshToken))

	return fmt.Sprintf("oauth:refresh:%s", b64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func fragmentOrQuery(responseMode string) string {
	switch responseMode {
	case responseModeQuery:
		return "?"
	case responseModeFragment:
		return "#"
	default:
		return "?"
	}
}

func parseSigningKey(privateKeyPEM []byte, signatureAlg string) (jwk.Key, error) {
	block, _ := pem.Decode(privateKeyPEM)
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	signingKey, err := jwk.FromRaw(rsaKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key from raw: %w", err)
	}

	_ = signingKey.Set("kid", signingKeyID)
	_ = signingKey.Set("use", "sig")
	_ = signingKey.Set("alg", signatureAlg)

	return signingKey, nil
}

func sha256B64Safe(x string) string {
	h := sha256.New()
	h.Write([]byte(x))
	z := b64.StdEncoding.EncodeToString(h.Sum(nil))
	z = strings.ReplaceAll(z, "/", "_")
	z = strings.ReplaceAll(z, "+", "-")
	z = strings.ReplaceAll(z, "=", "")
	return z
}

func randomString() string {
	t := make([]byte, refreshTokenLength)
	_, _ = rand.Read(t)
	return hex.EncodeToString(t)
}
