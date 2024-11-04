package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Boolean 表示可以为空的bool
type Boolean struct {
	Valid bool
	V     bool
}

// String 表示可以为空的string。
type String struct {
	Valid bool
	V     string
}

// Short 表示可以为空的int16
type Short struct {
	Valid bool
	V     int16
}

// Integer 表示可以为空的int32
type Integer struct {
	Valid bool
	V     int32
}

// Long 表示可以为空的int64
type Long struct {
	Valid bool
	V     int64
}

// Timestamp 表示可以为空的Time
type Timestamp struct {
	Valid bool
	V     time.Time
}

func parseBoolean(s string) (v, ok bool) {
	s0 := strings.ToLower(strings.TrimSpace(s))
	if s0 == "null" || s0 == "" || s0 == `''` || s0 == `""` {
		// null
	} else if s0 == "true" || s0 == "1" || s0 == "yes" || s0 == "on" || s0 == "t" || s0 == "y" {
		v = true
		ok = true
	} else if s0 == "false" || s0 == "0" || s0 == "no" || s0 == "off" || s0 == "f" || s0 == "n" {
		v = false
		ok = true
	} else if s1, err := strconv.Unquote(s0); err != nil {
		// illegal
	} else {
		s1 = strings.TrimSpace(s1)
		if s1 == "true" || s1 == "1" || s1 == "yes" || s1 == "on" || s1 == "t" || s1 == "y" {
			v = true
			ok = true
		} else if s1 == "false" || s1 == "0" || s1 == "no" || s1 == "off" || s1 == "f" || s1 == "n" {
			v = false
			ok = true
		}
	}

	return
}

func parseInt16(s string) (v int16, ok bool) {
	if v0, err := strconv.ParseInt(s, 10, 16); err != nil {
		// illegal
	} else {
		v = int16(v0)
		ok = true
	}

	return
}

func parseInt32(s string) (v int32, ok bool) {
	if v0, err := strconv.ParseInt(s, 10, 32); err != nil {
		// illegal
	} else {
		v = int32(v0)
		ok = true
	}

	return
}

func parseInt64(s string) (v int64, ok bool) {
	if v0, err := strconv.ParseInt(s, 10, 64); err != nil {
		// illegal
	} else {
		v = int64(v0)
		ok = true
	}

	return
}

func parseTime(s string) (v time.Time, ok bool) {
	if v0, err := time.Parse(time.RFC3339, s); err != nil {
		// illegal
	} else {
		v = v0
		ok = true
	}

	return
}

func (b Boolean) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(b.V)
	}
}

func (b *Boolean) UnmarshalJSON(data []byte) error {
	if v, ok := parseBoolean(string(data)); ok {
		b.Valid = true
		b.V = v
	} else {
		b.Valid = false
	}

	return nil
}

func (b *Boolean) Merge(o Boolean) *Boolean {
	if o.Valid {
		b.Valid = true
		b.V = o.V
	}

	return b
}

func (b *Boolean) Eq(ov bool) bool {
	return b.Valid && b.V == ov
}

func (b Boolean) Coalsece(o bool) bool {
	if b.Valid {
		return b.V
	} else {
		return o
	}
}

// Scan implements the [Scanner] interface.
func (b *Boolean) Scan(value any) error {
	if value == nil {
		b.V, b.Valid = false, false
		return nil
	}
	b.Valid = true
	switch v := value.(type) {
	case bool:
		b.V = v
	case int8:
		b.V = v != 0
	case int16:
		b.V = v != 0
	case int32:
		b.V = v != 0
	case int64:
		b.V = v != 0
	case uint8:
		b.V = v != 0
	case uint16:
		b.V = v != 0
	case uint32:
		b.V = v != 0
	case uint64:
		b.V = v != 0
	case float32:
		b.V = v != 0
	case float64:
		b.V = v != 0
	case string:
		if v, ok := parseBoolean(v); ok {
			b.V = v
		} else {
			b.Valid = false
		}
	case []byte:
		if v, ok := parseBoolean(string(v)); ok {
			b.V = v
		} else {
			b.Valid = false
		}
	default:
		b.Valid = false
		return fmt.Errorf("illegal db type: %T", value)
	}

	return nil
}

// Value implements the [driver.Valuer] interface.
func (b Boolean) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.V, nil
}

func (b Boolean) String() string {
	if !b.Valid {
		return "null"
	} else if b.V {
		return "true"
	} else {
		return "false"
	}
}

func ParseBoolean(s string) (result Boolean) {
	result.V, result.Valid = parseBoolean(s)
	return
}

func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(s.V)
	}
}

func (s *String) UnmarshalJSON(data []byte) error {
	rs := string(data)
	if rs == "null" || rs == "" {
		s.Valid = false
		s.V = ""
		return nil
	} else if s0, err := strconv.Unquote(rs); err != nil {
		return err
	} else {
		s.Valid = true
		s.V = s0
		return nil
	}
}

func (s *String) Merge(o String) *String {
	if o.Valid {
		s.Valid = true
		s.V = o.V
	}

	return s
}

func (s *String) Eq(ov string) bool {
	return s.Valid && s.V == ov
}

