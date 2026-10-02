package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"familyquest-backend/internal/domain/entity"
)

type HMACGenerator struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewHMACGenerator(secret string, issuer string, ttl time.Duration) *HMACGenerator {
	return &HMACGenerator{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (g *HMACGenerator) Generate(user entity.User, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(g.ttl)
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := map[string]any{
		"sub":   user.ID.String(),
		"login": user.Login,
		"name":  user.DisplayName,
		"iss":   g.issuer,
		"iat":   now.Unix(),
		"exp":   expiresAt.Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal token header: %w", err)
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal token claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsignedToken := strings.Join([]string{encodedHeader, encodedClaims}, ".")

	mac := hmac.New(sha256.New, g.secret)
	mac.Write([]byte(unsignedToken))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsignedToken + "." + signature, expiresAt, nil
}
