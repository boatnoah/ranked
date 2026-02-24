package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"crypto/sha256"
	"encoding/hex"
)

const expiryOffset = 168

type SessionStore struct {
	db *sql.DB
}

type Session struct {
	Token string
}

func (s *SessionStore) CreateSession(ctx context.Context, uuid uuid.UUID) (*Session, error) {

	query := `
		INSERT INTO sessions (user_uuid, session_token, expires_at)
		VALUES ($1, $2, $3)
	`

	var session Session

	sessionToken, err := newToken()
	if err != nil {
		return nil, err
	}

	session.Token = sessionToken

	expiryDate := time.Now().Add(expiryOffset * time.Hour)

	_, err = s.db.ExecContext(ctx, query, uuid, hashTokenSHA256(sessionToken), expiryDate)
	if err != nil {
		return nil, err
	}

	return &session, nil

}

func (s *SessionStore) GetIDbyToken(ctx context.Context, token string) (string, error) {
	query := `		
		SELECT user_uuid FROM sessions
		WHERE session_token = $1
	`

	var userID string

	hashedToken := hashTokenSHA256(token)
	err := s.db.QueryRowContext(ctx, query, hashedToken).Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, nil

}
func newToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func hashTokenSHA256(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
