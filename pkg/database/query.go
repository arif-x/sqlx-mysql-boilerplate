package database

import "fmt"

func Search(columns []string, search string, deleted_at string) string {
	conditions := ""
	if search != "" {
		conditions = "WHERE"
		if len(columns) > 0 {
			conditions += " ( "
		}
		for i := 0; i < len(columns); i++ {
			if i == (len(columns) - 1) {
				conditions += fmt.Sprintf(" %s LIKE '%s' ", columns[i], "%"+search+"%")
				conditions += " ) "
			} else {
				conditions += fmt.Sprintf(" %s LIKE '%s' OR ", columns[i], "%"+search+"% ")
			}
		}
		conditions += fmt.Sprintf(" AND %s IS NULL", deleted_at)
	} else {
		conditions += fmt.Sprintf(" WHERE %s IS NULL", deleted_at)
	}

	return conditions
}

func SearchOther(columns []string, search string, deleted_at string) string {
	conditions := ""
	if search != "" {
		conditions = "AND ("
		for i := 0; i < len(columns); i++ {
			if i == (len(columns) - 1) {
				conditions += fmt.Sprintf(" lower(%s) LIKE lower('%s') ", columns[i], "%"+search+"%")
			} else {
				conditions += fmt.Sprintf(" lower(%s) LIKE lower('%s') OR ", columns[i], "%"+search+"%")
			}
		}
		conditions += ")"
		conditions += fmt.Sprintf(" AND %s IS NULL", deleted_at)
	} else {
		conditions += fmt.Sprintf(" AND %s IS NULL", deleted_at)
	}

	return conditions
}

func OrderBy(sort_by string, sort string) string {
	return fmt.Sprintf(" ORDER BY %s %s", sort_by, sort)
}

func Limit(limit int, offset uint) string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}
