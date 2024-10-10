package redishelper

import (
	"context"
	"testing"
	"time"
)

func TestIncr(t *testing.T) {
	if err := InitRedis("localhost:6379", "", "", "", 0); err != nil {
		t.Fatalf("init redis failed, err: %v", err)
	}

	k1 := Incr(context.Background(), "test", 10*time.Second)
	t.Logf("incr result: %v", k1)

	time.Sleep(6 * time.Second)

	k2 := Incr(context.Background(), "test", 10*time.Second)
	t.Logf("incr result: %v", k2)

	time.Sleep(6 * time.Second)

	k3 := Incr(context.Background(), "test", 10*time.Second)
	t.Logf("incr result: %v", k3)
}

func TestHMSet(t *testing.T) {
	if err := InitRedis("127.0.0.1:6379", "", "", "", 0); err != nil {
		t.Fatalf("init redis failed, err: %v", err)
	}

	n1 := HashSet(context.Background(), "test", map[string]any{"a": 791, "b": "aab", "c": 3}, 10*time.Second)
	t.Logf("hset result: %v", n1)

	time.Sleep(6 * time.Second)

	n2 := HashSet(context.Background(), "test", map[string]any{"a": 781, "d": "acb"}, 10*time.Second)
	t.Logf("hset result: %v", n2)
}

func TestHMSetIfExists(t *testing.T) {
	if err := InitRedis("127.0.0.1:6379", "", "", "", 0); err != nil {
		t.Fatalf("init redis failed, err: %v", err)
	}

	n1 := HashSet(context.Background(), "test", map[string]any{"a": 791, "b": "aab", "c": 3}, 10*time.Second)
	t.Logf("hset result: %v", n1)

	n2 := HashSetIfExists(context.Background(), "test", map[string]any{"z": 25, "b": "jj"})
	t.Logf("hset xx result: %v", n2)
}

func TestHMGet(t *testing.T) {
	if err := InitRedis("127.0.0.1:6379", "", "", "", 0); err != nil {
		t.Fatalf("init redis failed, err: %v", err)
	}

	n1 := HashSet(context.Background(), "test", map[string]any{"a": 791, "b": "aab", "c": 3}, 5*time.Second)
	t.Logf("hset result: %v", n1)

	mv1 := HashGet(context.Background(), "test", "a", "c", "d")
	t.Logf("hget result: %v=>%v", "a", mv1["a"])
	t.Logf("hget result: %v=>%v", "c", mv1["c"])
	t.Logf("hget result: %v=>%v", "d", mv1["d"])

	time.Sleep(6 * time.Second)

	mv2 := HashGet(context.Background(), "test", "a", "c", "d")
	t.Logf("hget result: %v=>%v", "a", mv2["a"])
	t.Logf("hget result: %v=>%v", "c", mv2["c"])
	t.Logf("hget result: %v=>%v", "d", mv2["d"])
}

func TestDel(t *testing.T) {
	if err := InitRedis("127.0.0.1:6379", "", "", "", 0); err != nil {
		t.Fatalf("init redis failed, err: %v", err)
	}

	n1 := HashSet(context.Background(), "test", map[string]any{"a": 791, "b": "aab", "c": 3}, 5*time.Second)
	t.Logf("hset result: %v", n1)

	mv1 := HashGet(context.Background(), "test", "a", "c", "d")
	t.Logf("hget result: %v=>%v", "a", mv1["a"])
	t.Logf("hget result: %v=>%v", "c", mv1["c"])
	t.Logf("hget result: %v=>%v", "d", mv1["d"])

	n2 := Del(context.Background(), "test", "aab")
	t.Logf("del result: %v", n2)

	mv2 := HashGet(context.Background(), "test", "a", "c", "d")
	t.Logf("hget result: %v=>%v", "a", mv2["a"])
	t.Logf("hget result: %v=>%v", "c", mv2["c"])
	t.Logf("hget result: %v=>%v", "d", mv2["d"])
}
