package dashboard

import (
	"database/sql"
	"fmt"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	repo "github.com/arif-x/sqlx-mysql-boilerplate/app/repository/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/paginate"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// TagIndex func gets all post tag.
// @Description Get all post tag.
// @Summary Get all post tag
// @Tags Post Tag
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Param search query string false "Search"
// @Param sort_by query string false "Sort By" Enums(id, name)
// @Param sort query string false "Sort" Enums(ASC, DESC)
// @Success 200 {object} response.TagsResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/tags [get]
func TagIndex(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)
	repository := repo.NewTagRepo(database.GetDB())

	tags, count, err := repository.Index(limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Index(c, page, limit, count, tags)
}

// TagShow func gets single post tag.
// @Description Get single post tag.
// @Summary Get single post tag
// @Tags Post Tag
// @Accept json
// @Produce json
// @Param id path string true "Post Tag ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.TagResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/tags/{id} [get]
func TagShow(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewTagRepo(database.GetDB())
	tag, err := repository.Show(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Show(c, tag)
}

// TagStore func create post tag.
// @Description Create post tag.
// @Summary Create post tag
// @Tags Post Tag
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name" default(Tag Name)
// @Param is_active formData bool true "Is Active"
// @Success 200 {object} response.TagResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/tags [post]
func TagStore(c *fiber.Ctx) error {
	tag := &model.StoreTag{}

	if err := c.BodyParser(tag); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewTagRepo(database.GetDB())

	tag.Slug = repository.GetSlug(tag.Name, nil)

	res, err := repository.Store(tag)

	fmt.Print(err)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Store(c, res)
}

// TagUpdate func update post tag.
// @Description Update post tag.
// @Summary Update post tag
// @Tags Post Tag
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Post Tag ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Param name formData string true "Name" default(Tag Name Update)
// @Param is_active formData bool true "Is Active"
// @Success 200 {object} response.TagResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/tags/{id} [put]
func TagUpdate(c *fiber.Ctx) error {
	ID := c.Params("id")

	tag := &model.UpdateTag{}

	if err := c.BodyParser(tag); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewTagRepo(database.GetDB())

	tag.Slug = repository.GetSlug(tag.Name, &ID)

	res, err := repository.Update(ID, tag)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Update(c, res)
}

// TagDestroy func delete post tag.
// @Description Delete post tag.
// @Summary Delete post tag
// @Tags Post Tag
// @Accept json
// @Produce json
// @Param id path string true "Post Tag ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.TagResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/tags/{id} [delete]
func TagDestroy(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewTagRepo(database.GetDB())
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
