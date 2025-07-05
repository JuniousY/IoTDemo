package svc

import (
	"api/internal/config"
	"api/internal/constants"
	"api/internal/model"
	"api/internal/repo"
	"api/internal/utils"
	"context"
	"fmt"
	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	GormDB   *gorm.DB
	Redis    *redis.Redis
	RawRedis *red.Client

	// 展示两种repo构建方式，实际项目里最好统一
	ProductRepo *repo.ProductRepo
	DeviceRepo  repo.DeviceRepo

	ProductCache *utils.Cache[model.Product, int]
	DeviceCache  *utils.Cache[model.Device, int64]
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.User,
		c.Mysql.Password,
		c.Mysql.Host,
		c.Mysql.Port,
		c.Mysql.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect gorm mysql: %v", err))
	}

	redisCli := redis.MustNewRedis(c.Redis)
	rawRedisCli := red.NewClient(&red.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       0, // 默认DB
	})

	svc := &ServiceContext{
		Config:   c,
		GormDB:   db,
		Redis:    redisCli,
		RawRedis: rawRedisCli,

		DeviceRepo:  repo.NewDeviceRepo(db),
		ProductRepo: repo.NewProductRepo(db),
	}
	cacheInit(svc)
	return svc
}

func cacheInit(svc *ServiceContext) {
	productCache, err := utils.NewCache(utils.CacheConfig[model.Product, int]{
		Redis:    svc.Redis,
		RawRedis: svc.RawRedis,
		KeyType:  constants.ServerCacheKeyProduct,
		GetData: func(ctx context.Context, key int) (*model.Product, error) {
			return svc.ProductRepo.FindOneByFilter(ctx, repo.ProductFilter{ID: key})
		},
	})
	if err != nil {
		panic(err)
	}
	svc.ProductCache = productCache

	deviceCache, err := utils.NewCache(utils.CacheConfig[model.Device, int64]{
		Redis:    svc.Redis,
		RawRedis: svc.RawRedis,
		KeyType:  constants.ServerCacheKeyDevice,
		GetData: func(ctx context.Context, key int64) (*model.Device, error) {
			return svc.DeviceRepo.FindOneByFilter(ctx, repo.DeviceFilter{ID: &key})
		},
	})
	if err != nil {
		panic(err)
	}
	svc.DeviceCache = deviceCache
}
