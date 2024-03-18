package public

import (
	"context"
	"fmt"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/public"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
)

type PostRepository interface {
	Index(limit int, offset uint, search string, sort_by string, sort string) ([]model.Post, int, error)
	TagPost(slug string, limit int, offset uint, search string, sort_by string, sort string) (model.TagWithPost, int, error)
	UserPost(username string, limit int, offset uint, search string, sort_by string, sort string) (model.UserWithPost, int, error)
	Show(slug string) (model.PostSingle, error)
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
	_order := database.OrderBy(sort_by, sort)
	_limit := database.Limit(limit, offset)

	count_query := fmt.Sprintf(`SELECT count(*) FROM posts LEFT JOIN users ON users.uuid = posts.user_uuid LEFT JOIN tags ON tags.uuid = posts.tag_uuid %s AND posts.deleted_at IS NULL`, _conditions)
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

	if items == nil {
		return []model.Post{}, count, nil
	}

	return items, count, nil
}

func (repo *PostRepo) TagPost(slug string, limit int, offset uint, search string, sort_by string, sort string) (model.TagWithPost, int, error) {
	_limit := database.Limit(limit, offset)
	_conditions := database.SearchOther([]string{"title", "content", "users.name"}, search, "posts.deleted_at")
	_order := database.OrderBy(sort_by, sort)
	_select := fmt.Sprintf(`
	tags.uuid,
    tags.name,
    tags.slug,
    IFNULL(
        CONCAT('[', 
            IF(
                COUNT(posts.uuid) > 0,
                GROUP_CONCAT(
                    DISTINCT JSON_OBJECT(
                        'uuid', posts.uuid,
                        'tag_uuid', posts.tag_uuid,
                        'user_uuid', posts.user_uuid,
                        'title', posts.title,
                        'thumbnail', posts.thumbnail,
                        'content', posts.content,
                        'keyword', posts.keyword,
                        'slug', posts.slug,
                        'is_active', posts.is_active,
                        'is_highlight', posts.is_highlight,
                        'user', IF(
                            user_uuid_u IS NULL,
                            NULL,
                            JSON_OBJECT(
                                'uuid', user_uuid_u,
                                'name', user_name,
                                'username', user_username,
                                'email', user_email,
                                'created_at', user_created_at,
                                'updated_at', user_updated_at
                            )
                        ),
                        'tag', IF(
                            tag_uuid_u IS NULL,
                            NULL,
                            JSON_OBJECT(
                                'uuid', tag_uuid_u,
                                'name', tag_name,
                                'slug', tag_slug,
                                'created_at', tag_created_at,
                                'updated_at', tag_updated_at
                            )
                        )
                    )
                ),
                ''
            )
        , ']'),
        '[]'
    ) AS post,
    tags.created_at,
    tags.updated_at
	`)

	count_query := fmt.Sprintf(`SELECT count(*) FROM posts LEFT JOIN tags ON tags.uuid = posts.tag_uuid LEFT JOIN users ON users.uuid = posts.user_uuid WHERE tags.slug = ? AND posts.is_active = true AND posts.deleted_at IS NULL %s`, _conditions)
	var count int
	_ = repo.db.QueryRow(count_query, slug).Scan(&count)

	query := fmt.Sprintf(`SELECT %s FROM tags 
	LEFT JOIN (
		SELECT 
			posts.*, 
            users.uuid as user_uuid_u, users.name as user_name, users.username as user_username, users.email as user_email, users.created_at as user_created_at, users.updated_at as user_updated_at, 
        	tags.uuid as tag_uuid_u, tags.name as tag_name, tags.slug as tag_slug, tags.created_at as tag_created_at, tags.updated_at as tag_updated_at 
    	FROM posts
	    LEFT JOIN users ON users.uuid = posts.user_uuid
	    LEFT JOIN tags ON tags.uuid = posts.tag_uuid
	    WHERE tags.slug = ? AND posts.is_active = true %s
	    %s %s
	) as posts ON posts.tag_uuid = tags.uuid
	WHERE tags.slug = ? AND tags.is_active = true
	GROUP BY tags.uuid, tags.name, tags.slug, tags.created_at, tags.updated_at
	`, _select, _conditions, _order, _limit)

	var items model.TagWithPost

	err := repo.db.QueryRowContext(context.Background(), query, slug, slug).Scan(
		&items.UUID,
		&items.Name,
		&items.Slug,
		&items.Post,
		&items.CreatedAt,
		&items.UpdatedAt,
	)

	if err != nil {
		return model.TagWithPost{}, 0, err
	}

	return items, count, nil
}

