package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Redis     redis.RedisConf
	Mysql     *MysqlConf
	Mqtt      *MqttConf
	AuthWhite AuthConf
}

type MysqlConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type MqttConf struct {
	ClientID string   //在mqtt中的clientID
	Brokers  []string //mqtt服务器节点
	User     string   `json:",default=root"` //用户名
	Pass     string   `json:",optional"`     //密码
	ConnNum  int      `json:",default=1"`    //默认连接数
}

type AuthConf struct {
	IpRange []string `json:",optional"` //白名单ip 及ip段
	Users   []AuthUserInfo
}

type AuthUserInfo struct {
	UserName string // 管理员名称
	Password string // 密码
}