func (s String) Coalsece(o string) string {
	if s.Valid {
		return s.V
	} else {
		return o
	}
}

// Scan implements the [Scanner] interface.
func (s *String) Scan(value any) error {
	if value == nil {
		s.V, s.Valid = "", false
		return nil
	}
	s.Valid = true
	switch v := value.(type) {
	case bool:
		s.V = strconv.FormatBool(v)
	case int8:
		s.V = strconv.FormatInt(int64(v), 10)
	case int16:
		s.V = strconv.FormatInt(int64(v), 10)
	case int32:
		s.V = strconv.FormatInt(int64(v), 10)
	case int64:
		s.V = strconv.FormatInt(int64(v), 10)
	case uint8:
		s.V = strconv.FormatUint(uint64(v), 10)
	case uint16:
		s.V = strconv.FormatUint(uint64(v), 10)
	case uint32:
		s.V = strconv.FormatUint(uint64(v), 10)
	case uint64:
		s.V = strconv.FormatUint(uint64(v), 10)
	case float32:
		s.V = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		s.V = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case string:
		s.V = v
	case []byte:
		s.V = string(v)
	default:
		s.Valid = false
		return fmt.Errorf("illegal db type: %T", value)
	}

	return nil
}

// Value implements the [driver.Valuer] interface.
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.V, nil
}

func (s String) String() string {
	if !s.Valid {
		return "null"
	} else {
		return s.V
	}
}

func ParseString(s string) (result String) {
	if s == "" {
		result.V, result.Valid = "", false
	} else {
		result.V, result.Valid = s, true
	}
	return
}

func ParseStringEnum(s string) (result String) {
	s0 := strings.ToUpper(strings.TrimSpace(s))
	if s0 == "" {
		result.V, result.Valid = "", false
	} else {
		result.V, result.Valid = s0, true
	}
	return
}

func (i Integer) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(i.V)
	}
}

func (i *Integer) UnmarshalJSON(data []byte) error {
	if v, ok := parseInt32(string(data)); ok {
		i.Valid = true
		i.V = v
	} else {
		i.Valid = false
	}

	return nil
}

func (i *Integer) Merge(o Integer) *Integer {
	if o.Valid {
		i.Valid = true
		i.V = o.V
	}

	return i
}

func (i *Integer) Eq(ov int32) bool {
	return i.Valid && i.V == ov
}

func (i Integer) Coalsece(o int32) int32 {
	if i.Valid {
		return i.V
	} else {
		return o
	}
}

func (i Integer) CoalseceInt(o int) int {
	if i.Valid {
		return int(i.V)
	} else {
		return o
	}
}

func (i *Integer) Scan(value any) error {
	if value == nil {
		i.V, i.Valid = 0, false
		return nil
	}
	i.Valid = true
	switch v := value.(type) {
	case bool:
		if v {
			i.V = 1
		} else {
			i.V = 0
		}
	case int8:
		i.V = int32(v)
	case int16:
		i.V = int32(v)
	case int32:
		i.V = v
	case int64:
		i.V = int32(v)
	case uint8:
		i.V = int32(v)
	case uint16:
		i.V = int32(v)
	case uint32:
		i.V = int32(v)
	case uint64:
		i.V = int32(v)
	case float32:
		i.V = int32(v)
	case float64:
		i.V = int32(v)
	case string:
		if v, ok := parseInt32(v); ok {
			i.V = v
		} else {
			i.Valid = false
		}
	case []byte:
		if v, ok := parseInt32(string(v)); ok {
			i.V = v
		} else {
			i.Valid = false
		}
	default:
		i.Valid = false
		return fmt.Errorf("illegal db type: %T", value)
	}

	return nil
}

func (i Integer) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.V), nil
}

func (i Integer) String() string {
	if !i.Valid {
		return "null"
	} else {
		return strconv.FormatInt(int64(i.V), 10)
	}
}

func ParseInteger(s string) (result Integer) {
	result.V, result.Valid = parseInt32(s)
	return
}

func (s Short) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(s.V)
	}
}

func (i *Short) UnmarshalJSON(data []byte) error {
	if v, ok := parseInt16(string(data)); ok {
		i.Valid = true
		i.V = v
	} else {
		i.Valid = false
	}

	return nil
}

func (s *Short) Merge(o Short) *Short {
	if o.Valid {
		s.Valid = true
		s.V = o.V
	}

	return s
}

func (s *Short) Eq(ov int16) bool {
	return s.Valid && s.V == ov
}

func (s Short) Coalsece(o int16) int16 {
	if s.Valid {
		return s.V
	} else {
		return o
	}
}

