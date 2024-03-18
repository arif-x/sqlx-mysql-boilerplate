package public

import (
	"database/sql"
	"log"

	repo "github.com/arif-x/sqlx-mysql-boilerplate/app/repository/public"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/paginate"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// PublicPostIndex func gets all post.
// @Description Get all post.
// @Summary Get all post
// @Tags Public Post
// @Accept json
// @Produce json
// @Success 200 {object} response.PublicPostsResponse
// @Failure 400,403,404 {object} response.ErrorResponse "Error"
// @Router /api/v1/public/post [get]
func PostIndex(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)

	if sort_by == "id" {
		sort_by = "posts.id"
	}

	repository := repo.NewPostRepo(database.GetDB())

	posts, count, err := repository.Index(limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Index(c, page, limit, count, posts)
}

// PublicPostByTag func gets post by tag.
// @Description Get post by tag.
// @Summary Get post by tag
// @Tags Public Post
// @Accept json
// @Produce json
// @Param slug path string true "Tag Slug" default(tag-1)
// @Success 200 {object} response.PublicPostsByTagResponse
// @Failure 400,403,404 {object} response.ErrorResponse "Error"
// @Router /api/v1/public/post/tag/{slug} [get]
func TagPost(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)
	slug := c.Params("slug")
	repository := repo.NewPostRepo(database.GetDB())

	posts, count, err := repository.TagPost(slug, limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			log.Println(err)
			return response.InternalServerError(c, err)
		}
	}

	return response.Index(c, page, limit, count, posts)
}

// PublicPostByUser func gets post by user.
// @Description Get post by user.
// @Summary Get post by user
// @Tags Public Post
// @Accept json
// @Produce json
// @Param username path string true "Username" default(username1)
// @Success 200 {object} response.PublicPostsByUserResponse
// @Failure 400,403,404 {object} response.ErrorResponse "Error"
// @Router /api/v1/public/post/user/{username} [get]
func UserPost(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)
	username := c.Params("username")
	repository := repo.NewPostRepo(database.GetDB())

	posts, count, err := repository.UserPost(username, limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Index(c, page, limit, count, posts)
}

// PublicPostShow func gets single post.
// @Description Get single post.
// @Summary Get single post
// @Tags Public Post
// @Accept json
// @Produce json
// @Param slug path string true "Post Slug" default(title-1)
// @Success 200 {object} response.PostResponse
// @Failure 400,403,404 {object} response.ErrorResponse "Error"
// @Router /api/v1/public/post/{slug} [get]
func PostShow(c *fiber.Ctx) error {
	slug := c.Params("slug")

	repository := repo.NewPostRepo(database.GetDB())
	post, err := repository.Show(slug)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Show(c, post)
}
