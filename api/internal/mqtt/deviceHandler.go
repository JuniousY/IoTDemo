package mqtt

import (
	"api/internal/constants"
	"api/internal/model"
	"api/internal/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"strings"
)

type DeviceHandler struct {
	GormDB *gorm.DB
	Redis  *redis.Redis
	Cli    mqtt.Client
}

func (d *DeviceHandler) SubscribeTopic() {
	d.subscribeConnectStatus()
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
	logx.Infof("【处理mqtt请求】handleConnect*** %v", utils.MarshalString(msg))
}

func (d *DeviceHandler) handleDisconnect(msg model.ConnectMsg) {
	logx.Infof("【处理mqtt请求】handleDisConnect %v", utils.MarshalString(msg))
}
