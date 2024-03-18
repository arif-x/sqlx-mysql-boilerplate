package public

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
	Keyword     string            `db:"keyword" json:"keyword"`
	Slug        string            `db:"slug" json:"slug"`
	IsActive    string            `db:"is_active" json:"is_active"`
	IsHighlight string            `db:"is_highlight" json:"is_highlight"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time        `db:"deleted_at" json:"deleted_at"`
	User        *jsonutil.JSONRaw `db:"user" json:"user"`
	Tag         *jsonutil.JSONRaw `db:"tag" json:"tag"`
}

type PostSingle struct {
	UUID        uuid.UUID         `db:"uuid" json:"uuid"`
	TagUUID     uuid.UUID         `db:"tag_uuid" json:"tag_uuid"`
	UserUUID    uuid.UUID         `db:"user_uuid" json:"user_uuid"`
	Title       string            `db:"title" json:"title"`
	Thumbnail   string            `db:"thumbnail" json:"thumbnail"`
	Content     string            `db:"content" json:"content"`
	Keyword     string            `db:"keyword" json:"keyword"`
	Slug        string            `db:"slug" json:"slug"`
	IsActive    string            `db:"is_active" json:"is_active"`
	IsHighlight string            `db:"is_highlight" json:"is_highlight"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time        `db:"deleted_at" json:"deleted_at"`
	User        *jsonutil.JSONRaw `db:"user" json:"user"`
	Tag         *jsonutil.JSONRaw `db:"tag" json:"tag"`
	Similiar    []*Post           `db:"similiar" json:"similiar"`
}

type UserWithPost struct {
	UUID      uuid.UUID         `db:"uuid" json:"uuid"`
	Name      string            `db:"name" json:"name"`
	Username  string            `db:"username" json:"username"`
	Email     string            `db:"email" json:"email"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time        `db:"updated_at" json:"updated_at"`
	Post      *jsonutil.JSONRaw `db:"post" json:"post"`
}

type TagWithPost struct {
	UUID      uuid.UUID         `db:"uuid" json:"uuid"`
	Name      string            `db:"name" json:"name"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time        `db:"updated_at" json:"updated_at"`
	Post      *jsonutil.JSONRaw `db:"post" json:"post"`
	Slug      string            `db:"slug" json:"slug"`
}
