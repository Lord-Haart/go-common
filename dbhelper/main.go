package dbhelper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DbRow interface {
	Scan(dest ...any) error
}

// ctxKey 用于在上下文中记录事务对象的key。
type ctxKey struct{}

type ctxRef struct {
	tx    *sql.Tx
	alive bool
}

func (cr *ctxRef) Commit() error {
	if cr.alive {
		if err := cr.tx.Commit(); err != nil {
			return err
		} else {
			cr.alive = false
			log.Printf("[DEBUG] Committed transaction\n")
			return nil
		}
	} else {
		return nil
	}
}

func (cr *ctxRef) Close() error {
	if cr.alive {
		if err := cr.tx.Rollback(); err != nil {
			return err
		} else {
			cr.alive = false
			log.Printf("[DEBUG] Rollback transaction\n")
			return nil
		}
	} else {
		return nil
	}
}

var (
	db              *sql.DB
	SQL_ARG_PATTERN = regexp.MustCompile(`:[1|2|3|4|5|6|7|8|9](0|1|2|3|4|5|6|7|8|9)?`)
)

// InitMySqlDb 初始化MySql数据源。
func InitMySqlDb(addr, username, password, dbname string) error {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = addr
	cfg.User = username
	cfg.Passwd = password
	cfg.DBName = dbname
	cfg.Loc = time.Local
	cfg.ParseTime = true
	if db_, err := sql.Open("mysql", cfg.FormatDSN()); err != nil {
		return err
	} else {
		db_.SetConnMaxLifetime(300 * time.Second)
		db_.SetMaxIdleConns(2)
		if err := db_.Ping(); err != nil {
			return err
		}

		db = db_
		return nil
	}
}

func logSql(query string, args []any) {
	buf := make([]string, 0, 10)

	buf = append(buf, query)
	for i, arg := range args {
		av := any(nil)
		if ns, ok := arg.(sql.NullString); ok {
			if ns.Valid {
				av = ns.String
			}
		} else if ni, ok := arg.(sql.NullInt64); ok {
			if ni.Valid {
				av = ni.Int64
			}
		} else if ni, ok := arg.(sql.NullInt32); ok {
			if ni.Valid {
				av = ni.Int32
			}
		} else if ni, ok := arg.(sql.NullInt16); ok {
			if ni.Valid {
				av = ni.Int16
			}
		} else if ni, ok := arg.(sql.NullByte); ok {
			if ni.Valid {
				av = ni.Byte
			}
		} else if ni, ok := arg.(sql.NullFloat64); ok {
			if ni.Valid {
				av = ni.Float64
			}
		} else if ni, ok := arg.(sql.NullBool); ok {
			if ni.Valid {
				av = ni.Bool
			}
		} else if ni, ok := arg.(sql.NullTime); ok {
			if ni.Valid {
				av = ni.Time
			}
		} else {
			av = arg
		}

		if t, ok := av.(time.Time); ok {
			buf = append(buf, fmt.Sprintf("  [%d] %s", i+1, t.Format("2006-01-02T15:04:05-0700")))
		} else {
			buf = append(buf, fmt.Sprintf("  [%d] %#v", i+1, av))
		}
	}

	log.Printf("[DEBUG] Execute sql: %s\n", strings.Join(buf, "\n"))
}

func prepareSql(ctx context.Context, query string, args []any) (*sql.Stmt, []any) {
	oargs := make([]any, 0, len(args))
	oquery := SQL_ARG_PATTERN.ReplaceAllStringFunc(query, func(s string) string {
		s0 := s[1:]
		if v, err := strconv.ParseInt(s0, 10, 64); err != nil {
			panic(fmt.Errorf("illegal sql arg: %#v", s0))
		} else {
			si := int(v)
			if si <= len(args) {
				oargs = append(oargs, args[si-1])
				return "?"
			} else {
				panic(fmt.Errorf("no enough args, wants %d", si))
			}
		}
	})

	logSql(oquery, oargs)

	if tx, ok := ctx.Value(ctxKey{}).(*sql.Tx); ok {
		if stmt, err := tx.Prepare(oquery); err != nil {
			panic(err)
		} else {
			return stmt, oargs
		}
	} else {
		if stmt, err := db.Prepare(oquery); err != nil {
			panic(err)
		} else {
			return stmt, oargs
		}
	}
}

func BeginTx(ctx context.Context, serializable bool) context.Context {
	isolation := sql.LevelDefault
	if serializable {
		isolation = sql.LevelSerializable
	}

	if tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: isolation}); err != nil {
		panic(err)
	} else {
		log.Printf("[DEBUG] Begin transaction\n")
		return context.WithValue(ctx, ctxKey{}, &ctxRef{tx: tx, alive: true})
	}
}

func CommitTx(ctx context.Context) {
	if cr, ok := ctx.Value(ctxKey{}).(*ctxRef); ok {
		if err := cr.Commit(); err != nil {
			panic(err)
		}
	}
}

func CloseTx(ctx context.Context) {
	if cr, ok := ctx.Value(ctxKey{}).(*ctxRef); ok {
		cr.Close()
	}
}

