package dashboard

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/google/uuid"
)

type UserRepository interface {
	Index(limit int, offset uint, search string, sort_by string, sort string) ([]model.User, int, error)
	Show(UUID string) (model.UserShow, error)
	Store(model *model.StoreUser) (model.User, error)
	Update(UUID string, request *model.UpdateUser) (model.User, error)
	Destroy(UUID string) (model.User, error)
}

type UserRepo struct {
	db *database.DB
}

func (repo *UserRepo) Index(limit int, offset uint, search string, sort_by string, sort string) ([]model.User, int, error) {
	_select := "users.uuid, users.name, email, username, role_uuid, roles.name as role_name, users.created_at, users.updated_at, users.deleted_at"
	_conditions := database.Search([]string{"users.name", "email", "username"}, search, "users.deleted_at")
	sort_by = "users.id"
	_order := database.OrderBy(sort_by, sort)
	_limit := database.Limit(limit, offset)

	count_query := fmt.Sprintf(`SELECT count(*) FROM users %s`, _conditions)
	var count int
	_ = repo.db.QueryRow(count_query).Scan(&count)

	query := fmt.Sprintf(`SELECT %s FROM users LEFT JOIN roles ON roles.uuid = users.role_uuid %s %s %s`, _select, _conditions, _order, _limit)

	rows, err := repo.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	items := []model.User{}
	for rows.Next() {
		var i model.User
		if err := rows.Scan(
			&i.UUID,
			&i.Name,
			&i.Email,
			&i.Username,
			&i.RoleUUID,
			&i.RoleName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

func (repo *UserRepo) Show(ID string) (model.UserShow, error) {
	var user model.UserShow
	query := "SELECT users.uuid, users.name, email, username, role_uuid, roles.name as role_name, users.created_at, users.updated_at, users.deleted_at FROM users LEFT JOIN roles ON roles.uuid = users.role_uuid WHERE users.uuid = ? AND users.deleted_at IS NULL LIMIT 1"
	err := repo.db.QueryRowContext(context.Background(), query, ID).Scan(
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.RoleUUID,
		&user.RoleName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return model.UserShow{}, err
	}
	return user, err
}

func (repo *UserRepo) Store(request *model.StoreUser) (model.User, error) {
	query := `INSERT INTO users (uuid, name, username, email, role_uuid, password, created_at) VALUES(?, ?, ?, ?, ?, ?, ?) 
	RETURNING uuid, name, username, email, role_uuid, created_at`
	var user model.User
	err := repo.db.QueryRowContext(context.Background(), query, uuid.New(), request.Name, request.Username, request.Email, request.RoleUUID, request.Password, time.Now()).Scan(
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.RoleUUID,
		&user.CreatedAt,
	)
	if err != nil {
		return model.User{}, err
	}
	return user, err
}

func (repo *UserRepo) Update(ID string, request *model.UpdateUser) (model.User, error) {
	if request.Password == "" {
		query := `UPDATE users SET name = ?, username = ?, email = ?, role_uuid = ?, updated_at = ? WHERE uuid = ?`
		result, err := repo.db.ExecContext(context.Background(), query, request.Name, request.Username, request.Email, request.RoleUUID, time.Now(), ID)
		if err != nil {
			return model.User{}, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return model.User{}, err
		}

		if rowsAffected == 0 {
			return model.User{}, errors.New("no rows updated")
		}

		var user model.User
		err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, name, username, email, role_uuid, created_at, updated_at FROM users WHERE uuid = ?", ID).Scan(
			&user.UUID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.RoleUUID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return model.User{}, err
		}

		return user, nil
	} else {
		query := `UPDATE users SET name = ?, username = ?, email = ?, role_uuid = ?, password = ?, updated_at = ? WHERE uuid = ?`
		result, err := repo.db.ExecContext(context.Background(), query, request.Name, request.Username, request.Email, request.RoleUUID, request.Password, time.Now(), ID)
		if err != nil {
			return model.User{}, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return model.User{}, err
		}

		if rowsAffected == 0 {
			return model.User{}, errors.New("no rows updated")
		}

		var user model.User
		err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, name, username, email, role_uuid, created_at, updated_at FROM users WHERE uuid = ?", ID).Scan(
			&user.UUID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.RoleUUID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return model.User{}, err
		}

		return user, nil
	}
}

func (repo *UserRepo) Destroy(ID string) (model.User, error) {
	query := `UPDATE users SET updated_at = ?, deleted_at = ? WHERE uuid = ?`
	result, err := repo.db.ExecContext(context.Background(), query, time.Now(), time.Now(), ID)
	if err != nil {
		return model.User{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.User{}, err
	}

	if rowsAffected == 0 {
		return model.User{}, errors.New("no rows updated")
	}

	var user model.User
	err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, name, username, email, role_uuid, created_at, updated_at, deleted_at FROM users WHERE uuid = ?", ID).Scan(
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.RoleUUID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func NewUserRepo(db *database.DB) UserRepository {
	return &UserRepo{db}
}
