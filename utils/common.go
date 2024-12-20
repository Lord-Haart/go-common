package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	salt = "1234xyz"
)

// Set 表示一个集合。
type Set[T comparable] struct {
	data map[T]bool
}

func (ss *Set[T]) Len() int {
	if ss == nil {
		return 0
	}
	return len(ss.data)
}

func (ss *Set[T]) IsEmpty() bool {
	if ss == nil {
		return true
	}
	return len(ss.data) == 0
}

func (ss *Set[T]) Add(t T) {
	ss.data[t] = true
}

func (ss *Set[T]) Remove(t T) {
	delete(ss.data, t)
}

func (ss *Set[T]) Contains(t T) bool {
	if ss == nil {
		return false
	}
	return ss.data[t]
}

func (ss *Set[T]) ContainsAny(t ...T) bool {
	if ss == nil {
		return false
	}
	for _, t0 := range t {
		if ss.data[t0] {
			return true
		}
	}

	return false
}

func (ss *Set[T]) String() string {
	result := make([]string, 0, ss.Len())
	for k, v := range ss.data {
		if v {
			result = append(result, fmt.Sprintf("%v", k))
		}
	}
	return "[" + strings.Join(result, ",") + "]"
}

func (ss *Set[T]) AllKeys() []T {
	result := make([]T, 0, ss.Len())
	for k, v := range ss.data {
		if v {
			result = append(result, k)
		}
	}
	return result
}

func NewSet[T comparable](t ...T) *Set[T] {
	result := &Set[T]{data: make(map[T]bool, 8)}

	for _, t0 := range t {
		result.data[t0] = true
	}

	return result
}

type AtomicTime struct {
	l sync.Mutex
	t time.Time
}

func (t *AtomicTime) CompareAndMax(oa ...time.Time) bool {
	t.l.Lock()
	defer t.l.Unlock()

	r := t.t
	result := false
	for _, o := range oa {
		if o.After(r) {
			r = o
			result = true
		}
	}

	if result {
		t.t = r
	}

	return result
}

func (t *AtomicTime) CompareAndMin(oa ...time.Time) bool {
	t.l.Lock()
	defer t.l.Unlock()

	r := t.t
	result := false
	for _, o := range oa {
		if r.IsZero() || o.Before(r) {
			r = o
			result = true
		}
	}

	if result {
		t.t = r
	}

	return result
}

func (t *AtomicTime) Time() time.Time {
	return t.t
}

func NewAtomicTime(t time.Time) *AtomicTime {
	return &AtomicTime{
		t: t,
		l: sync.Mutex{},
	}
}

type Ele[T any] struct {
	d   *T
	err error
}

func (e *Ele[T]) Data() *T {
	return e.d
}

func (e *Ele[T]) Err() error {
	return e.err
}

func NewEle[T any](t *T) Ele[T] {
	return Ele[T]{d: t, err: nil}
}

func NewEleErr[T any](err any) Ele[T] {
	if err_, ok := err.(error); ok {
		return Ele[T]{d: nil, err: err_}
	} else {
		return Ele[T]{d: nil, err: fmt.Errorf("%v", err)}
	}
}

func RecoverAsEleErr[T any](ch chan<- Ele[T]) {
	defer close(ch)

	if err := recover(); err != nil {
		ch <- NewEleErr[T](recover())
	}
}

// ParseDate 解析日期。如果解析不成功则返回默认值。
func ParseDate(s string, d time.Time) time.Time {
	if r, err := time.Parse("20060102", strings.TrimSpace(s)); err != nil {
		return d
	} else {
		return r
	}
}

// TruncateToDay 将日期截断到日。
func TruncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// CeilingToDay 将日期改为日的结束。
func CeilingToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, t.Location())
}

// Today 获取表示今天0点0分0秒的日期。
func Today() time.Time {
	return TruncateToDay(time.Now())
}

