package svc

import (
	"api/internal/config"
	"api/internal/repo"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	GormDB *gorm.DB
	Redis  *redis.Redis

	DeviceInfoRepo repo.DeviceInfoRepo
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

	return &ServiceContext{
		Config: c,
		GormDB: db,
		Redis:  redis.MustNewRedis(c.Redis),

		DeviceInfoRepo: repo.NewDeviceInfoRepo(db),
	}
}
