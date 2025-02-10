package repository

import (
	config "api/configs"
	"api/internal/users/domain"
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/tentone/mssql-uuid"

	_ "github.com/denisenkom/go-mssqldb" // Import SQL Server driver
)

// Interface for user's data operations
type UserRepository interface {
	// Stores a new user in the database
	CreateUser(ctx context.Context, user *domain.User) error
	// Updates the details of an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
	// Get the users info
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	// Get any users info
	ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error)
	// Get user by email and password
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, error)
}

// Performs user's data operations using GORM to interact with the database
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
		INSERT INTO Users (id, name, email, password, picture, phone, role)
		VALUES (NEWID(), @name, @email, @password, @picture, @phone, @role)
	`

	_, err := r.DB.ExecContext(ctx, query,
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
		UPDATE Users
		SET 
			name = COALESCE(NULLIF(@name, ''), name),
			email = COALESCE(NULLIF(@email, ''), email),
			phone = COALESCE(NULLIF(@phone, ''), phone),
			picture = COALESCE(NULLIF(@picture, ''), picture),
			password = COALESCE(NULLIF(@password, ''), password)
		WHERE id = @id
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("name", user.Name),
		sql.Named("email", user.Email),
		sql.Named("phone", user.Phone),
		sql.Named("picture", user.Picture),
		sql.Named("password", user.Password),
		sql.Named("id", user.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, name, email, picture, phone, role
		FROM Users
		WHERE id = @id
	`

	row := r.DB.QueryRowContext(ctx, query, sql.Named("id", userID))

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

func (r *UserRepositoryImpl) ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error) {
	query := `
		SELECT id, name, picture
		FROM Users
		WHERE name LIKE '%' + @search + '%' OR email LIKE '%' + @search + '%'
		ORDER BY CASE WHEN @sortDirection = 1 THEN name END ASC,
				 CASE WHEN @sortDirection = -1 THEN name END DESC
	`

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

func (r *UserRepositoryImpl) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, error) {
	query := `
		SELECT id, name, email, picture, phone, role
		FROM Users
		WHERE email = @email AND password = @password
	`

	row := r.DB.QueryRowContext(ctx, query, sql.Named("email", email), sql.Named("password", password))

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Picture, &user.Phone, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return &user, nil
}
