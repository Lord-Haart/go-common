package dbhelper

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

func ParseNullInt64(s string) sql.NullInt64 {
	if i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err != nil {
		return sql.NullInt64{}
	} else {
		return sql.NullInt64{Valid: true, Int64: i}
	}
}

func ParseNullInt32(s string) sql.NullInt32 {
	if i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err != nil {
		return sql.NullInt32{}
	} else {
		return sql.NullInt32{Valid: true, Int32: int32(i)}
	}
}

func ParseNullInt16(s string) sql.NullInt16 {
	if i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err != nil {
		return sql.NullInt16{}
	} else {
		return sql.NullInt16{Valid: true, Int16: int16(i)}
	}
}

func CoalseceString(v sql.NullString, dv string) string {
	if v.Valid {
		return v.String
	} else {
		return dv
	}
}

func CoalseceBool(v sql.NullBool, dv bool) bool {
	if v.Valid {
		return v.Bool
	} else {
		return dv
	}
}

func CoalseceInt64(v sql.NullInt64, dv int64) int64 {
	if v.Valid {
		return v.Int64
	} else {
		return dv
	}
}

func CoalseceInt32(v sql.NullInt32, dv int32) int32 {
	if v.Valid {
		return v.Int32
	} else {
		return dv
	}
}

func CoalseceInt16(v sql.NullInt16, dv int16) int16 {
	if v.Valid {
		return v.Int16
	} else {
		return dv
	}
}

func CoalseceByte(v sql.NullByte, dv byte) byte {
	if v.Valid {
		return v.Byte
	} else {
		return dv
	}
}

func CoalseceTime(v sql.NullTime, dv time.Time) time.Time {
	if v.Valid {
		return v.Time
	} else {
		return dv
	}
}

func CoalseceFloat64(v sql.NullFloat64, dv float64) float64 {
	if v.Valid {
		return v.Float64
	} else {
		return dv
	}
}