// LocalDate 获取指定的本地日期。
func LocalDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// Sha256Salt 加盐然后获取Sha256的hash值。
func Sha256Salt(s string) string {
	var ps string
	if len(s) == 0 {
		ps = "s"
	} else if len(s) <= 3 {
		ps = s[:1] + salt + s[1:]
	} else {
		ps = s[:3] + salt + s[3:] + salt
	}

	data := sha256.New().Sum([]byte(ps))
	return base64.StdEncoding.EncodeToString(data)
}

func DiffDays(st, et time.Time) float64 {
	return math.Ceil(st.Sub(et).Abs().Hours() / 24)
}

func Sum[T any, R int | int8 | int16 | int32 | int64 | float32 | float64](a []T, r0 R, fn func(T) R) R {
	for _, item := range a {
		r0 += fn(item)
	}

	return r0
}

func Compute[T comparable, R any](a map[T]R, k T, r0 R, fn func(T, R) R) R {
	var nv R
	if ov, exists := a[k]; exists {
		nv = ov
	} else {
		nv = r0
	}

	nv = fn(k, nv)
	a[k] = nv
	return nv
}

// SplitAsInt
func SplitAsInt[T int | int8 | int16 | int32 | int64](s, sep string) []T {
	tmp := strings.Split(s, sep)
	result := make([]T, 0, len(tmp))
	for _, ts := range tmp {
		ts := strings.TrimSpace(ts)
		if ts == "" {
			continue
		}
		if ti, err := strconv.ParseInt(ts, 10, 64); err != nil {
			continue
		} else {
			result = append(result, T(ti))
		}
	}
	return result
}

// SplitAndTrim
func SplitAndTrim(s, sep string) []string {
	tmp := strings.Split(s, sep)
	result := make([]string, 0, len(tmp))
	for _, ts := range tmp {
		ts := strings.TrimSpace(ts)
		if ts == "" {
			continue
		}
		result = append(result, ts)
	}
	return result
}

func StringListToUpper(sa []string) (result []string) {
	result = make([]string, 0, len(sa))
	for _, si := range sa {
		result = append(result, strings.ToUpper(si))
	}
	return
}

// ToStr 将对象转化为字符串。
// 如果o表示nil则返回ds。
func ToStr(o any, ds string) string {
	if o == nil {
		return ds
	} else if cs, ok := o.(string); ok {
		return cs
	} else {
		return fmt.Sprintf("%v", o)
	}
}

func IntListToStr[T int | int8 | int16 | int32 | int64](l []T) string {
	if len(l) == 0 {
		return ""
	} else {
		k := make([]string, 0, len(l))
		for _, v := range l {
			k = append(k, strconv.FormatInt(int64(v), 10))
		}
		return strings.Join(k, ",")
	}
}

// EncodeJson 将字符串进行json编码。
func EncodeJson(s string) string {
	return "\"" + strings.ReplaceAll(s, `"`, `\"`) + "\""
}

type ISODate time.Time

func (t ISODate) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return json.Marshal(nil)
	} else {
		return json.Marshal(time.Time(t).Format("20060102"))
	}
}

func (t *ISODate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if s == "null" || s == "" {
		*t = ISODate{}
		return nil
	} else {
		if dt, err := time.ParseInLocation("20060102", strings.TrimSpace(s), time.Local); err != nil {
			return err
		} else {
			*t = ISODate(dt)
			return nil
		}
	}
}

func (t ISODate) String() string {
	if time.Time(t).IsZero() {
		return "null"
	} else {
		return time.Time(t).Format("2006-01-02")
	}
}

func FileExists(fn string) (bool, error) {
	if _, err := os.Stat(fn); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func HasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

func TrimPrefixFold(s, prefix string) string {
	if HasPrefixFold(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func LeftStr(s string, n int) string {
	if n >= len(s) {
		return s
	} else if n <= 0 {
		return ""
	} else {
		return s[:n]
	}
}

func RightStr(s string, n int) string {
	if n >= len(s) {
		return s
	} else if n <= 0 {
		return ""
	} else {
		return s[len(s)-n:]
	}
}

func MidStr(s string, start, end int) string {
	if start < 0 {
		start = 0
	}

	if end > len(s) {
		end = len(s)
	}

	if start >= end {
		return ""
	} else {
		return s[start:end]
	}
}
