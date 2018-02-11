package BaseCache

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

type Base struct {
	Client *redis.Client
}

func NewClient(DB int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       DB,
	})
	return client
}

func (c *Base) IsNoSet() (bool, error) {
	keys, _, err := c.Client.Scan(0, "", 2).Result()
	if err != nil {
		return false, err
	}
	return len(keys) == 0, nil
}
func (c *Base) Close() {
	c.Client.Close()
}
