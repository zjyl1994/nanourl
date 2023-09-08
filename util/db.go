package util

import (
	"github.com/mattn/go-sqlite3"
)

func IsSqliteDuplicateError(err error) bool {
	if err, ok := err.(sqlite3.Error); ok {
		if err.ExtendedCode == sqlite3.ErrConstraintUnique {
			return true
		}
	}
	return false
}

func PageHelper(page, pageSize int) (limit, offset int) {
	page, pageSize = PageNormalize(page, pageSize)

	limit = pageSize
	offset = (page - 1) * pageSize
	return limit, offset
}

func TotalPage(totalRows, pageSize int) int {
	return (totalRows + pageSize - 1) / pageSize
}

func PageNormalize(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return page, pageSize
}
