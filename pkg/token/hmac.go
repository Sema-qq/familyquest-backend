package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

type HMACGenerator struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type claims struct {
	Subject string `json:"sub"`
	Issuer  string `json:"iss"`
	Expires int64  `json:"exp"`
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
		return "", time.Time{}, fmt.Errorf("can't marshal token header: %w", err)
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("can't marshal token claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsignedToken := strings.Join([]string{encodedHeader, encodedClaims}, ".")

	mac := hmac.New(sha256.New, g.secret)
	mac.Write([]byte(unsignedToken))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsignedToken + "." + signature, expiresAt, nil
}

func (g *HMACGenerator) Parse(token string) (entity.UserID, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return entity.UserID{}, errors.New("invalid token")
	}

	unsignedToken := strings.Join(parts[:2], ".")
	if !g.validSignature(unsignedToken, parts[2]) {
		return entity.UserID{}, errors.New("invalid token signature")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return entity.UserID{}, fmt.Errorf("can't decode token claims: %w", err)
	}

	var tokenClaims claims
	if err = json.Unmarshal(claimsBytes, &tokenClaims); err != nil {
		return entity.UserID{}, fmt.Errorf("can't unmarshal token claims: %w", err)
	}

	if tokenClaims.Issuer != g.issuer {
		return entity.UserID{}, errors.New("invalid token issuer")
	}
	if time.Now().Unix() >= tokenClaims.Expires {
		return entity.UserID{}, errors.New("token expired")
	}

	userID, err := uuid.Parse(tokenClaims.Subject)
	if err != nil {
		return entity.UserID{}, fmt.Errorf("can't parse token subject: %w", err)
	}

	return entity.UserID(userID), nil
}

func (g *HMACGenerator) validSignature(unsignedToken string, signature string) bool {
	mac := hmac.New(sha256.New, g.secret)
	mac.Write([]byte(unsignedToken))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
