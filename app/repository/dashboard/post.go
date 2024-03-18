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

type PostRepository interface {
	Index(limit int, offset uint, search string, sort_by string, sort string) ([]model.Post, int, error)
	Show(UUID string) (model.PostShow, error)
	Store(model *model.StorePost) (model.Post, error)
	Update(UUID string, request *model.UpdatePost) (model.Post, error)
	Destroy(UUID string) (model.Post, error)
	GetSlug(Title string, UUID *string) string
}

type PostRepo struct {
	db *database.DB
}

func (repo *PostRepo) Index(limit int, offset uint, search string, sort_by string, sort string) ([]model.Post, int, error) {
	_select := `
	posts.uuid,
    user_uuid,
    tag_uuid,
    title,
    thumbnail,
    content,
    keyword,
    posts.slug,
    posts.is_active,
    is_highlight,
    posts.created_at,
    posts.updated_at,
    IFNULL(
        JSON_OBJECT(
            'uuid', users.uuid,
            'name', users.name,
            'username', users.username,
            'email', users.email,
            'created_at', users.created_at,
            'updated_at', users.updated_at
        ),
        NULL
    ) AS user,
    IFNULL(
        JSON_OBJECT(
            'uuid', tags.uuid,
            'name', tags.name,
            'slug', tags.slug,
            'is_active', tags.is_active,
            'created_at', tags.created_at,
            'updated_at', tags.updated_at
        ),
        NULL
    ) AS tag
	`
	_conditions := database.Search([]string{"title", "content", "users.name", "tags.name"}, search, "posts.deleted_at")
	_order := ""
	if sort_by == "id" {
		_order = database.OrderBy("posts.id", sort)
	} else {
		_order = database.OrderBy(sort_by, sort)
	}

	_limit := database.Limit(limit, offset)

	count_query := fmt.Sprintf(`SELECT count(*) FROM posts LEFT JOIN users ON users.uuid = posts.user_uuid LEFT JOIN tags ON tags.uuid = posts.tag_uuid %s`, _conditions)
	var count int
	_ = repo.db.QueryRow(count_query).Scan(&count)

	query := fmt.Sprintf(`SELECT %s FROM posts LEFT JOIN users ON users.uuid = posts.user_uuid LEFT JOIN tags ON tags.uuid = posts.tag_uuid %s %s %s`, _select, _conditions, _order, _limit)

	rows, err := repo.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	items := []model.Post{}
	for rows.Next() {
		var i model.Post
		err := rows.Scan(
			&i.UUID,
			&i.TagUUID,
			&i.UserUUID,
			&i.Title,
			&i.Thumbnail,
			&i.Content,
			&i.Keyword,
			&i.Slug,
			&i.IsActive,
			&i.IsHighlight,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.User,
			&i.Tag,
		)
		if err != nil {
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

func (repo *PostRepo) Show(UUID string) (model.PostShow, error) {
	var post model.PostShow
	query := `
	SELECT 
	posts.uuid,
    user_uuid,
    tag_uuid,
    title,
    thumbnail,
    content,
    keyword,
    posts.slug,
    posts.is_active,
    is_highlight,
    posts.created_at,
    posts.updated_at,
    IFNULL(
        JSON_OBJECT(
            'uuid', users.uuid,
            'name', users.name,
            'username', users.username,
            'email', users.email,
            'created_at', users.created_at,
            'updated_at', users.updated_at
        ),
        NULL
    ) AS user,
    IFNULL(
        JSON_OBJECT(
            'uuid', tags.uuid,
            'name', tags.name,
            'slug', tags.slug,
            'is_active', tags.is_active,
            'created_at', tags.created_at,
            'updated_at', tags.updated_at
        ),
        NULL
    ) AS tag
	FROM posts LEFT JOIN users ON users.uuid = posts.user_uuid LEFT JOIN tags ON tags.uuid = posts.tag_uuid
	WHERE posts.uuid = ? AND posts.deleted_at IS NULL LIMIT 1
	`

	err := repo.db.QueryRowContext(context.Background(), query, UUID).Scan(
		&post.UUID,
		&post.TagUUID,
		&post.UserUUID,
		&post.Title,
		&post.Thumbnail,
		&post.Content,
		&post.Keyword,
		&post.Slug,
		&post.IsActive,
		&post.IsHighlight,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.User,
		&post.Tag,
	)

	return post, err
}

func (repo *PostRepo) Store(request *model.StorePost) (model.Post, error) {
	query := `INSERT INTO posts (uuid, tag_uuid, user_uuid, title, thumbnail, content, keyword, slug, is_active, is_highlight, created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?) 
	RETURNING uuid, tag_uuid, user_uuid, title, thumbnail, content, keyword, slug, is_active, is_highlight, created_at`
	var post model.Post
	err := repo.db.QueryRowContext(context.Background(), query, uuid.New(), request.TagUUID, request.UserUUID, request.Title, request.Thumbnail, request.Content, request.Keyword, request.Slug, request.IsActive, request.IsHighlight, time.Now()).Scan(
		&post.UUID,
		&post.TagUUID,
		&post.UserUUID,
		&post.Title,
		&post.Thumbnail,
		&post.Content,
		&post.Keyword,
		&post.Slug,
		&post.IsActive,
		&post.IsHighlight,
		&post.CreatedAt,
	)
	if err != nil {
		return model.Post{}, err
	}
	return post, err
}

func (repo *PostRepo) Update(ID string, request *model.UpdatePost) (model.Post, error) {
	if request.Thumbnail == "" {
		query := `UPDATE posts SET tag_uuid = ?, user_uuid = ?, title = ?, content = ?, keyword = ?, slug = ?, is_active = ?, is_highlight = ?, updated_at = ? WHERE uuid = ?`
		result, err := repo.db.ExecContext(context.Background(), query, request.TagUUID, request.UserUUID, request.Title, request.Content, request.Keyword, request.Slug, request.IsActive, request.IsHighlight, time.Now(), ID)
		if err != nil {
			return model.Post{}, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return model.Post{}, err
		}

		if rowsAffected == 0 {
			return model.Post{}, sql.ErrNoRows
		}

		var post model.Post
		err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, tag_uuid, user_uuid, title, thumbnail, content, keyword, slug, is_active, is_highlight, created_at, updated_at FROM posts WHERE uuid = ?", ID).Scan(
			&post.UUID,
			&post.TagUUID,
			&post.UserUUID,
			&post.Title,
			&post.Thumbnail,
			&post.Content,
			&post.Keyword,
			&post.Slug,
			&post.IsActive,
			&post.IsHighlight,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return model.Post{}, err
		}

		return post, nil
	} else {
		query := `UPDATE posts SET tag_uuid = ?, user_uuid = ?, title = ?, thumbnail = ?, content = ?, keyword = ?, slug = ?, is_active = ?, is_highlight = ?, updated_at = ? WHERE uuid = ?`
		result, err := repo.db.ExecContext(context.Background(), query, request.TagUUID, request.UserUUID, request.Title, request.Thumbnail, request.Content, request.Keyword, request.Slug, request.IsActive, request.IsHighlight, time.Now(), ID)
		if err != nil {
			return model.Post{}, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return model.Post{}, err
		}

		if rowsAffected == 0 {
			return model.Post{}, sql.ErrNoRows
		}

		var post model.Post
		err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, tag_uuid, user_uuid, title, thumbnail, content, keyword, slug, is_active, is_highlight, created_at, updated_at FROM posts WHERE uuid = ?", ID).Scan(
			&post.UUID,
			&post.TagUUID,
			&post.UserUUID,
			&post.Title,
			&post.Thumbnail,
			&post.Content,
			&post.Keyword,
			&post.Slug,
			&post.IsActive,
			&post.IsHighlight,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return model.Post{}, err
		}

		return post, nil
	}
}

func (repo *PostRepo) Destroy(UUID string) (model.Post, error) {
	query := `UPDATE posts SET updated_at = ?, deleted_at = ? WHERE uuid = ? 
	RETURNING uuid, tag_uuid, user_uuid, title, thumbnail, content, keyword, slug, is_active, is_highlight, created_at, updated_at, deleted_at`
	result, err := repo.db.ExecContext(context.Background(), query, time.Now(), time.Now(), UUID)
	if err != nil {
		return model.Post{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Post{}, err
	}

	if rowsAffected == 0 {
		return model.Post{}, sql.ErrNoRows
	}

	var post model.Post
	err = repo.db.QueryRowContext(context.Background(), "SELECT uuid, tag_uuid, user_uuid, title, thumbnail, content, keyword, slug, is_active, is_highlight, created_at, updated_at, deleted_at FROM posts WHERE uuid = ?", UUID).Scan(
		&post.UUID,
		&post.TagUUID,
		&post.UserUUID,
		&post.Title,
		&post.Thumbnail,
		&post.Content,
		&post.Keyword,
		&post.Slug,
		&post.IsActive,
		&post.IsHighlight,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.DeletedAt,
	)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (repo *PostRepo) GetSlug(Title string, UUID *string) string {
	count := 0
	first_slug := slug.Make(Title)
	var slug_check string

	query := ""
	if UUID == nil {
		query = `
			SELECT 
			slug
			FROM posts 
			WHERE posts.slug = ? LIMIT 1
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
			new_slug = slug.Make(Title) + "-" + strconv.Itoa(count)
			next_query := `SELECT slug FROM posts WHERE slug = ? LIMIT 1`
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
			FROM posts 
			WHERE posts.slug = ? AND uuid != ? LIMIT 1
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
			new_slug = slug.Make(Title) + "-" + strconv.Itoa(count)
			next_query := `SELECT slug FROM posts WHERE slug = ? AND uuid != ? LIMIT 1`
			err_again := repo.db.QueryRowContext(context.Background(), next_query, new_slug, UUID).Scan(&slug_check)
			if err_again != nil {
				break
			}
		}
		return new_slug
	}
}

func NewPostRepo(db *database.DB) PostRepository {
	return &PostRepo{db}
}
