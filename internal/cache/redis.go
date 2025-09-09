package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	rc         *redis.Client
	defaultTTL time.Duration
}

func NewDefaultCache(rc *redis.Client) *Cache {
	return &Cache{rc: rc}
}

func (c *Cache) GetDataStr(ctx context.Context, key string) (str string, err error) {
	data, err := c.rc.Get(ctx, key).Bytes()
	if err != nil {
		return
	}

	if len(data) == 0 {
		return
	}

	str = string(data)
	return
}

func (c *Cache) GetData(ctx context.Context, key string, output interface{}) (ok bool, err error) {
	data, err := c.rc.Get(ctx, key).Bytes()
	if err != nil {
		return
	}

	if len(data) == 0 {
		return
	}

	err = json.Unmarshal(data, &output)
	if err != nil {
		c.DelKey(ctx, key)
		return
	}

	ok = true
	return
}

func (c *Cache) DelKey(ctx context.Context, key string) error {
	_, err := c.rc.Del(ctx, key).Result()

	if err != nil {
		fmt.Println("del key err:", err)
		return err
	}

	return nil
}

func (c *Cache) SetKeyNx(ctx context.Context, key string, value interface{}, expiration time.Duration) bool {
	ok, err := c.rc.SetNX(ctx, key, value, expiration).Result()

	if err != nil {
		fmt.Println("set key err:", err)
		return false
	}

	return ok
}
