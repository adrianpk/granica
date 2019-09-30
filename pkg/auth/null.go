package auth

import (
	"database/sql"
	"strconv"
)

func toNullInt64(s string) sql.NullInt64 {
	i, err := strconv.Atoi(s)
	return sql.NullInt64{
		Int64: int64(i),
		Valid: err == nil,
	}
}

func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
