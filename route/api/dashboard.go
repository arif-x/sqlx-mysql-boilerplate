package api

import (
	controllers "github.com/arif-x/sqlx-mysql-boilerplate/app/http/controller/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/app/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func Dashboard(a *fiber.App) {
	dashboard := a.Group("/api/v1/dashboard", middleware.JWTProtected(), middleware.Email(), middleware.IsActive())

	user := dashboard.Group("/user")
	user.Get("/", middleware.Permission("user-index"), controllers.UserIndex)
	user.Get("/:id", middleware.Permission("user-show"), controllers.UserShow)
	user.Post("/", middleware.Permission("user-store"), controllers.UserStore)
	user.Put("/:id", middleware.Permission("user-update"), controllers.UserUpdate)
	user.Delete("/:id", middleware.Permission("user-destroy"), controllers.UserDestroy)

	tag := dashboard.Group("/tags")
	tag.Get("/", middleware.Permission("tags-index"), controllers.TagIndex)
	tag.Get("/:id", middleware.Permission("tags-show"), controllers.TagShow)
	tag.Post("/", middleware.Permission("tags-store"), controllers.TagStore)
	tag.Put("/:id", middleware.Permission("tags-update"), controllers.TagUpdate)
	tag.Delete("/:id", middleware.Permission("tags-destroy"), controllers.TagDestroy)

	post := dashboard.Group("/post")
	post.Get("/", middleware.Permission("post-index"), controllers.PostIndex)
	post.Get("/:id", middleware.Permission("post-show"), controllers.PostShow)
	post.Post("/", middleware.Permission("post-store"), controllers.PostStore)
	post.Put("/:id", middleware.Permission("post-update"), controllers.PostUpdate)
	post.Delete("/:id", middleware.Permission("post-destroy"), controllers.PostDestroy)

	role := dashboard.Group("/role")
	role.Get("/", middleware.Permission("role-index"), controllers.RoleIndex)
	role.Get("/:id", middleware.Permission("role-show"), controllers.RoleShow)
	role.Post("/", middleware.Permission("role-store"), controllers.RoleStore)
	role.Put("/:id", middleware.Permission("role-update"), controllers.RoleUpdate)
	role.Delete("/:id", middleware.Permission("role-destroy"), controllers.RoleDestroy)

	permission := dashboard.Group("/permission")
	permission.Get("/", middleware.Permission("permission-index"), controllers.PermissionIndex)
	permission.Get("/:id", middleware.Permission("permission-show"), controllers.PermissionShow)
	permission.Post("/", middleware.Permission("permission-store"), controllers.PermissionStore)
	permission.Put("/:id", middleware.Permission("permission-update"), controllers.PermissionUpdate)
	permission.Delete("/:id", middleware.Permission("permission-destroy"), controllers.PermissionDestroy)

	sync_permission := dashboard.Group("/sync-permission")
	sync_permission.Get("/:id", middleware.Permission("permission-index"), controllers.SyncPermissionShow)
	sync_permission.Put("/:id", middleware.Permission("permission-index"), controllers.SyncPermissionUpdate)
}
