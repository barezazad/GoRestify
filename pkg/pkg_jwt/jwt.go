package pkg_jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"GoRestify/pkg/pkg_types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// TokenTypeAccess marks an access token.
	TokenTypeAccess = "access"
	// TokenTypeRefresh marks a refresh token.
	TokenTypeRefresh = "refresh"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	accessTTL  = time.Hour
	refreshTTL = 7 * 24 * time.Hour
)

// Init loads RSA keys from PEM files and sets token TTLs.
func Init(privateKeyPath, publicKeyPath string, accessMinutes, refreshHours int) error {
	privPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("read JWT private key: %w", err)
	}
	pubPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("read JWT public key: %w", err)
	}

	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privPEM)
	if err != nil {
		return fmt.Errorf("parse JWT private key: %w", err)
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubPEM)
	if err != nil {
		return fmt.Errorf("parse JWT public key: %w", err)
	}

	privateKey = priv
	publicKey = pub

	if accessMinutes > 0 {
		accessTTL = time.Duration(accessMinutes) * time.Minute
	}
	if refreshHours > 0 {
		refreshTTL = time.Duration(refreshHours) * time.Hour
	}

	return nil
}

// TokenPair holds access + refresh tokens.
type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"` // access token lifetime in seconds
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	AccessJTI        string `json:"-"`
	RefreshJTI       string `json:"-"`
}

// GenerateTokenPair creates RS256 access (short) and refresh (long) tokens.
func GenerateTokenPair(userID uint, username string) (TokenPair, error) {
	if privateKey == nil {
		return TokenPair{}, errors.New("JWT private key is not initialized")
	}

	now := time.Now()
	accessJTI := uuid.NewString()
	refreshJTI := uuid.NewString()

	accessClaims := &pkg_types.JWTClaims{
		UserID:    userID,
		Username:  username,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims).SignedString(privateKey)
	if err != nil {
		return TokenPair{}, err
	}

	refreshClaims := &pkg_types.JWTClaims{
		UserID:    userID,
		Username:  username,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTTL)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(privateKey)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        int64(accessTTL.Seconds()),
		RefreshExpiresIn: int64(refreshTTL.Seconds()),
		AccessJTI:        accessJTI,
		RefreshJTI:       refreshJTI,
	}, nil
}

// ParseAndValidate parses a token with the public key and requires RS256.
func ParseAndValidate(tokenStr string) (*pkg_types.JWTClaims, error) {
	if publicKey == nil {
		return nil, errors.New("JWT public key is not initialized")
	}

	claims := &pkg_types.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is not valid")
	}
	return claims, nil
}

// AccessTTL returns configured access token lifetime.
func AccessTTL() time.Duration { return accessTTL }

// RefreshTTL returns configured refresh token lifetime.
func RefreshTTL() time.Duration { return refreshTTL }

// RefreshRedisKey builds the Redis key for a refresh jti.
func RefreshRedisKey(jti string) string {
	return "jwt-refresh-" + jti
}
