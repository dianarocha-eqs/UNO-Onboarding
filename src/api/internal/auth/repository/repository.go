package repository

import (
	config "api/configs"
	auth_domain "api/internal/auth/domain"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // Import SQL Server driver
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
	StoreTokenToPasswordRecovery(ctx context.Context, auth *auth_domain.AuthToken) error
	// Get user's password reset token
	GetUserByTokenToResetPassword(ctx context.Context, token string) (uuid.UUID, error)
	// Delete token
	DeteteToken(ctx context.Context, token string) error
}

// GORM to interact with the token's database
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
			DELETE FROM user_tokens WHERE user_id = @user_id;

			-- Insert new token
			INSERT INTO user_tokens (user_id, token, is_valid)
			VALUES (@user_id, @token, 1);
		END;
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("user_id", auth.UserID),
		sql.Named("token", auth.Token),
	)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	return nil
}

func (r *AuthRepositoryImpl) GetToken(ctx context.Context, tokenStr string) (*auth_domain.AuthToken, error) {
	query := `
		SELECT user_id, token, is_valid
		FROM user_tokens
		WHERE token = @token
	`

	var authToken auth_domain.AuthToken
	row := r.DB.QueryRowContext(ctx, query, sql.Named("token", tokenStr))

	// Scan the result into the authToken struct
	err := row.Scan(&authToken.UserID, &authToken.Token, &authToken.IsValid)
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
		UPDATE user_tokens
		SET is_valid = 0
		WHERE token = @token
	`

	_, err := r.DB.ExecContext(ctx, query, sql.Named("token", tokenStr))
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}

	return nil
}

func (r *AuthRepositoryImpl) StoreTokenToPasswordRecovery(ctx context.Context, auth *auth_domain.AuthToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, is_valid)
		VALUES (@user_id, @token, 1);
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("user_id", auth.UserID),
		sql.Named("token", auth.Token),
	)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	return nil
}

func (r *AuthRepositoryImpl) GetUserByTokenToResetPassword(ctx context.Context, token string) (uuid.UUID, error) {
	query := `
		SELECT user_id
		FROM password_reset_tokens
		WHERE token = @token
	`

	var userID uuid.UUID
	row := r.DB.QueryRowContext(ctx, query, sql.Named("token", token))
	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.NilUUID, fmt.Errorf("invalid or expired token")
		}
		return uuid.NilUUID, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return userID, nil
}

func (r *AuthRepositoryImpl) DeteteToken(ctx context.Context, token string) error {

	deleteQuery := `
		DELETE FROM password_reset_tokens
		WHERE token = @token
	`
	_, _ = r.DB.ExecContext(ctx, deleteQuery, sql.Named("token", token))

	return nil
}
