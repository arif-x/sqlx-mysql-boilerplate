package dashboard

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	repo "github.com/arif-x/sqlx-mysql-boilerplate/app/repository/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	fileHelper "github.com/arif-x/sqlx-mysql-boilerplate/pkg/file"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/paginate"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// PostIndex func gets all post.
// @Description Get all post.
// @Summary Get all post
// @Tags Post
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Param search query string false "Search"
// @Param sort_by query string false "Sort By" Enums(posts.id, title, content)
// @Param sort query string false "Sort" Enums(ASC, DESC)
// @Success 200 {object} response.PostResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/post [get]
func PostIndex(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)

	repository := repo.NewPostRepo(database.GetDB())

	posts, count, err := repository.Index(limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Index(c, page, limit, count, posts)
}

// PostShow func gets single post.
// @Description Get single post.
// @Summary Get single post
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post ID" default(f72cb686-2fc3-4147-8183-f93684780765)
// @Success 200 {object} response.PostResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/post/{id} [get]
func PostShow(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewPostRepo(database.GetDB())
	post, err := repository.Show(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Show(c, post)
}

// PostStore func create post.
// @Description Create post.
// @Summary Create post
// @Tags Post
// @Accept multipart/form-data
// @Produce json
// @Param user_uuid formData string true "User UUID" default(87c76e22-e2f0-4ebf-bda8-56802c0a0577)
// @Param tag_uuid formData string true "Post Tag UUID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Param title formData string true "Title" default(Title)
// @Param thumbnail formData file true "Thumbnail"
// @Param content formData string true "Content" default(Content)
// @Param keyword formData string true "Keyword" default(keyword 1, keyword 2)
// @Param is_active formData bool true "Is Active"
// @Param is_highlight formData bool true "Is Highlight"
// @Success 200 {object} response.PostResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/post [post]
func PostStore(c *fiber.Ctx) error {
	post := &model.StorePost{}

	if err := c.BodyParser(post); err != nil {
		return response.BadRequest(c, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return response.BadRequest(c, err)
	}

	thumbnail := form.File["thumbnail"]

	thumbnail_data := ""

	allowedExtensions := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".gif":  {},
		".svg":  {},
	}

	repository := repo.NewPostRepo(database.GetDB())

	for _, file := range thumbnail {
		ext := filepath.Ext(file.Filename)
		if _, allowed := allowedExtensions[ext]; !allowed {
			return response.InternalServerError(c, errors.New("Unsupported file extension"))
		}
		filename := fileHelper.GenerateUniqueFilename(file.Filename)
		if err := fileHelper.SaveFile(file, "./upload/post/", filename); err != nil {
			return response.InternalServerError(c, errors.New("Can't Upload File"))
		}
		thumbnail_data = os.Getenv("APP_FULL_URL") + "/upload/post/" + filename
	}

	post.Thumbnail = thumbnail_data
	post.Slug = repository.GetSlug(post.Title, nil)

	res, err := repository.Store(post)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Store(c, res)
}

// PostUpdate func update post.
// @Description Update post.
// @Summary Update post
// @Tags Post
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Post ID" default(f72cb686-2fc3-4147-8183-f93684780765)
// @Param user_uuid formData string true "User UUID" default(87c76e22-e2f0-4ebf-bda8-56802c0a0577)
// @Param tag_uuid formData string true "Post Tag UUID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Param title formData string true "Title" default(Title Update)
// @Param thumbnail formData file false "Thumbnail"
// @Param content formData string true "Content" default(Content Update)
// @Param keyword formData string true "Keyword" default(keyword 1, keyword 2)
// @Param is_active formData bool true "Is Active"
// @Param is_highlight formData bool true "Is Highlight"
// @Success 200 {object} response.PostResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/post/{id} [put]
func PostUpdate(c *fiber.Ctx) error {
	ID := c.Params("id")

	post := &model.UpdatePost{}

	if err := c.BodyParser(post); err != nil {
		return response.BadRequest(c, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return response.BadRequest(c, err)
	}

	thumbnail := form.File["thumbnail"]

	thumbnail_data := ""

	allowedExtensions := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".gif":  {},
		".svg":  {},
	}

	repository := repo.NewPostRepo(database.GetDB())

	for _, file := range thumbnail {
		ext := filepath.Ext(file.Filename)
		if _, allowed := allowedExtensions[ext]; !allowed {
			return response.InternalServerError(c, errors.New("Unsupported file extension"))
		}
		filename := fileHelper.GenerateUniqueFilename(file.Filename)
		if err := fileHelper.SaveFile(file, "./upload/post/", filename); err != nil {
			return response.InternalServerError(c, errors.New("Can't Upload File"))
		}
		if filename != "" {
			thumbnail_data = os.Getenv("APP_FULL_URL") + "/upload/post/" + filename
		}
	}

	post.Thumbnail = thumbnail_data
	post.Slug = repository.GetSlug(post.Title, &ID)

	res, err := repository.Update(ID, post)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			log.Println(err)
			return response.InternalServerError(c, err)
		}
	}

	return response.Update(c, res)
}

// PostDestroy func delete post.
// @Description Delete post.
// @Summary Delete post
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post ID" default(f72cb686-2fc3-4147-8183-f93684780765)
// @Success 200 {object} response.PostResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/post/{id} [delete]
func PostDestroy(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewPostRepo(database.GetDB())
	res, err := repository.Destroy(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Destroy(c, res)
}
