package auth

import (
	"api/internal/utils"
	"fmt"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	const productId = 1
	const deviceId = 1
	connId, _ := utils.SecureRandomString(6)
	timeStamp := time.Now().Unix()
	const secret = "RAbIee3mldOc_kdV0vmj"

	ld := &LoginDevice{
		ProductID: productId,
		DeviceID:  deviceId,
		ConnID:    connId,
		Timestamp: timeStamp,
	}
	username := fmt.Sprintf("%d;%d;%s;%d", ld.ProductID, ld.DeviceID, ld.ConnID, ld.Timestamp)

	password := utils.HmacSha256Pass(username, []byte(secret))
	fmt.Printf("username: %s\npassword: %s\n", username, password)
}
