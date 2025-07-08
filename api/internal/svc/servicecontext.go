package svc

import (
	"api/internal/config"
	"api/internal/constants"
	"api/internal/model"
	"api/internal/mq"
	mqttHandler "api/internal/mqtt"
	"api/internal/repo"
	"api/internal/utils"
	"context"
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"os"
	"time"
)

type ServiceContext struct {
	Config   config.Config
	GormDB   *gorm.DB
	Redis    *redis.Redis
	RawRedis *red.Client

	DeviceHandler mqttHandler.DeviceHandler

	RabbitConn *amqp091.Connection

	//mqtt cli
	MqttCli mqtt.Client

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

	deviceHandler := mqttHandler.DeviceHandler{
		GormDB: svc.GormDB,
		Redis:  svc.Redis,
	}
	svc.DeviceHandler = deviceHandler

	svc.initRabbitMQ(c.RabbitMQ)

	svc.initCache()

	// 延迟执行 mqtt client 的 connect 操作.
	go func() {
		svc.initMQTT(c)
	}()

	return svc
}

func (svc *ServiceContext) initMQTT(c config.Config) {
	var tryTime = 5
	var cli mqtt.Client
	var err error
	for i := tryTime; i > 0; i-- {
		cli, err = initMqtt(c.Mqtt)
		if err != nil { //出现并发情况的时候可能联犀的http还没启动完毕
			logx.Errorf("mqtt 连接失败 重试剩余次数:%v", i-1)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}

	if err != nil {
		logx.Errorf("mqtt 连接失败 conf:%#v  err:%v", c.Mqtt, err)
		os.Exit(-1)
	}
	svc.MqttCli = cli
	svc.DeviceHandler.Cli = cli

	// 注册方法
	svc.DeviceHandler.SubscribeTopic()
}

func initMqtt(conf *config.MqttConf) (mc mqtt.Client, err error) {
	opts := mqtt.NewClientOptions()
	for _, broker := range conf.Brokers {
		opts.AddBroker(broker)
	}
	cliUUID := uuid.NewString()
	opts.SetClientID(conf.ClientID + "/" + cliUUID).SetUsername(conf.User).SetPassword(conf.Pass)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Info("mqtt client Connected")

	})

	opts.SetAutoReconnect(true).SetMaxReconnectInterval(30 * time.Second) //意外离线的重连参数
	opts.SetConnectRetry(true).SetConnectRetryInterval(5 * time.Second)   //首次连接的重连参数

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logx.Infof("mqtt 正在尝试连接 broker:%v", broker)
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt 连接丢失 err:%v", err)
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false {
		logx.Info("mqtt 连接失败")
		err = fmt.Errorf("mqtt 连接失败")
		return
	}
	return
}

func (svc *ServiceContext) initRabbitMQ(conf *config.RabbitMQConf) {
	conn, err := amqp091.Dial(conf.URL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
	}
	svc.RabbitConn = conn

	svc.registerRabbitMQ()
}

func (svc *ServiceContext) registerRabbitMQ() {
	conf := svc.Config.RabbitMQ
	ch, err := svc.RabbitConn.Channel()
	if err != nil {
		panic(fmt.Sprintf("create channel err: %s", err))
	}
	err = ch.ExchangeDeclare(conf.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		panic(fmt.Sprintf("create exchange err: %s", err))
	}

	// 注册并启动消费者
	ch.QueueDeclare(constants.DeviceUpQueue, true, false, false, false, nil)
	ch.QueueBind(
		constants.DeviceUpQueue,
		constants.DeviceUpRouteKey,
		conf.Exchange,
		false, nil)
	handler := svc.DeviceHandler
	consumer := &mq.Consumer{
		Channel:   ch,
		QueueName: constants.DeviceUpQueue,
		WorkerNum: 3,
		AutoAck:   false,
		Handler:   handler.ConsumeDataUploadMsg,
	}
	consumer.Start()

	// 注册生产者
	producer := &mq.Producer{
		Channel:     ch,
		Exchange:    conf.Exchange,
		RoutingKey:  constants.DeviceUpRouteKey,
		Mandatory:   false,
		Immediate:   false,
		ContentType: "application/json",
	}
	svc.DeviceHandler.Producer = producer
}

func (svc *ServiceContext) initCache() {
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
