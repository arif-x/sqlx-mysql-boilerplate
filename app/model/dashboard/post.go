package dashboard

import (
	"time"

	jsonutil "github.com/arif-x/sqlx-mysql-boilerplate/pkg/json"
	"github.com/google/uuid"
)

type Post struct {
	UUID        uuid.UUID         `db:"uuid" json:"uuid"`
	TagUUID     uuid.UUID         `db:"tag_uuid" json:"tag_uuid"`
	UserUUID    uuid.UUID         `db:"user_uuid" json:"user_uuid"`
	Title       string            `db:"title" json:"title"`
	Thumbnail   string            `db:"thumbnail" json:"thumbnail"`
	Content     string            `db:"content" json:"content"`
	Slug        string            `db:"slug" json:"slug"`
	Keyword     string            `db:"keyword" json:"keyword"`
	IsActive    bool              `db:"is_active" json:"is_active"`
	IsHighlight bool              `db:"is_highlight" json:"is_highlight"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time        `db:"deleted_at" json:"deleted_at"`
	User        *jsonutil.JSONRaw `db:"user" json:"user"`
	Tag         *jsonutil.JSONRaw `db:"tag" json:"tag"`
}

type PostShow struct {
	UUID        uuid.UUID         `db:"uuid" json:"uuid"`
	TagUUID     uuid.UUID         `db:"tag_uuid" json:"tag_uuid"`
	UserUUID    uuid.UUID         `db:"user_uuid" json:"user_uuid"`
	Thumbnail   string            `db:"thumbnail" json:"thumbnail"`
	Title       string            `db:"title" json:"title"`
	Content     string            `db:"content" json:"content"`
	Slug        string            `db:"slug" json:"slug"`
	Keyword     string            `db:"keyword" json:"keyword"`
	IsActive    bool              `db:"is_active" json:"is_active"`
	IsHighlight bool              `db:"is_highlight" json:"is_highlight"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time        `db:"updated_at" json:"updated_at"`
	User        *jsonutil.JSONRaw `db:"user" json:"user"`
	Tag         *jsonutil.JSONRaw `db:"tag" json:"tag"`
}

type StorePost struct {
	TagUUID     uuid.UUID `json:"tag_uuid" form:"tag_uuid"`
	UserUUID    uuid.UUID `json:"user_uuid" form:"user_uuid"`
	Title       string    `json:"title" form:"title"`
	Thumbnail   string    `json:"thumbnail" form:"thumbnail"`
	Content     string    `json:"content" form:"content"`
	Slug        string    `json:"slug" form:"slug"`
	Keyword     string    `json:"keyword" form:"keyword"`
	IsActive    bool      `json:"is_active" form:"is_active"`
	IsHighlight bool      `json:"is_highlight" form:"is_highlight"`
}

type UpdatePost struct {
	TagUUID     uuid.UUID `json:"tag_uuid" form:"tag_uuid"`
	UserUUID    uuid.UUID `json:"user_uuid" form:"user_uuid"`
	Title       string    `json:"title" form:"title"`
	Thumbnail   string    `json:"thumbnail" form:"thumbnail"`
	Content     string    `json:"content" form:"content"`
	Slug        string    `json:"slug" form:"slug"`
	Keyword     string    `json:"keyword" form:"keyword"`
	IsActive    bool      `json:"is_active" form:"is_active"`
	IsHighlight bool      `json:"is_highlight" form:"is_highlight"`
}
