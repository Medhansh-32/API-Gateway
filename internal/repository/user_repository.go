package repository

import (
	"database/sql"

	"github.com/medhansh-32/api-gateway/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepository *UserRepository) GetUserByUserId(userId int64) (*models.User, error) {
	
	query := `
		SELECT
			id,
			username,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = ?;
		`

	var user models.User

	err := userRepository.db.QueryRow(query, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
