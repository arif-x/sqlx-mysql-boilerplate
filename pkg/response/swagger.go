package response

import (
	dashboard "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	public "github.com/arif-x/sqlx-mysql-boilerplate/app/model/public"
)

type AuthResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type AuthWithPermissionResponse struct {
	Status     bool     `json:"status"`
	Message    string   `json:"message"`
	Data       string   `json:"data"`
	Permission []string `json:"permission"`
}

type ErrorResponse struct {
	Status  bool   `json:"status" example:"false"`
	Message string `json:"message"`
}

type UserResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    dashboard.User `json:"data"`
}

type UsersResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    []dashboard.User `json:"data"`
	Limit   int              `json:"limit"`
	Page    int              `json:"page"`
	Total   int              `json:"total"`
}

type TagResponse struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    dashboard.Tag `json:"data"`
}

type TagsResponse struct {
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    []dashboard.Tag `json:"data"`
	Limit   int             `json:"limit"`
	Page    int             `json:"page"`
	Total   int             `json:"total"`
}

type PostResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    dashboard.Post `json:"data"`
}

type PostsResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    []dashboard.Post `json:"data"`
	Limit   int              `json:"limit"`
	Page    int              `json:"page"`
	Total   int              `json:"total"`
}

type RoleResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    dashboard.Role `json:"data"`
}

type RolesResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    []dashboard.Role `json:"data"`
	Limit   int              `json:"limit"`
	Page    int              `json:"page"`
	Total   int              `json:"total"`
}

type PermissionResponse struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message"`
	Data    dashboard.Permission `json:"data"`
}

type PermissionsResponse struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    []dashboard.Permission `json:"data"`
	Limit   int                    `json:"limit"`
	Page    int                    `json:"page"`
	Total   int                    `json:"total"`
}

type PublicPostResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    public.Post `json:"data"`
}

type PublicPostsResponse struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    []public.Post `json:"data"`
	Limit   int           `json:"limit"`
	Page    int           `json:"page"`
	Total   int           `json:"total"`
}

type PublicPostsByTagResponse struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    public.TagWithPost `json:"data"`
	Limit   int                `json:"limit"`
	Page    int                `json:"page"`
	Total   int                `json:"total"`
}

type PublicPostsByUserResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    public.UserWithPost `json:"data"`
	Limit   int                 `json:"limit"`
	Page    int                 `json:"page"`
	Total   int                 `json:"total"`
}

type SyncPermissionResponse struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message"`
	Data    dashboard.Permission `json:"data"`
}