func (s *Short) Scan(value any) error {
	if value == nil {
		s.V, s.Valid = 0, false
		return nil
	}
	s.Valid = true
	switch v := value.(type) {
	case bool:
		if v {
			s.V = 1
		} else {
			s.V = 0
		}
	case int8:
		s.V = int16(v)
	case int16:
		s.V = v
	case int32:
		s.V = int16(v)
	case int64:
		s.V = int16(v)
	case uint8:
		s.V = int16(v)
	case uint16:
		s.V = int16(v)
	case uint32:
		s.V = int16(v)
	case uint64:
		s.V = int16(v)
	case float32:
		s.V = int16(v)
	case float64:
		s.V = int16(v)
	case string:
		if v, ok := parseInt16(v); ok {
			s.V = v
		} else {
			s.Valid = false
		}
	case []byte:
		if v, ok := parseInt16(string(v)); ok {
			s.V = v
		} else {
			s.Valid = false
		}
	default:
		s.Valid = false
		return fmt.Errorf("illegal db type: %T", value)
	}

	return nil
}

func (s Short) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return int64(s.V), nil
}

func (s Short) String() string {
	if !s.Valid {
		return "null"
	} else {
		return strconv.FormatInt(int64(s.V), 10)
	}
}

func ParseShort(s string) (result Short) {
	result.V, result.Valid = parseInt16(s)
	return
}

func (i Long) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(i.V)
	}
}

func (i *Long) UnmarshalJSON(data []byte) error {
	if v, ok := parseInt64(string(data)); ok {
		i.Valid = true
		i.V = v
	} else {
		i.Valid = false
	}

	return nil
}

func (l *Long) Merge(o Long) *Long {
	if o.Valid {
		l.Valid = true
		l.V = o.V
	}

	return l
}

func (l *Long) Eq(ov int64) bool {
	return l.Valid && l.V == ov
}

func (s Long) Coalsece(o int64) int64 {
	if s.Valid {
		return s.V
	} else {
		return o
	}
}

func (l *Long) Scan(value any) error {
	if value == nil {
		l.V, l.Valid = 0, false
		return nil
	}
	l.Valid = true
	switch v := value.(type) {
	case bool:
		if v {
			l.V = 1
		} else {
			l.V = 0
		}
	case int8:
		l.V = int64(v)
	case int16:
		l.V = int64(v)
	case int32:
		l.V = int64(v)
	case int64:
		l.V = v
	case uint8:
		l.V = int64(v)
	case uint16:
		l.V = int64(v)
	case uint32:
		l.V = int64(v)
	case uint64:
		l.V = int64(v)
	case float32:
		l.V = int64(v)
	case float64:
		l.V = int64(v)
	case string:
		if v, ok := parseInt64(v); ok {
			l.V = v
		} else {
			l.Valid = false
		}
	case []byte:
		if v, ok := parseInt64(string(v)); ok {
			l.V = v
		} else {
			l.Valid = false
		}
	default:
		l.Valid = false
		return fmt.Errorf("illegal db type: %T", value)
	}

	return nil
}

func (l Long) Value() (driver.Value, error) {
	if !l.Valid {
		return nil, nil
	}
	return l.V, nil
}

func (l Long) String() string {
	if !l.Valid {
		return "null"
	} else {
		return strconv.FormatInt(l.V, 10)
	}
}

func ParseLong(s string) (result Long) {
	result.V, result.Valid = parseInt64(s)
	return
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(t.V.Unix())
	}
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if s0 := string(data); s0 == `null` || s0 == `''` || s0 == `""` {
		*t = Timestamp{}
		return nil
	}

	var l int64
	if err := json.Unmarshal(data, &l); err != nil {
		return err
	} else {
		t.Valid = true
		t.V = time.Unix(l, 0)
		return nil
	}
}

func (t *Timestamp) Merge(o Timestamp) *Timestamp {
	if o.Valid {
		t.Valid = true
		t.V = o.V
	}

	return t
}

func (t *Timestamp) Eq(ov time.Time) bool {
	return t.Valid && t.V == ov
}

func (t Timestamp) Coalsece(o time.Time) time.Time {
	if t.Valid {
		return t.V
	} else {
		return o
	}
}

func (t *Timestamp) Scan(value any) error {
	if value == nil {
		t.V, t.Valid = time.Time{}, false
		return nil
	}
	t.Valid = true
	switch v := value.(type) {
	case int8:
		t.V = time.Unix(int64(v), 0)
	case int16:
		t.V = time.Unix(int64(v), 0)
	case int32:
		t.V = time.Unix(int64(v), 0)
	case int64:
		t.V = time.Unix(v, 0)
	case uint8:
		t.V = time.Unix(int64(v), 0)
	case uint16:
		t.V = time.Unix(int64(v), 0)
	case uint32:
		t.V = time.Unix(int64(v), 0)
	case uint64:
		t.V = time.Unix(int64(v), 0)
	case string:
		if v, ok := parseTime(v); ok {
			t.V = v
		} else {
			t.Valid = false
		}
	case []byte:
		if v, ok := parseTime(string(v)); ok {
			t.V = v
		} else {
			t.Valid = false
		}
	case time.Time:
		t.Valid = !v.IsZero()
		t.V = v
	default:
		t.Valid = false
		return fmt.Errorf("illegal db type: %T", value)
	}

	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.V, nil
}

func (t Timestamp) String() string {
	if !t.Valid {
		return "null"
	} else {
		return t.V.Format(time.RFC3339)
	}
}

func ParseTimestamp(s string) (result Timestamp) {
	result.V, result.Valid = parseTime(s)
	return
}
