package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/maypok86/otter/v2"
	rd "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/sync/singleflight"
	"math/rand"
	"sync"
	"time"
)

type Cache[dataT any, keyT comparable] struct {
	keyType       string
	expireTime    time.Duration
	localCache    *otter.Cache[string, *dataT]
	redis         *redis.Redis
	rawRedis      *rd.Client
	getData       func(ctx context.Context, key keyT) (*dataT, error)
	pubSubChannel string
	sf            singleflight.Group
}

type CacheConfig[dataT any, keyT comparable] struct {
	KeyType    string
	GetData    func(ctx context.Context, key keyT) (*dataT, error)
	ExpireTime time.Duration
	Redis      *redis.Redis
	RawRedis   *rd.Client
}

var (
	cacheMap   = map[string]any{}
	cacheMutex sync.Mutex
)

func NewCache[dataT any, keyT comparable](cfg CacheConfig[dataT, keyT]) (*Cache[dataT, keyT], error) {
	cacheMutex.Lock() //单例模式
	defer cacheMutex.Unlock()

	if v, ok := cacheMap[cfg.KeyType]; ok {
		return v.(*Cache[dataT, keyT]), nil
	}

	localCache := otter.Must(&otter.Options[string, *dataT]{
		MaximumSize:     1000,
		InitialCapacity: 100,
	})

	ret := Cache[dataT, keyT]{
		keyType:       cfg.KeyType,
		localCache:    localCache,
		getData:       cfg.GetData,
		expireTime:    cfg.ExpireTime,
		redis:         cfg.Redis,
		rawRedis:      cfg.RawRedis,
		pubSubChannel: GenPubSubChannel(cfg.KeyType),
	}
	if ret.expireTime == 0 {
		ret.expireTime = time.Minute*10 + time.Second*time.Duration(rand.Int63n(60))
	}

	cacheMap[cfg.KeyType] = &ret

	ctx := context.Background()
	ret.StartPubSub(ctx)
	return &ret, nil
}

func (c *Cache[dataT, keyT]) StartPubSub(ctx context.Context) {

	if c.pubSubChannel == "" {
		return
	}
	go func() {
		sub := c.rawRedis.Subscribe(ctx, c.pubSubChannel)
		defer sub.Close()

		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					logx.Error("pubsub channel closed")
					return
				}
				// 监听到后删除本地缓存
				logx.Infof("PubSub: invalidate local cache for key: %s", msg.Payload)
				c.localCache.Invalidate(msg.Payload)
			}
		}
	}()
}

func (c *Cache[dataT, keyT]) GenCacheKey(key any) string {
	return fmt.Sprintf("cache:%s:%v", c.keyType, key)
}

func GenKeyStr(v any) string {
	return fmt.Sprintf("%v", v)
}

func GenPubSubChannel(keyType string) string {
	return "cache_channel:" + keyType
}

func (c *Cache[dataT, keyT]) SetData(ctx context.Context, key keyT, data *dataT) error {
	cacheKey := c.GenCacheKey(key)
	keyStr := GenKeyStr(key)

	if data != nil {
		dataStr, err := json.Marshal(data)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		err = c.redis.SetexCtx(ctx, cacheKey, string(dataStr), int(c.expireTime.Seconds()))
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	} else {
		// 短时间缓存空值 防止缓存穿透
		c.localCache.Set(keyStr, nil)
		err := c.redis.SetexCtx(ctx, cacheKey, "NULL", 10)
		if err != nil {
			return err
		}
		return nil
	}

	// 删除本地缓存，通知删除其他节点的本地缓存
	c.localCache.Invalidate(keyStr)
	if c.pubSubChannel != "" {
		go func() {
			_, _ = c.redis.PublishCtx(ctx, c.pubSubChannel, keyStr)
		}()
	}

	return nil
}

func (c *Cache[dataT, keyT]) GetData(ctx context.Context, key keyT) (*dataT, error) {
	cacheKey := c.GenCacheKey(key)
	keyStr := GenKeyStr(key)
	temp, _ := c.localCache.GetIfPresent(keyStr)
	if temp != nil {
		return temp, nil
	}

	//并发获取的情况下避免击穿
	ret, err, _ := c.sf.Do(keyStr, func() (any, error) {
		{ //内存中没有就从redis上获取
			val, err := c.redis.GetCtx(ctx, cacheKey)
			if len(val) > 0 {
				var ret dataT
				err = json.Unmarshal([]byte(val), &ret)
				if err != nil {
					return nil, err
				}
				c.localCache.Set(keyStr, &ret)
				return &ret, nil
			}
		}
		if c.getData == nil {
			// 如果没有设置从数据库读，则直接设置该参数为空并返回
			c.localCache.Set(keyStr, nil)
			return nil, nil
		}
		//redis上没有就读数据库
		data, err := c.getData(ctx, key)
		if err != nil {
			return nil, err
		}
		//读到了之后设置缓存
		var newData = data
		c.localCache.Set(keyStr, newData)
		if data == nil {
			return data, err
		}
		dataStr, err := json.Marshal(data)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		} else {
			_, err = c.redis.SetnxExCtx(ctx, cacheKey, string(dataStr), int(c.expireTime.Seconds()))
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
		}

		return newData, nil
	})
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*dataT), nil
}
