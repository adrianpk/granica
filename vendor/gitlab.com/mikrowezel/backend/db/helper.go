package db

import (
	"database/sql"
	"strconv"
)

func ToNullInt64(s string) sql.NullInt64 {
	i, err := strconv.Atoi(s)
	return sql.NullInt64{
		Int64: int64(i),
		Valid: err == nil,
	}
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}
