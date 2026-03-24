package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"gogo/internal/user"
	"strconv"
	"strings"
	"time"
)

var errInvalidJWT = errors.New("invalid jwt")

type tokenClaims struct {
	Subject string
	Email   string
	Type    string
	Issued  int64
	Expires int64
	ID      string
}

func (s *Service) signAccessToken(dbUser user.User) (string, error) {
	now := time.Now().UTC()

	return s.signJWT(tokenClaims{
		Subject: strconv.FormatUint(uint64(dbUser.ID), 10),
		Email:   dbUser.Email,
		Type:    "access",
		Issued:  now.Unix(),
		Expires: now.Add(accessTokenTTL).Unix(),
	})
}

func (s *Service) signRefreshToken(dbUser user.User) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(refreshTokenTTL)
	tokenID, err := generateTokenID()
	if err != nil {
		return "", time.Time{}, err
	}

	token, err := s.signJWT(tokenClaims{
		Subject: strconv.FormatUint(uint64(dbUser.ID), 10),
		Email:   dbUser.Email,
		Type:    "refresh",
		Issued:  now.Unix(),
		Expires: expiresAt.Unix(),
		ID:      tokenID,
	})
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (s *Service) signJWT(claims tokenClaims) (string, error) {
	header, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	payload := map[string]any{
		"sub":   claims.Subject,
		"email": claims.Email,
		"type":  claims.Type,
		"iat":   claims.Issued,
		"exp":   claims.Expires,
	}
	if claims.ID != "" {
		payload["jti"] = claims.ID
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(header)
	claimsJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsignedToken := encodedHeader + "." + encodedClaims

	signature := hmac.New(sha256.New, s.jwtSecret)
	if _, err := signature.Write([]byte(unsignedToken)); err != nil {
		return "", err
	}

	return unsignedToken + "." + base64.RawURLEncoding.EncodeToString(signature.Sum(nil)), nil
}

func (s *Service) parseRefreshToken(token string) (*tokenClaims, error) {
	claims, err := s.parseJWT(token)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, errInvalidJWT
	}

	if claims.Subject == "" || claims.ID == "" {
		return nil, errInvalidJWT
	}

	return claims, nil
}

func (s *Service) parseJWT(token string) (*tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errInvalidJWT
	}

	unsignedToken := parts[0] + "." + parts[1]
	signature := hmac.New(sha256.New, s.jwtSecret)
	if _, err := signature.Write([]byte(unsignedToken)); err != nil {
		return nil, err
	}

	expectedSignature := base64.RawURLEncoding.EncodeToString(signature.Sum(nil))
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, errInvalidJWT
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errInvalidJWT
	}

	var rawClaims struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
		Type    string `json:"type"`
		Issued  int64  `json:"iat"`
		Expires int64  `json:"exp"`
		ID      string `json:"jti"`
	}
	if err := json.Unmarshal(payload, &rawClaims); err != nil {
		return nil, errInvalidJWT
	}

	if rawClaims.Subject == "" || rawClaims.Expires == 0 {
		return nil, errInvalidJWT
	}

	if time.Now().UTC().Unix() >= rawClaims.Expires {
		return nil, errInvalidJWT
	}

	return &tokenClaims{
		Subject: rawClaims.Subject,
		Email:   rawClaims.Email,
		Type:    rawClaims.Type,
		Issued:  rawClaims.Issued,
		Expires: rawClaims.Expires,
		ID:      rawClaims.ID,
	}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func generateTokenID() (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}
