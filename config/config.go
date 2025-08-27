package config

import "github.com/mmtbak/microlibrary/rdb"

type SkyflowConfig struct {
	Persistence      *rdb.Config         `yaml:"persistence"`
	MessageQueue     *MessageQueueConfig `yaml:"message_queue"`
	DelayMesageQueue *MessageQueueConfig `yaml:"delay_message_queue"`
	API              *APIConfig          `yaml:"api"`
	Dispatcher       *DispatcherConfig   `yaml:"dispatcher"`
	Exporter         *ExporterConfig     `yaml:"exporter"`
	Resource         *ResourceLimiter    `yaml:"resource"`
}

type MessageQueueConfig struct {
	DSN string `yaml:"dsn"`
}
type APIConfig struct {
	QPSLimit int `yaml:"qps_limit"` //qps limit
}

type DispatcherConfig struct {
	MaxConcurrency int `yaml:"max_concurrency"`
	MaxQueueSize   int `yaml:"max_queue_size"`
}

type ExporterConfig struct {
	// Listeners is a list of URIs where the exporter will send workflow event data.
	// supported formats (including formats future will support):
	// for mysql "mysql://user:password@tcp(host:port)/dbname",
	// for clickhouse "clickhouse://user:password@tcp(host:port)/dbname",
	// for http "http://host:port/path"
	// for cloudevent "cloud://provider/service/resource",
	// for kafka "kafka://broker1,broker2/topic"
	// for nats "nats://host:port",
	// for redis "redis://host:port/db"
	// for gRPC "grpc://host:port",
	// for AMQP "amqp://user:password@host:port/vhost",
	// for SQS "sqs://region/queue",
	// for Pub/Sub "pubsub://project/topic",
	// for Webhook "webhook://host:port/path",
	// for File "file:///path/to/file",
	// for Local "local:///path/to/local/directory",
	// for Prometheus "prometheus://host:port/metrics",
	// for OpenTelemetry "otel://host:port",
	// for Elasticsearch "elasticsearch://host:port/index",
	// for InfluxDB "influxdb://host:port/database",
	// for MongoDB "mongodb://user:password@host:port/dbname",
	// for Redis Streams "redis-streams://host:port/stream",
	// for RabbitMQ "rabbitmq://user:password@host:port/vhost",
	// for Azure Event Hubs "azure-eventhubs://namespace/eventhub",
	// for Google Cloud Pub/Sub "google-cloud-pubsub://project/topic",
	// for AWS Kinesis "aws-kinesis://region/stream",
	// for Apache Pulsar "pulsar://host:port/topic",
	Listeners []string
}

// ResourceLimiter configuration for resource limits
type ResourceLimiter struct {
	// MemoryLimit and CPULimit are resource limits for the workflow execution.
	// k8s style format, e.g. "512Mi" for memory and "2" for CPU.
	MemoryLimit string `yaml:"memory_limit"`
	CPULimit    string `yaml:"cpu_limit"`
}

func NewConfig() *SkyflowConfig {
	return &SkyflowConfig{}
}
