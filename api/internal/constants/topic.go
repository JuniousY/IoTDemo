package constants

const (
	ShareSubTopicPrefix = "$queue/"
	TopicConnectStatus  = ShareSubTopicPrefix + "$SYS/brokers/+/clients/#"

	TopicExt = ShareSubTopicPrefix + "$ext/up/#"
)
