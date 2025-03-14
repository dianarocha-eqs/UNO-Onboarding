package repository

import (
	config "api/configs"
	auth_domain "api/internal/auth/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for token's data operations
type AuthRepository interface {
	// Stores a new token in the database
	StoreToken(ctx context.Context, auth *auth_domain.AuthToken) error
	// Gets a specific token and returns it
	GetToken(ctx context.Context, tokenStr string) (*auth_domain.AuthToken, error)
	// Sets token to false (invalid)
	InvalidateToken(ctx context.Context, tokenStr string) error
	// Stores new token for password recovery
	StoreTokenToPasswordRecovery(ctx context.Context, userID uuid.UUID, token string, expirationTime time.Time) error
}

// database/sql to interact with the token's database
type AuthRepositoryImpl struct {
	DB *sql.DB
}

// Connects with the database
func NewAuthRepository() (AuthRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &AuthRepositoryImpl{DB: db}, nil
}

func (r *AuthRepositoryImpl) StoreToken(ctx context.Context, auth *auth_domain.AuthToken) error {
	query := `
		BEGIN
			-- Delete old tokens before inserting a new one
			DELETE FROM users_tokens WHERE userUuid = @userUuid;

			-- Insert new token
			INSERT INTO users_tokens (uuid, userUuid, token, is_valid, created_at, expired_at)
			VALUES (@uuid, @userUuid, @token, 1, @createdAt, @expiredAt);
		END;
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("userUuid", auth.UserID),
		sql.Named("token", auth.Token),
		sql.Named("createdAt", auth.CreatedAt),
		sql.Named("expiredAt", auth.ExpiredAt),
		sql.Named("uuid", auth.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	return nil
}

func (r *AuthRepositoryImpl) GetToken(ctx context.Context, tokenStr string) (*auth_domain.AuthToken, error) {
	query := `
		SELECT userUuid, token, is_valid, created_at, expired_at
		FROM users_tokens
		WHERE token = @token
	`

	var authToken auth_domain.AuthToken
	row := r.DB.QueryRowContext(ctx, query, sql.Named("token", tokenStr))

	// Scan the result into the authToken struct
	err := row.Scan(&authToken.UserID, &authToken.Token, &authToken.IsValid, &authToken.CreatedAt, &authToken.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid or expired token")
		}
		return nil, fmt.Errorf("failed to retrieve token: %v", err)
	}

	return &authToken, nil
}

func (r *AuthRepositoryImpl) InvalidateToken(ctx context.Context, tokenStr string) error {
	query := `
		UPDATE users_tokens
		SET is_valid = 0
		WHERE token = @token
	`

	_, err := r.DB.ExecContext(ctx, query, sql.Named("token", tokenStr))
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}

	return nil
}

func (r *AuthRepositoryImpl) StoreTokenToPasswordRecovery(ctx context.Context, userID uuid.UUID, token string, expirationTime time.Time) error {
	// Insert only the token and expiration time into the Users_Tokens table based on userUuid
	query := `
		UPDATE users_tokens 
		SET password_recovery_token = @token, 
			password_recovery_expiration = @expiredAt
		WHERE userUuid = @userUuid;
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("userUuid", userID),
		sql.Named("token", token),
		sql.Named("expiredAt", expirationTime),
	)

	if err != nil {
		return fmt.Errorf("failed to store password recovery token: %v", err)
	}

	return nil
}