func (repo *PostRepo) UserPost(username string, limit int, offset uint, search string, sort_by string, sort string) (model.UserWithPost, int, error) {
	_limit := database.Limit(limit, offset)
	_conditions := database.SearchOther([]string{"title", "content", "users.name"}, search, "posts.deleted_at")
	_order := database.OrderBy(sort_by, sort)
	_select := fmt.Sprintf(`
	users.uuid,
    users.name,
    users.username,
    users.email,
    IFNULL(
        CONCAT('[', 
            IF(
                COUNT(posts.uuid) > 0,
                GROUP_CONCAT(
                    DISTINCT JSON_OBJECT(
                        'uuid', posts.uuid,
                        'tag_uuid', posts.tag_uuid,
                        'user_uuid', posts.user_uuid,
                        'title', posts.title,
                        'thumbnail', posts.thumbnail,
                        'content', posts.content,
                        'keyword', posts.keyword,
                        'slug', posts.slug,
                        'is_active', posts.is_active,
                        'is_highlight', posts.is_highlight,
                        'user', IF(
                            user_uuid_u IS NULL,
                            NULL,
                            JSON_OBJECT(
                                'uuid', user_uuid_u,
                                'name', user_name,
                                'username', user_username,
                                'email', user_email,
                                'created_at', user_created_at,
                                'updated_at', user_updated_at
                            )
                        ),
                        'tag', IF(
                            tag_uuid_u IS NULL,
                            NULL,
                            JSON_OBJECT(
                                'uuid', tag_uuid_u,
                                'name', tag_name,
                                'slug', tag_slug,
                                'created_at', tag_created_at,
                                'updated_at', tag_updated_at
                            )
                        )
                    )
                ),
                ''
            )
        , ']'),
        '[]'
    ) AS post,
    users.created_at,
    users.updated_at
	`)

	count_query := fmt.Sprintf(`SELECT count(*) FROM posts JOIN users ON users.uuid = posts.user_uuid WHERE username = ? AND posts.is_active = true AND posts.deleted_at IS NULL %s`, _conditions)
	var count int
	_ = repo.db.QueryRow(count_query, username).Scan(&count)

	query := fmt.Sprintf(`SELECT %s FROM users
	LEFT JOIN (
		SELECT 
			posts.*, 
            users.uuid as user_uuid_u, users.name as user_name, users.username as user_username, users.email as user_email, users.created_at as user_created_at, users.updated_at as user_updated_at, 
        	tags.uuid as tag_uuid_u, tags.name as tag_name, tags.slug as tag_slug, tags.created_at as tag_created_at, tags.updated_at as tag_updated_at 
    	FROM posts
	    LEFT JOIN users ON users.uuid = posts.user_uuid
	    LEFT JOIN tags ON tags.uuid = posts.tag_uuid
	    WHERE users.username = ? AND posts.is_active = true %s
	    %s %s
	) as posts ON posts.user_uuid = users.uuid
	WHERE username = ?
	GROUP BY users.uuid, users.name, users.username, users.email, users.created_at, users.updated_at
	`, _select, _conditions, _order, _limit)

	var items model.UserWithPost

	err := repo.db.QueryRowContext(context.Background(), query, username, username).Scan(
		&items.UUID,
		&items.Name,
		&items.Username,
		&items.Email,
		&items.Post,
		&items.CreatedAt,
		&items.UpdatedAt,
	)

	if err != nil {
		return model.UserWithPost{}, 0, err
	}

	return items, count, nil
}

func (repo *PostRepo) Show(slug string) (model.PostSingle, error) {
	var post model.PostSingle
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
	WHERE posts.slug = $1 AND posts.is_active = true AND posts.deleted_at IS NULL LIMIT 1
	`

	err := repo.db.QueryRowContext(context.Background(), query, slug).Scan(
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

	if err != nil {
		return model.PostSingle{}, err
	}

	similiar_query := `
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
	CASE 
    	WHEN users.uuid IS NULL THEN null
    	ELSE jsonb_build_object(
        	'uuid', users.uuid,
        	'name', users.name,
        	'username', users.username,
        	'email', users.email,
        	'created_at', users.created_at,
        	'updated_at', users.updated_at
    	)
	END AS user,
	CASE 
    	WHEN tags.uuid IS NULL THEN null
		ELSE jsonb_build_object(
        	'uuid', tags.uuid,
        	'name', tags.name,
			'slug', tags.slug,
        	'created_at', tags.created_at,
        	'updated_at', tags.updated_at
		)
	END AS tag
	FROM posts LEFT JOIN users ON users.uuid = posts.user_uuid LEFT JOIN tags ON tags.uuid = posts.tag_uuid
	WHERE posts.title LIKE $1 AND posts.slug != $2 AND posts.is_active = true AND posts.deleted_at IS NULL LIMIT 5
	`

	rows, err := repo.db.QueryContext(context.Background(), similiar_query, "%"+post.Title+"%", post.Slug)
	if err != nil {
		return post, err
	}

	defer rows.Close()
	var items []*model.Post
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
			return post, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return post, err
	}
	if err := rows.Err(); err != nil {
		return post, err
	}

	post.Similiar = items

	return post, err
}

func NewPostRepo(db *database.DB) PostRepository {
	return &PostRepo{db}
}
