package apiserver

import (
	"async_api/config"
	"fmt"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var signingMethod = jwt.SigningMethodHS256

type JwtManager struct {
	config *config.Config
}

func NewJwtManager(config *config.Config) *JwtManager {
	return &JwtManager{
		config: config,
	}
}

type TokenPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}

type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// Parse takes in a token as string, parses it and returns jwt.Token object
func (j *JwtManager) Parse(token string) (*jwt.Token, error) {
	parser := jwt.NewParser()
	jwtToken, err := parser.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.config.JwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %w", err)
	}
	return jwtToken, nil
}

func (j *JwtManager) IsAccessToken(token *jwt.Token) bool {
	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	if tokenType, ok := jwtClaims["token_type"]; ok {
		return tokenType == "access"
	}
	return false
}

// GenerateTokenPair generates a new struct of TokenPair
func (j *JwtManager) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	now := time.Now()
	issuer := "http://" + net.JoinHostPort(j.config.ApiServerHost, j.config.ApiServerPort)
	// jwtAccessTokenLifetime, err := strconv.Atoi(j.config.JwtAccessTokenLifetime)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to convert access token lifetime %w", err)
	// }
	jwtAccessToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	key := []byte(j.config.JwtSecret)
	signedAccessToken, err := jwtAccessToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token %w", err)
	}

	accessToken, err := j.Parse(signedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token %w", err)
	}

	// jwtRefreshTokenLifetime, err := strconv.Atoi(j.config.JwtRefreshTokenLifetime)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to convert refresh token lifetime %w", err)
	// }
	jwtRefreshToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 5)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	key = []byte(j.config.JwtSecret)
	signedRefreshToken, err := jwtRefreshToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token %w", err)
	}

	refreshToken, err := j.Parse(signedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
