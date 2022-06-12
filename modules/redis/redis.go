package redis

import (
	"github.com/go-redis/redis/v8"
)

type Client struct {
	//sync.RWMutex
	inner *redis.Client
	cfg   *Configuration
}
