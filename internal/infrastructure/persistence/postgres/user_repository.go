package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RodriguezMjs/tasks-tracking/internal/domain/entities"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error al consultar usuario: %w", err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	roles, err := r.getUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener roles: %w", err)
	}

	user.Roles = roles
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	user := &entities.User{}

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error al consultar usuario: %w", err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	roles, err := r.getUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener roles: %w", err)
	}

	user.Roles = roles
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Name, user.Email, user.PasswordHash,
	)

	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}

func (r *UserRepository) getUserRoles(ctx context.Context, userID string) ([]entities.Role, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT r.id, r.name FROM roles r
		 INNER JOIN user_roles ur ON r.id = ur.role_id
		 WHERE ur.user_id = $1`,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("error al consultar roles: %w", err)
	}
	defer rows.Close()

	roles := []entities.Role{}
	for rows.Next() {
		var role entities.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, fmt.Errorf("error al mapear roles: %w", err)
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error estructurando roles: %w", err)
	}

	return roles, nil
}
