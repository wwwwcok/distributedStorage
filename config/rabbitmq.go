package config

const (
	AsyncTransferEnable  = true
	RabbitURL            = "amqp://guest:guest@192.168.0.90:5672/"
	TransExchangeName    = "uploadserver.trans"
	TransOSSQueueName    = "uploadserver.trans.oss"
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	TransOSSRoutingKey   = "oss"
)
