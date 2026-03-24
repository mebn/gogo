package auth

import (
	"errors"
	"gogo/internal/user"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrEmailAlreadyRegistered = errors.New("email already registered")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

const (
	defaultJWTSecret = "change-me-in-production"
	accessTokenTTL   = time.Hour
	refreshTokenTTL  = 24 * 30 * time.Hour
)

type Service struct {
	db        *gorm.DB
	jwtSecret []byte
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		jwtSecret: []byte(defaultJWTSecret),
	}
}

func (s *Service) Register(email string, password string, name string, age *uint) (*AuthTokens, error) {
	var existingUser user.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, ErrEmailAlreadyRegistered
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	newUser := &user.User{
		Email:    email,
		Password: password,
		Name:     name,
	}
	if age != nil {
		newUser.Age = *age
	}

	if err := s.db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return s.issueTokens(*newUser)
}

func (s *Service) Login(email string, password string) (*AuthTokens, error) {
	var dbUser user.User
	if err := s.db.Where("email = ?", email).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if dbUser.Password != password {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(dbUser)
}

func (s *Service) Refresh(refreshToken string) (*AuthTokens, error) {
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	tokenHash := hashToken(refreshToken)
	var storedToken RefreshToken
	if err := s.db.Where("token_hash = ?", tokenHash).First(&storedToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}

		return nil, err
	}

	if time.Now().UTC().After(storedToken.ExpiresAt.UTC()) {
		if err := s.db.Delete(&storedToken).Error; err != nil {
			return nil, err
		}

		return nil, ErrInvalidRefreshToken
	}

	subjectUserID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || uint(subjectUserID) != storedToken.UserID {
		return nil, ErrInvalidRefreshToken
	}

	var dbUser user.User
	if err := s.db.First(&dbUser, storedToken.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}

		return nil, err
	}

	if err := s.db.Delete(&storedToken).Error; err != nil {
		return nil, err
	}

	return s.issueTokens(dbUser)
}

func (s *Service) issueTokens(dbUser user.User) (*AuthTokens, error) {
	accessToken, err := s.signAccessToken(dbUser)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.storeRefreshToken(dbUser.ID)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) storeRefreshToken(userID uint) (string, error) {
	if err := s.db.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error; err != nil {
		return "", err
	}

	var dbUser user.User
	if err := s.db.First(&dbUser, userID).Error; err != nil {
		return "", err
	}

	token, expiresAt, err := s.signRefreshToken(dbUser)
	if err != nil {
		return "", err
	}

	refreshToken := &RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(token),
		ExpiresAt: expiresAt,
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return "", err
	}

	return token, nil
}
