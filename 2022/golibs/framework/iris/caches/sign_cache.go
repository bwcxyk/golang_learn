package caches

import (
	"fmt"
	"golibs/config"
	"golibs/db/redis"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type StorageStrategy interface {
	Set(s string) bool
}

type StorageContext interface {
	SetStrategy(strategy *StorageStrategy)
	Set(v string) bool
}

type storageContext struct {
	storageStrategy StorageStrategy
}

func (s *storageContext) SetStrategy(strategy *StorageStrategy) {
	s.storageStrategy = *strategy
}

func (s *storageContext) Set(v string) bool {
	return s.storageStrategy.Set(v)
}

type redisDb struct {
	db redis.GoLibRedis
}

type redisClusterDb struct {
	db redis.GoLibRedisCluster
}

const (
	signRoot = "sign_%s"
)

func (r *redisDb) Set(s string) bool {
	k := fmt.Sprintf(signRoot, s)
	is, err := r.db.Exists(context.Background(), k).Result()
	if err != nil {
		panic(err)
	}
	if is > 0 {
		return false
	}
	err = r.db.SetEX(context.Background(), k, 1, time.Minute*5).Err()
	if err != nil {
		panic(err)
	}
	return true
}

func (r *redisClusterDb) Set(s string) bool {
	k := fmt.Sprintf(signRoot, s)
	is, err := r.db.Exists(context.Background(), k).Result()
	if err != nil {
		panic(err)
	}
	if is > 0 {
		return false
	}
	err = r.db.SetEX(context.Background(), k, 1, time.Minute*5).Err()
	if err != nil {
		panic(err)
	}
	return true
}

var signOnce sync.Once
var signNew StorageContext

func SignSet(v string) bool {
	signOnce.Do(func() {
		var redisDbs StorageStrategy
		if config.Config.GetRedisUse() == "redis_cluster" {
			redisDbs = &redisClusterDb{redis.NewRedisCluster()}
		} else {
			redisDbs = &redisDb{redis.NewRedis()}
		}
		signNew = &storageContext{}
		signNew.SetStrategy(&redisDbs)
	})
	return signNew.Set(v)
}
