package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenStore struct {
	db *sqlx.DB
}

func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

type RefreshToken struct {
	UserID      uuid.UUID `db:"user_id"`
	HashedToken string    `db:"hashed_token"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiresAt   time.Time `db:"expires_at"`
}

func (s *RefreshTokenStore) getBase64HashFromToken(token *jwt.Token) (string, error) {
	// hashedToken, err := bcrypt.GenerateFromPassword([]byte(token.Raw), bcrypt.DefaultCost)
	// if err != nil {
	// 	return "", fmt.Errorf("bcrypt.GenerateFromPassward: %w", err)
	// }

	h := sha256.New()
	h.Write([]byte(token.Raw))
	hashedBytes := h.Sum((nil))
	base64TokenHash := base64.StdEncoding.EncodeToString(hashedBytes)
	return base64TokenHash, nil
}

// Create inserts a new record into refresh_tokens table
func (s *RefreshTokenStore) Create(ctx context.Context, userID uuid.UUID, token *jwt.Token) (*RefreshToken, error) {
	const stmt = `INSERT INTO refresh_tokens (user_id, hashed_token, expires_at) VALUES ($1, $2, $3) RETURNING *;`
	base64TokenHash, err := s.getBase64HashFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get base64 encoded token hash: %w", err)
	}

	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to extract expiration time: %w", err)
	}

	var refreshToken RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, stmt, userID, base64TokenHash, expiresAt.Time); err != nil {
		return nil, fmt.Errorf("failed to create refresh token record: %w", err)
	}

	return &refreshToken, nil
}

// ByPrimaryKey extracts the refresh_token record by userID and refresh_token
func (s *RefreshTokenStore) ByPrimaryKey(ctx context.Context, userID uuid.UUID, token *jwt.Token) (*RefreshToken, error) {
	const stmt = `SELECT * FROM refresh_tokens WHERE user_id = $1 AND hashed_token = $2;`
	base64TokenHash, err := s.getBase64HashFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get base64 encoded token hash: %w", err)
	}

	var refreshToken RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, stmt, userID, base64TokenHash); err != nil {
		return nil, fmt.Errorf("failed to fetch refreshed_token %s record for user %s: %w", base64TokenHash, userID, err)
	}

	return &refreshToken, nil
}

func (s *RefreshTokenStore) DeleteUserTokens(ctx context.Context, userID uuid.UUID) (sql.Result, error) {
	const stmt = `DELETE FROM refresh_tokens WHERE user_id = $1;`
	result, err := s.db.ExecContext(ctx, stmt, userID)
	if err != nil {
		return result, fmt.Errorf("failed to delete refresh_tokens record: %w", err)
	}
	return result, nil
}
