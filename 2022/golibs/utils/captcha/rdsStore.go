package captcha

import (
	"golibs/db/redis"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/net/context"
	"time"
)

type rdsStore struct {
	expiration time.Duration
}

var ctx = context.Background()

func (r rdsStore) Set(id string, value string) {
	if err := redis.NewRedis().Set(ctx, id, value, r.expiration).Err(); err != nil {
		panic(err)
	}
}

func (r rdsStore) Get(id string, clear bool) string {
	code, err := redis.NewRedis().Get(ctx, id).Result()
	if err != nil {
		return ""
	}
	if clear {
		redis.NewRedis().Del(ctx, id)
	}
	return code
}

func (r rdsStore) Verify(id, answer string, clear bool) bool {
	return r.Get(id, clear) == answer
}

func newRdsStore(expiration time.Duration) base64Captcha.Store {
	s := new(rdsStore)
	s.expiration = expiration
	return s
}

var rdsMemStore = newRdsStore(time.Minute)
