package repository

import (
	config "api/configs"
	"api/internal/users/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // Import SQL Server driver
	uuid "github.com/tentone/mssql-uuid"
)

// Interface for user's data operations
type UserRepository interface {
	// Stores a new user in the database
	CreateUser(ctx context.Context, user *domain.User) error
	// Updates the details of an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
	// Get any users info
	ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error)
	// Get user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	// Authenticate user through email and password
	AuthenticateUser(ctx context.Context, email, password string) error
	// Checks user's role and uuid from token
	GetRoutesAuthorization(ctx context.Context, tokenStr string, getRole *bool, getUserID *uuid.UUID) error
	// Updates user's password and deletes token for password reset
	ResetPassword(ctx context.Context, token string, password string) error
	// Get user id by token
	GetUserByToken(ctx context.Context, tokenStr string) (uuid.UUID, error)
}

// Performs user's data operations using database/sql to interact with the database
type UserRepositoryImpl struct {
	DB *sql.DB
}

// Connects with the database
func NewUserRepository() (UserRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &UserRepositoryImpl{DB: db}, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (uuid, name, email, password, picture, phone, role)
		VALUES (@uuid, @name, @email, @password, @picture, @phone, @role)
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("uuid", user.ID),
		sql.Named("name", user.Name),
		sql.Named("email", user.Email),
		sql.Named("password", user.Password),
		sql.Named("picture", user.Picture),
		sql.Named("phone", user.Phone),
		sql.Named("role", user.Role),
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET 
			name = COALESCE(NULLIF(@name, ''), name),
			email = COALESCE(NULLIF(@email, ''), email),
			phone = COALESCE(NULLIF(@phone, ''), phone),
			picture = @picture,
			password = COALESCE(NULLIF(@password, ''), password)
		WHERE uuid = @uuid
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("name", user.Name),
		sql.Named("email", user.Email),
		sql.Named("phone", user.Phone),
		sql.Named("picture", user.Picture),
		sql.Named("password", user.Password),
		sql.Named("uuid", user.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error) {
	query := `
		SELECT uuid, name, picture
		FROM users
		WHERE name LIKE '%' + @search + '%' OR email LIKE '%' + @search + '%'
		ORDER BY CASE WHEN @sortDirection = 1 THEN name END ASC,
				 CASE WHEN @sortDirection = -1 THEN name END DESC
	`

	// Execute the query with named parameters
	rows, err := r.DB.QueryContext(ctx, query,
		sql.Named("search", search),
		sql.Named("sortDirection", sortDirection),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Picture); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %v", err)
	}

	return users, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	query := "SELECT uuid, name, email, picture, phone, role FROM users WHERE email = @email"
	row := r.DB.QueryRowContext(ctx, query, sql.Named("email", email))

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Picture, &user.Phone, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) AuthenticateUser(ctx context.Context, email, password string) error {

	query := "SELECT 1 FROM users WHERE email = @email AND password = @password"
	row := r.DB.QueryRowContext(ctx, query, sql.Named("email", email), sql.Named("password", password))

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid credentials")
		}
		return fmt.Errorf("authentication error: %v", err)
	}

	return nil
}

func (r *UserRepositoryImpl) GetRoutesAuthorization(ctx context.Context, tokenStr string, getRole *bool, getUserID *uuid.UUID) error {
	query := `
		SELECT users.role, users.uuid
		FROM users
		INNER JOIN users_tokens
		ON users_tokens.userUuid = users.uuid
		WHERE users_tokens.token = @token
	`

	var role bool
	var userid uuid.UUID
	row := r.DB.QueryRowContext(ctx, query, sql.Named("token", tokenStr))
	// Scan the result for role and uuid
	err := row.Scan(&role, &userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("token not found")
		}
		return fmt.Errorf("failed to retrieve role and user id: %v", err)
	}

	// Assign only if the caller wants these values
	if getRole != nil {
		*getRole = role
	}
	if getUserID != nil {
		*getUserID = userid
	}

	return nil
}

func (r *UserRepositoryImpl) ResetPassword(ctx context.Context, token string, password string) error {

	// Get a Tx for making transaction requests.
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// Gets user id and expiration date from token
	query := `
		SELECT userUuid, password_recovery_expiration
		FROM users_tokens
		WHERE password_recovery_token = @password_recovery_token 
	`

	var userUuid uuid.UUID
	var passwordRecoveryExpiration time.Time
	row := r.DB.QueryRowContext(ctx, query, sql.Named("password_recovery_token", token))
	err = row.Scan(&userUuid, &passwordRecoveryExpiration)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid or expired token")
		}
		return fmt.Errorf("failed to retrieve user and token data: %v", err)
	}

	// Updates user's password from id retrived
	query = `
		UPDATE users
		SET password = @password
		WHERE uuid = @uuid
	`
	_, err = tx.ExecContext(ctx, query,
		sql.Named("password", password),
		sql.Named("uuid", userUuid),
	)
	if err != nil {
		return fmt.Errorf("failed to update user's password: %v", err)
	}

	// Deletes token for this action
	query = `UPDATE users_tokens SET password_recovery_token = NULL, password_recovery_expiration = NULL WHERE password_recovery_token = @token`
	_, err = tx.ExecContext(ctx, query, sql.Named("token", token))
	if err != nil {
		return fmt.Errorf("failed to delete password recovery token: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *UserRepositoryImpl) GetUserByToken(ctx context.Context, tokenStr string) (uuid.UUID, error) {
	query := `
		SELECT users.uuid
		FROM users
		INNER JOIN users_tokens
		ON users_tokens.userUuid = users.uuid
		WHERE users_tokens.token = @token
	`

	var userUuid uuid.UUID
	row := r.DB.QueryRowContext(ctx, query, sql.Named("token", tokenStr))
	err := row.Scan(&userUuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.NilUUID, fmt.Errorf("token not found")
		}
		return uuid.NilUUID, fmt.Errorf("failed to retrieve role and user id: %v", err)
	}

	return userUuid, nil
}
