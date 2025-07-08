package constants

const (
	ShareSubTopicPrefix = "$share/group1/"
	TopicConnectStatus  = ShareSubTopicPrefix + "$SYS/brokers/+/clients/#"

	TopicDataUpload = ShareSubTopicPrefix + "device/+/data/up"
)
