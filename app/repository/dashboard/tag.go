package dashboard

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type TagRepository interface {
	Index(limit int, offset uint, search string, sort_by string, sort string) ([]model.Tag, int, error)
	Show(UUID string) (model.TagShow, error)
	Store(model *model.StoreTag) (model.Tag, error)
	Update(UUID string, request *model.UpdateTag) (model.Tag, error)
	Destroy(UUID string) (model.Tag, error)
	GetSlug(Name string, UUID *string) string
}

type TagRepo struct {
	db *database.DB
}

func (repo *TagRepo) Index(limit int, offset uint, search string, sort_by string, sort string) ([]model.Tag, int, error) {
	_select := "uuid, name, slug, is_active, created_at, updated_at, deleted_at"
	_conditions := database.Search([]string{"name"}, search, "tags.deleted_at")
	_order := database.OrderBy(sort_by, sort)
	_limit := database.Limit(limit, offset)

	count_query := fmt.Sprintf(`SELECT count(*) FROM tags %s`, _conditions)
	var count int
	_ = repo.db.QueryRow(count_query).Scan(&count)

	query := fmt.Sprintf(`SELECT %s FROM tags %s %s %s`, _select, _conditions, _order, _limit)

	rows, err := repo.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	items := []model.Tag{}
	for rows.Next() {
		var i model.Tag
		if err := rows.Scan(
			&i.UUID,
			&i.Name,
			&i.Slug,
			&i.IsActive,
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

func (repo *TagRepo) Show(UUID string) (model.TagShow, error) {
	var Tag model.TagShow
	query := "SELECT uuid, name, slug, is_active, created_at, updated_at, deleted_at FROM tags WHERE uuid = ? AND tags.deleted_at IS NULL LIMIT 1"
	err := repo.db.QueryRowContext(context.Background(), query, UUID).Scan(
		&Tag.UUID,
		&Tag.Name,
		&Tag.Slug,
		&Tag.IsActive,
		&Tag.CreatedAt,
		&Tag.UpdatedAt,
		&Tag.DeletedAt,
	)
	if err != nil {
		return model.TagShow{}, err
	}
	return Tag, err
}

func (repo *TagRepo) Store(request *model.StoreTag) (model.Tag, error) {
	query := `INSERT INTO tags (uuid, name, slug, is_active, created_at) VALUES(?,?,?,?,?) 
	RETURNING uuid, name, slug, is_active, created_at`
	var Tag model.Tag
	err := repo.db.QueryRowContext(context.Background(), query, uuid.New(), request.Name, request.Slug, request.IsActive, time.Now()).Scan(
		&Tag.UUID,
		&Tag.Name,
		&Tag.Slug,
		&Tag.IsActive,
		&Tag.CreatedAt,
	)
	if err != nil {
		return model.Tag{}, err
	}
	return Tag, err
}

func (repo *TagRepo) Update(UUID string, request *model.UpdateTag) (model.Tag, error) {
	query := `UPDATE tags SET name = ?, slug = ?, is_active = ?, updated_at = ? WHERE uuid = ?`
	result, err := repo.db.ExecContext(context.Background(), query, request.Name, request.Slug, request.IsActive, time.Now(), UUID)
	if err != nil {
		return model.Tag{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Tag{}, err
	}

	if rowsAffected == 0 {
		return model.Tag{}, sql.ErrNoRows
	}

	var tag model.Tag
	err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, name, slug, is_active, created_at, updated_at FROM tags WHERE uuid = ?", UUID).Scan(
		&tag.UUID,
		&tag.Name,
		&tag.Slug,
		&tag.IsActive,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		return model.Tag{}, err
	}

	return tag, nil
}

func (repo *TagRepo) Destroy(UUID string) (model.Tag, error) {
	query := `UPDATE tags SET updated_at = ?, deleted_at = ? WHERE uuid = ?`
	result, err := repo.db.ExecContext(context.Background(), query, time.Now(), time.Now(), UUID)
	if err != nil {
		return model.Tag{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Tag{}, err
	}

	if rowsAffected == 0 {
		return model.Tag{}, sql.ErrNoRows
	}

	var tag model.Tag
	err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, name, slug, is_active, created_at, updated_at FROM tags WHERE uuid = ?", UUID).Scan(
		&tag.UUID,
		&tag.Name,
		&tag.Slug,
		&tag.IsActive,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		return model.Tag{}, err
	}

	return tag, nil
}

func (repo *TagRepo) GetSlug(Name string, UUID *string) string {
	count := 0
	first_slug := slug.Make(Name)
	var slug_check string

	query := ""
	if UUID == nil {
		query = `
			SELECT 
			slug
			FROM tags 
			WHERE tags.slug = ? LIMIT 1
		`

		err := repo.db.QueryRowContext(context.Background(), query, first_slug).Scan(&slug_check)

		new_slug := ""

		if err != nil {
			if err == sql.ErrNoRows {
				return first_slug
			}
		}

		for true {
			count++
			new_slug = slug.Make(Name) + "-" + strconv.Itoa(count)
			next_query := `SELECT slug FROM tags WHERE slug = ? LIMIT 1`
			err_again := repo.db.QueryRowContext(context.Background(), next_query, new_slug).Scan(&slug_check)
			if err_again != nil {
				break
			}
		}
		return new_slug
	} else {
		query = `
			SELECT 
			slug
			FROM tags 
			WHERE tags.slug = ? AND uuid != ? LIMIT 1
		`

		err := repo.db.QueryRowContext(context.Background(), query, first_slug, UUID).Scan(&slug_check)

		new_slug := ""

		if err != nil {
			if err == sql.ErrNoRows {
				return first_slug
			}
		}

		for true {
			count++
			new_slug = slug.Make(Name) + "-" + strconv.Itoa(count)
			next_query := `SELECT slug FROM tags WHERE slug = $1 AND uuid != $2 LIMIT 1`
			err_again := repo.db.QueryRowContext(context.Background(), next_query, new_slug, UUID).Scan(&slug_check)
			if err_again != nil {
				break
			}
		}
		return new_slug
	}
}

func NewTagRepo(db *database.DB) TagRepository {
	return &TagRepo{db}
}
