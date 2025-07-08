package mqtt

import (
	"api/internal/constants"
	"api/internal/model"
	"api/internal/mq"
	"api/internal/utils"
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type DeviceHandler struct {
	GormDB   *gorm.DB
	Redis    *redis.Redis
	Cli      mqtt.Client
	Producer *mq.Producer
}

func (d *DeviceHandler) extractDeviceIDFromTopic(topic string) int64 {
	parts := strings.Split(topic, "/")
	if len(parts) < 3 {
		logx.Info("topic format error")
		return 0
	}
	// 例如 device/1/data/up -> parts[1] == "1"
	parseInt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0
	}
	return parseInt
}

func (d *DeviceHandler) SubscribeTopic() {
	d.subscribeConnectStatus()
	d.subscribeUploadData()
}

func (d *DeviceHandler) subscribeConnectStatus() {
	d.Cli.Subscribe(constants.TopicConnectStatus, 1, func(client mqtt.Client, message mqtt.Message) {
		topic := message.Topic()
		payload := message.Payload()

		msg := utils.Unmarshal[model.ConnectMsg](payload)
		if strings.HasSuffix(topic, "/disconnected") {
			d.handleDisconnect(msg)
		} else {
			d.handleConnect(msg)
		}
	})
}

func (d *DeviceHandler) handleConnect(msg model.ConnectMsg) {
	logx.Infof("【处理mqtt消息】handleConnect ---- %v", utils.MarshalString(msg))
}

func (d *DeviceHandler) handleDisconnect(msg model.ConnectMsg) {
	logx.Infof("【处理mqtt消息】handleDisConnect --- %v", utils.MarshalString(msg))
}

func (d *DeviceHandler) subscribeUploadData() {
	d.Cli.Subscribe(constants.TopicDataUpload, 1, func(client mqtt.Client, message mqtt.Message) {
		logx.Infof("Received mqtt message on topic: %s, payload: %v", constants.TopicDataUpload, string(message.Payload()))

		// 发送消息
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		msg, err := utils.UnmarshalWithErr[model.UploadDataMsg](message.Payload())
		if err != nil {
			logx.Errorf("unmarshal upload data err: %v", err)
			return
		}
		msg.Topic = message.Topic()

		err = d.Producer.PublishWithContext(ctx, utils.Marshal(msg))
		if err != nil {
			logx.Infof("Publish failed: %v", err)
		}
	})
}

func (d *DeviceHandler) ConsumeDataUploadMsg(msg []byte) error {
	logx.Infof("Received message: %s", msg)
	// 业务处理逻辑

	var base model.UploadDataMsg
	base, err := utils.UnmarshalWithErr[model.UploadDataMsg](msg)
	if err != nil {
		logx.Errorf("handleUploadData json unmarshal fail: %v, payload: %s", err, string(msg))
		return nil
	}

	switch base.Type {
	case "sensor":
		sensor, err := utils.UnmarshalWithErr[model.SensorData](base.Data)
		if err != nil {
			logx.Errorf("handleUploadData json unmarshal fail: %v, payload: %s", err, string(msg))
			return nil
		}
		d.handleSensor(d.extractDeviceIDFromTopic(base.Topic), base.Device, sensor)

	default:
		logx.Infof("unknown type %s", base.Type)
	}
	return nil
}

func (d *DeviceHandler) handleSensor(deviceId int64, deviceName string, sensor model.SensorData) {
	logx.Infof("【处理mqtt消息】handleSensor")
}
