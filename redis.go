package gm

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Context context.Context
	client  *redis.Client
}

func InitRedis[T number](ctx context.Context, username, password, address string, db T) (redisClient *RedisClient, err error) {
	return InitRedisFromURL(ctx, fmt.Sprintf("redis://%s:%s@%s/%d", username, password, address, db))
}

func InitRedisFromURL(ctx context.Context, url string) (redisClient *RedisClient, err error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ping := client.Ping(ctx)
	if ping.Val() != "PONG" || ping.Err() != nil {
		return nil, fmt.Errorf("unable to ping redis server client, got %v", ping.Err())
	}

	return &RedisClient{Context: ctx, client: client}, nil
}

/*Base Commands*/
// ------------------------------------------------------------------------------
func (redisClient *RedisClient) Client() *redis.Client {
	return redisClient.client
}

func (redisClient *RedisClient) DBSize() *redis.IntCmd {
	return redisClient.client.DBSize(redisClient.Context)
}

func (redisClient *RedisClient) Keys(parttern string) *redis.StringSliceCmd {
	return redisClient.client.Keys(redisClient.Context, parttern)
}

func (redisClient *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redisClient.client.Set(redisClient.Context, key, value, expiration)
}

func (redisClient *RedisClient) Get(key string) *redis.StringCmd {
	return redisClient.client.Get(redisClient.Context, key)
}

func (redisClient *RedisClient) Del(keys ...string) *redis.IntCmd {
	return redisClient.client.Del(redisClient.Context, keys...)
}

// ------------------------------------------------------------------------------
func (redisClient *RedisClient) MSet(values ...interface{}) *redis.StatusCmd {
	return redisClient.client.MSet(redisClient.Context, values...)
}

func (redisClient *RedisClient) MGet(keys ...string) *redis.SliceCmd {
	return redisClient.client.MGet(redisClient.Context, keys...)
}

// ------------------------------------------------------------------------------
func (redisClient *RedisClient) HKeys(key string) *redis.StringSliceCmd {
	return redisClient.client.HKeys(redisClient.Context, key)
}

func (redisClient *RedisClient) HExists(key string, field string) *redis.BoolCmd {
	return redisClient.client.HExists(redisClient.Context, key, field)
}

func (redisClient *RedisClient) HLen(key string) *redis.IntCmd {
	return redisClient.client.HLen(redisClient.Context, key)
}

func (redisClient *RedisClient) HSet(key string, value ...interface{}) *redis.IntCmd {
	return redisClient.client.HSet(redisClient.Context, key, value...)
}

func (redisClient *RedisClient) HGetAll(key string) *redis.StringStringMapCmd {
	return redisClient.client.HGetAll(redisClient.Context, key)
}

func (redisClient *RedisClient) HGet(key string, field string) *redis.StringCmd {
	return redisClient.client.HGet(redisClient.Context, key, field)
}

func (redisClient *RedisClient) HDel(key string, field ...string) *redis.IntCmd {
	return redisClient.client.HDel(redisClient.Context, key, field...)
}

/* Redis Transaction Pipeline*/
type RedisTxPipeline struct {
	Context context.Context
	tx      redis.Pipeliner
}

func (redisClient *RedisClient) NewTxPipeline() *RedisTxPipeline {
	return &RedisTxPipeline{Context: redisClient.Context, tx: redisClient.client.TxPipeline()}
}

func (pipeline *RedisTxPipeline) Pipeliner() redis.Pipeliner {
	return pipeline.tx
}

func (pipeline *RedisTxPipeline) Keys(parttern string) *redis.StringSliceCmd {
	return pipeline.tx.Keys(pipeline.Context, parttern)
}

func (pipeline *RedisTxPipeline) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return pipeline.tx.Set(pipeline.Context, key, value, expiration)
}

func (pipeline *RedisTxPipeline) Get(key string) *redis.StringCmd {
	return pipeline.tx.Get(pipeline.Context, key)
}

func (pipeline *RedisTxPipeline) Del(keys ...string) *redis.IntCmd {
	return pipeline.tx.Del(pipeline.Context, keys...)
}

// ------------------------------------------------------------------------------
func (pipeline *RedisTxPipeline) MSet(values ...interface{}) *redis.StatusCmd {
	return pipeline.tx.MSet(pipeline.Context, values...)
}

func (pipeline *RedisTxPipeline) MGet(keys ...string) *redis.SliceCmd {
	return pipeline.tx.MGet(pipeline.Context, keys...)
}

// ------------------------------------------------------------------------------
func (pipeline *RedisTxPipeline) HKeys(key string) *redis.StringSliceCmd {
	return pipeline.tx.HKeys(pipeline.Context, key)
}

func (pipeline *RedisTxPipeline) HExists(key string, field string) *redis.BoolCmd {
	return pipeline.tx.HExists(pipeline.Context, key, field)
}

func (pipeline *RedisTxPipeline) HLen(key string) *redis.IntCmd {
	return pipeline.tx.HLen(pipeline.Context, key)
}

func (pipeline *RedisTxPipeline) HSet(key string, value ...interface{}) *redis.IntCmd {
	return pipeline.tx.HSet(pipeline.Context, key, value...)
}

func (pipeline *RedisTxPipeline) HGetAll(key string) *redis.StringStringMapCmd {
	return pipeline.tx.HGetAll(pipeline.Context, key)
}

func (pipeline *RedisTxPipeline) HGet(key string, field string) *redis.StringCmd {
	return pipeline.tx.HGet(pipeline.Context, key, field)
}

func (pipeline *RedisTxPipeline) HDel(key string, field ...string) *redis.IntCmd {
	return pipeline.tx.HDel(pipeline.Context, key, field...)
}

func (pipeline *RedisTxPipeline) Commit() ([]redis.Cmder, error) {
	return pipeline.tx.Exec(pipeline.Context)
}
