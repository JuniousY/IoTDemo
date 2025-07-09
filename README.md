## 项目介绍
相关技术：
- 服务框架：go-zero
- 关系型数据库：Mysql
- 缓存：Redis
- MQTT 消息服务器：emqx
- 消息队列：RabbitMQ
- 部署：docker-compose

## 功能点
- 设备管理
    - [x] 设备信息维护
    - [x] 设备认证
    - [x] 设备上报数据
- 通用cache层

技术细节：
- 通过`MQTT`协议实现设备与服务端的连接。`MQTT`设备认证常见两种方式，一种为`证书认证`，一种为`密钥认证`，本项目目前支持密钥认证。
- 通用cache层包含本地缓存、Redis缓存的读写，如果缓存中不存在，则从关系型数据库中获取

## 部署
部署流程：
1. 部署`MySQL`后跑`./deploy/migration/init.sql`
2. 部署`Reids`
3. 修改配置 `./api/etc/api.yaml`
4. 构建镜像 `docker build -f Dockerfile -t iot-demo .`
5. 进入到`./deploy/docker`文件夹中进行部署：`docker compose up -d`

测试设备连接：
- 使用`MQTTX｀进行设备连接测试
- `username`、`password`: 在`auth_test.go`文件中手动生成
- 设备上传`topic`为`device/$deviceId/data/up`
