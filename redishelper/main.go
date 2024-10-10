package redishelper

import (
	"context"
	"errors"
	"fmt"
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
	if rprefix != "" {
		return rprefix + ":" + key
	} else {
		return key
	}
}

// HashSetIfExists 设置Hash，如果指定的key存在，同时保留ttl。
// 返回值：新增的字段数。
func HashSetIfExists(ctx context.Context, key string, mv map[string]any) int64 {
	if len(mv) == 0 {
		return 0
	}

	key = getRedisKey(key)
	argv := make([]any, 0, len(mv)*2)
	for k, v := range mv {
		argv = append(argv, k, v)
	}

	if r, err := redis.NewScript(`
	local key = KEYS[1]
	if redis.call("EXISTS", key) == 0 then
		return 0
	else
		return redis.call("HSET", key, unpack(ARGV))
	end
	`).Run(ctx, rdb, []string{key}, argv...).Int64(); err != nil {
		panic(err)
	} else {
		return r
	}
}

// HashSet 设置Hash。
// 返回值：新增的字段数。
func HashSet(ctx context.Context, key string, mv map[string]any, expiration time.Duration) int64 {
	if len(mv) == 0 {
		return 0
	}

	key = getRedisKey(key)
	argv := make([]any, 0, len(mv)*2+1)
	for k, v := range mv {
		argv = append(argv, k, v)
	}
	argv = append(argv, int64(expiration.Seconds()))

	if r, err := redis.NewScript(`
	local key = KEYS[1]
	local r = redis.call("HSET", key, unpack(ARGV, 1, #ARGV - 1))
	local t = redis.call("TTL", key)
	if t < 0 then
	  redis.call("EXPIRE", key, ARGV[#ARGV])
	end
	return r
	`).Run(ctx, rdb, []string{key}, argv...).Int64(); err != nil {
		panic(err)
	} else {
		return r
	}
}

// HashGet 获取指定key的指定字段。
func HashGet(ctx context.Context, key string, fields ...string) map[string]string {
	if len(fields) == 0 {
		return nil
	}

	key = getRedisKey(key)

	if r0, err := rdb.HMGet(ctx, key, fields...).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			panic(err)
		} else {
			return nil
		}
	} else {
		result := make(map[string]string)
		for i, rv := range fields {
			if r0[i] != nil {
				result[rv] = r0[i].(string)
			}
		}
		return result
	}
}

// Incr 自增指定的键，并指定过期时间。
func Incr(ctx context.Context, key string, expiration time.Duration) int64 {
	if r, err := redis.NewScript(`
	local key = KEYS[1]
	local r = redis.call("INCRBY", key, 1)
	local t = redis.call("TTL", key)
	if t < 0 then
	  redis.call("EXPIRE", key, ARGV[1])
	end
	return r
	`).Run(ctx, rdb, []string{getRedisKey(key)}, int64(expiration.Seconds())).Result(); err != nil {
		panic(err)
	} else {
		return r.(int64)
	}
}

// Del 删除key。
func Del(ctx context.Context, key ...string) int64 {
	key2 := make([]string, 0, len(key))
	for _, kk := range key {
		key2 = append(key2, getRedisKey(kk))
	}

	if n, err := rdb.Del(ctx, key2...).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			panic(err)
		}
		return 0
	} else {
		return n
	}
}
