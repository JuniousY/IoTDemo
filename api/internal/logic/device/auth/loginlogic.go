package auth

import (
	"api/internal/config"
	"api/internal/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type LoginDevice struct {
	ProductID int    //产品id
	DeviceID  int64  //设备id
	ConnID    string //随机6字节字符串
	Timestamp int64  // unix时间戳
}

const ReqMaxDelay = 60 * 24 * 365 // 测试 正式上线时此处配置需要重构为在配置中心管理
//const ReqMaxDelay = 60 * 10 // 10分钟

// 设备登录认证
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 设备登录认证
// 参考文档 https://docs.emqx.com/zh/emqx/latest/access-control/authn/http.html
func (l *LoginLogic) Login(req *types.DeviceLoginReq) (resp *types.DeviceLoginResp, err error) {
	// 管理员校验
	if RootCheck(l.svcCtx.Config.AuthWhite, req.Username, req.Password, req.Ip) {
		return &types.DeviceLoginResp{
			Result:      "allow",
			IsSuperuser: true,
		}, nil
	}

	lg, err := GetLoginDevice(req.Username)
	if err != nil {
		return nil, err
	}
	device, err := l.svcCtx.DeviceCache.GetData(l.ctx, lg.DeviceID)
	if device == nil || err != nil {
		return nil, fmt.Errorf("device not found")
	}
	if device.ProductID != lg.ProductID {
		return nil, fmt.Errorf("device not found")
	}

	now := time.Now().Unix()
	if now-lg.Timestamp > ReqMaxDelay {
		return nil, fmt.Errorf("device login timeout")
	}

	checkPass := utils.HmacSha256Pass(req.Username, []byte(device.Secret))
	if checkPass != req.Password {
		return &types.DeviceLoginResp{Result: "deny"}, nil
	}

	return &types.DeviceLoginResp{Result: "allow"}, nil
}

// GetLoginDevice
// 生成 MQTT 的 username 部分, 格式为 ${productID};${deviceID};${connid};${timestamp}
func GetLoginDevice(userName string) (*LoginDevice, error) {
	keys := strings.Split(userName, ";")
	if len(keys) != 4 {
		return nil, fmt.Errorf("invalid username format, expected 4 parts but got %d", len(keys))
	}

	productID, err := strconv.Atoi(keys[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ProductID: %v", err)
	}
	deviceID, err := strconv.ParseInt(keys[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid DeviceID: %v", err)
	}
	connID := keys[2]

	timeStamp, err := strconv.ParseInt(keys[3], 10, 64)
	if err != nil {
		return nil, err
	}
	return &LoginDevice{
		ProductID: productID,
		DeviceID:  deviceID,
		ConnID:    connID,
		Timestamp: timeStamp,
	}, nil
}

// RootCheck 鉴定是否是root账号(提供给mqtt broker)
func RootCheck(a config.AuthConf, userName, password, ipaddr string) bool {
	var userCompare bool
	for _, user := range a.Users {
		if userName == user.UserName {
			userCompare = false
			if password == user.Password {
				userCompare = true
			}
			break
		}
	}
	if !userCompare {
		return false
	}
	if len(a.IpRange) == 0 {
		//如果没有,表示不开启ip白名单模式
		return true
	}
	for _, whiteIp := range a.IpRange {
		if utils.MatchIP(ipaddr, whiteIp) {
			return true
		}
	}
	return false
}