// Exec 执行指定的sql并返回受影响的行数。
func Exec[T int | int8 | int16 | int32 | int64](ctx context.Context, query string, args ...any) (T, error) {
	stmt, args := prepareSql(ctx, query, args)

	defer stmt.Close()

	if r, err := stmt.ExecContext(ctx, args...); err != nil {
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1452 /*违反外键约束的错误*/ {
			return 0, nil
		} else {
			return 0, err
		}
	} else {
		result, err := r.RowsAffected()
		return T(result), err
	}
}

// ExecLastInsertId 执行指定的sql并返回插入的ID值。只要sql执行成功，即使未插入任何记录也不会返回错误。
func ExecLastInsertId[T int | int8 | int16 | int32 | int64](ctx context.Context, query string, args ...any) (T, error) {
	stmt, args := prepareSql(ctx, query, args)

	defer stmt.Close()

	if r, err := stmt.ExecContext(ctx, args...); err != nil {
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == 1452 /*违反外键约束的错误*/ {
			return 0, nil
		} else {
			return 0, err
		}
	} else {
		if result, err := r.LastInsertId(); err != nil {
			return 0, nil
		} else {
			return T(result), nil
		}
	}
}

func Query[T bool | int | int8 | int16 | int32 | int64 | string | time.Time | sql.NullInt64 | sql.NullBool | sql.NullTime](ctx context.Context, query string, args ...any) (T, error) {
	stmt, args := prepareSql(ctx, query, args)

	defer stmt.Close()

	r := stmt.QueryRowContext(ctx, args...)

	var result T
	if err := r.Scan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		} else {
			return result, err
		}
	} else {
		return result, err
	}
}

type RowHandler[T any] interface {
	Scan(sc DbRow) (*T, error)
}

func QueryObj[T any](ctx context.Context, query string, rh RowHandler[T], args ...any) (*T, error) {
	stmt, args := prepareSql(ctx, query, args)

	defer stmt.Close()

	r := stmt.QueryRowContext(ctx, args...)

	if result, err := rh.Scan(r); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return result, err
	}
}

func QueryList[T bool | int | int8 | int16 | int32 | int64 | string | time.Time | sql.NullInt64 | sql.NullBool | sql.NullTime](ctx context.Context, query string, args ...any) ([]T, error) {
	stmt, args := prepareSql(ctx, query, args)

	defer stmt.Close()

	if r, err := stmt.QueryContext(ctx, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		result := make([]T, 0, 10)
		for r.Next() {
			var item T
			if err := r.Scan(&item); err != nil {
				return result, err
			} else {
				result = append(result, item)
			}
		}

		return result, nil
	}
}

// QueryObjList 查询对象列表。
func QueryObjList[T any](ctx context.Context, query string, rh RowHandler[T], args ...any) ([]*T, error) {
	stmt, args := prepareSql(ctx, query, args)

	defer stmt.Close()

	if r, err := stmt.QueryContext(ctx, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		result := make([]*T, 0, 10)
		for r.Next() {
			if item, err := rh.Scan(r); err != nil {
				return result, err
			} else {
				result = append(result, item)
			}
		}

		return result, nil
	}
}

func QueryKeyValueMap[T bool | int | int8 | int16 | int32 | int64 | string](ctx context.Context, query string, args ...any) (map[string]T, error) {
	if result, err := QueryObjList[KeyValuePairPo[T]](ctx, query, &KeyValuePairMapper[T]{}, args...); err != nil {
		return nil, err
	} else {
		m := make(map[string]T, len(result))
		for _, kv := range result {
			m[kv.Key] = kv.Value
		}
		return m, nil
	}
}

// insertBatch 批量插入记录。
// 返回成功插入的记录数。
func InsertBatch(ctx context.Context, query string, rows ...[]any) (int64, error) {
	if tx, err := db.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return 0, err
	} else if stmt, err := tx.Prepare(query); err != nil {
		return 0, err
	} else {
		c := int64(0)
		re := error(nil)

		for _, row := range rows {
			if row == nil {
				continue
			}

			if r, err := stmt.ExecContext(ctx, row...); err != nil {
				if re != nil {
					re = err
				}
			} else if c_, err := r.RowsAffected(); err != nil {
				panic(err)
			} else {
				c += c_
			}
		}

		return c, re
	}
}

func JoinInString(args []string) string {
	if len(args) == 0 {
		return "('')"
	} else {
		sa := make([]string, 0, len(args))
		for _, a := range args {
			sa = append(sa, "'"+a+"'")
		}
		return "(" + strings.Join(sa, ",") + ")"
	}
}

func JoinInInt[T int | int8 | int16 | int32 | int64](args []T) string {
	if len(args) == 0 {
		return "(0)"
	} else {
		sa := make([]string, 0, len(args))
		for _, a := range args {
			sa = append(sa, strconv.FormatInt(int64(a), 10))
		}
		return "(" + strings.Join(sa, ",") + ")"
	}
}

type KeyValuePairPo[T bool | int | int8 | int16 | int32 | int64 | string] struct {
	Key   string
	Value T
}

type KeyValuePairMapper[T bool | int | int8 | int16 | int32 | int64 | string] struct{}

func (m *KeyValuePairMapper[T]) Scan(r DbRow) (*KeyValuePairPo[T], error) {
	result := &KeyValuePairPo[T]{}
	err := r.Scan(&result.Key, &result.Value)
	return result, err
}
