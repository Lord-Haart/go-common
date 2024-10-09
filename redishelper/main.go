package redishelper

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	defaultRedisPost = 6379
)

var (
	rdb     *redis.Client
	rprefix string
)

// InitRedis 初始化Redis。
func InitRedis(addr, username, password, prefix string, db int) error {
	aa := strings.SplitN(addr, ":", 2)
	var ah, ap string
	if len(aa) == 1 {
		ah = strings.TrimSpace(aa[0])
		ap = ""
	} else {
		ah = strings.TrimSpace(aa[0])
		ap = strings.TrimSpace(aa[1])
	}
	if ah == "" {
		ah = "localhost"
	}
	if ap == "" {
		ap = fmt.Sprintf("%d", defaultRedisPost)
	}
	rdb_ := redis.NewClient(&redis.Options{
		Addr:     ah + ":" + ap,
		Username: username,
		Password: password,
		DB:       db,
	})
	if _, err := rdb_.Ping(context.Background()).Result(); err != nil {
		return err
	} else {
		rdb = rdb_
		rprefix = strings.TrimSpace(prefix)
		return nil
	}
}

func getRedisKey(key string) string {
	return rprefix + ":" + key
}

// HashSetIfExists 设置Hash，如果指定的key存在，同时保留ttl。
func HashSetIfExists(ctx context.Context, key string, mv map[string]any) bool {
	key = getRedisKey(key)

	if r0, err := rdb.TTL(ctx, key).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			panic(err)
		} else {
			// 无当前值，那么直接返回false。
			return false
		}
	} else if r0 > 1*time.Second {
		// 存在当前值，那么修改字段。
		if _, err := rdb.HSet(ctx, key, mv).Result(); err != nil {
			panic(err)
		} else {
			rdb.Expire(ctx, key, r0)
			return true
		}
	} else {
		return false
	}
}

// HashSetIfAbsent 设置Hash，如果指定的key不存在。
// 返回值 是否指定的key不存在且成功设置了值。
func HashSetIfAbsent(ctx context.Context, key string, mv map[string]any, expiration time.Duration) bool {
	if len(mv) == 0 {
		return false
	}

	key = getRedisKey(key)

	if err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		fields := make([]string, 0, len(mv))
		for rk := range mv {
			fields = append(fields, rk)
		}
		if r0, err := tx.HMGet(ctx, key, fields...).Result(); err != nil {
			if !errors.Is(err, redis.Nil) {
				return err
			} else {
				// 无当前值，那么设置当前值。
				_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
					pipe.HSet(ctx, key, mv)
					pipe.Expire(ctx, key, expiration)
					return nil
				})
				return err
			}
		} else {
			// 存在当前值。直接返回。
			for i, rv := range fields {
				mv[rv] = r0[i]
			}
			return nil
		}
	}, key); err != nil {
		if !errors.Is(err, redis.TxFailedErr) {
			panic(err)
		}

		return false
	} else {
		return true
	}
}

// HashGet 获取指定key的指定字段。
func HashGet(ctx context.Context, key string, fields ...string) map[string]any {
	if r0, err := rdb.HMGet(ctx, key, fields...).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			panic(err)
		} else {
			return nil
		}
	} else {
		result := make(map[string]any)
		for i, rv := range fields {
			result[rv] = r0[i]
		}
		return result
	}
}

// Del 删除key。
func Del(ctx context.Context, key ...string) bool {
	if n, err := rdb.Del(ctx, key...).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			panic(err)
		}
		return false
	} else {
		return n > 0
	}
}

// RedisValueToInt64 将redis返回的值转换为int64。
func RedisValueToInt64(s any, ds int64) int64 {
	if s == nil {
		return ds
	} else if rs, ok := s.(string); ok {
		if rr, err := strconv.ParseInt(strings.TrimSpace(rs), 10, 64); err != nil {
			return ds
		} else {
			return rr
		}
	} else {
		return ds
	}
}

// RedisValueToBoolean 将redis返回的值转换为bool。
func RedisValueToBoolean(s any, ds bool) bool {
	if s == nil {
		return ds
	} else if rs, ok := s.(string); ok {
		if rr, err := strconv.ParseBool(strings.TrimSpace(rs)); err != nil {
			return ds
		} else {
			return rr
		}
	} else {
		return ds
	}
}
