package model

import "time"

type Product struct {
	ID         int       `gorm:"primaryKey;autoIncrement;comment:ID"`
	Name       string    `gorm:"type:varchar(30);not null;comment:产品名称"`
	Status     int       `gorm:"type:int;default:0;not null;comment:状态 0正常，1已下线，-1删除"`
	CreateTime time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;onUpdate:CURRENT_TIMESTAMP;comment:更新时间"`
}

func (m *Product) TableName() string {
	return "product"
}

type Device struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;comment:ID"`
	ProductID  int       `gorm:"type:int;not null;index;comment:产品ID"` // 外键+索引
	Name       string    `gorm:"type:varchar(20);not null;comment:设备名称"`
	Secret     string    `gorm:"type:varchar(50);not null;comment:密码"`
	Info       string    `gorm:"type:varchar(64);default:'';not null;comment:设备信息"`
	Status     int       `gorm:"type:int;default:0;not null;comment:状态 0未激活，1已激活，-1删除"`
	IsOnline   int       `gorm:"type:int;default:0;not null;comment:是否在线 0下线 1在线"`
	CreateTime time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null"`
	UpdateTime time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;onUpdate:CURRENT_TIMESTAMP"`
}

func (m *Device) TableName() string {
	return "device"
}

// 状态常量
const (
	DeviceStatusInactive = 0
	DeviceStatusActive   = 1
	DeviceStatusDeleted  = -1
	DeviceOffline        = 0
	DeviceOnline         = 1
)
