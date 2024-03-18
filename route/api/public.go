package api

import (
	controllers "github.com/arif-x/sqlx-mysql-boilerplate/app/http/controller/public"
	"github.com/gofiber/fiber/v2"
)

func Public(a *fiber.App) {
	public := a.Group("/api/v1/public")

	public.Get("/post", controllers.PostIndex)
	public.Get("/post/tag/:slug", controllers.TagPost)
	public.Get("/post/user/:username", controllers.UserPost)
	public.Get("/post/:slug", controllers.PostShow)
}
