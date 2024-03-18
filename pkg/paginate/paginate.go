package paginate

import (
	"strconv"

	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

var logr = logger.GetLogger()

func Paginate(c *fiber.Ctx) (pageNo int, pageSize int, search string, sortBy string, sort string) {
	ps := c.Query("limit")
	pn := c.Query("page")
	search = c.Query("search")
	sortBy = c.Query("sort_by")
	sort = c.Query("sort")
	pageSize, pageNo = 10, 1

	if sortBy == "" {
		sortBy = "id"
	}

	if sort == "" {
		sort = "DESC"
	}

	if len(ps) > 0 {
		psInt, err := strconv.Atoi(ps)
		if err != nil {
			logr.Error(err)
		} else {
			pageSize = psInt
		}
	}

	if len(pn) > 0 {
		pnInt, err := strconv.Atoi(pn)
		if err != nil {
			logr.Error(err)
		} else {
			pageNo = pnInt
		}
	}

	return pageNo, pageSize, search, sortBy, sort
}
