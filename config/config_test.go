package config

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mmtbak/microlibrary/rdb"
	"gopkg.in/yaml.v3"
)

func TestValidateConfig(t *testing.T) {

	var testcases = []struct {
		name     string
		template string
		config   *SkyflowConfig
		expected bool
	}{
		{
			name: "valid config",
			template: `
persistence:
  dsn: mysql://user:password@tcp(localhost:3306)/skyflow
  max_idle_conns: 10
  max_open_conns: 10
  log_level: info
  max_idle_time: 10m

message_queue:
  dsn: kafka://localhost:9092/?topics=skyflow-events&numpartition=3&numreplica=1&autocommitsecond=1&initial=oldest&version=1.1.1
delay_message_queue:
  dsn: redis://host:port/db
api:
  qps_limit: 1000
dispatcher:
  max_concurrency: 10
  max_queue_size: 100
exporter:
  listeners:
    - "http://localhost:8080/export"
resource:
  memory_request: 512Mi
  cpu_request: 2
  rate: 1000`,
			config: &SkyflowConfig{
				Persistence: &rdb.Config{
					DSN:          "mysql://user:password@tcp(localhost:3306)/skyflow",
					MaxOpenConns: 10,
					MaxIdleConns: 10,
					LogLevel:     "info",
					MaxIdleTime:  "10m",
				},
				MessageQueue: &MessageQueueConfig{
					DSN: "kafka://localhost:9092/?topics=skyflow-events&numpartition=3&numreplica=1&autocommitsecond=1&initial=oldest&version=1.1.1",
				},
				DelayMesageQueue: &MessageQueueConfig{
					DSN: "redis://host:port/db",
				},
				API: &APIConfig{
					QPSLimit: 1000,
				},
				Dispatcher: &DispatcherConfig{
					MaxConcurrency: 10,
					MaxQueueSize:   100,
				},
				Exporter: &ExporterConfig{
					Listeners: []string{"http://localhost:8080/export"},
				},
				Resource: &ResourceLimiter{
					MemoryLimit: "512Mi",
					CPULimit:    "2",
				},
			},
			expected: true,
		},
		{
			name: "invalid config - missing persistence",
			template: `
message_queue:
  dsn: nats://localhost:4222
delay_message_queue:
  dsn: nats://localhost:4222
api:
  qps_limit: 1000
dispatcher:
  max_concurrency: 10
  max_queue_size: 100
exporter:
  listeners:
    - "http://localhost:8080/export"
limiter:
  memory: 512
  cpu: 2
  rate: 1000`,
			config: &SkyflowConfig{
				Persistence: nil,
				MessageQueue: &MessageQueueConfig{
					DSN: "nats://localhost:4222",
				},
				DelayMesageQueue: &MessageQueueConfig{
					DSN: "nats://localhost:4222",
				},
				API: &APIConfig{
					QPSLimit: 1000,
				},
				Dispatcher: &DispatcherConfig{
					MaxConcurrency: 10,
					MaxQueueSize:   100,
				},
				Exporter: &ExporterConfig{
					Listeners: []string{"http://localhost:8080/export"},
				},
				Resource: &ResourceLimiter{
					MemoryLimit: "512Mi",
					CPULimit:    "2",
				},
			},
			expected: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var newConfig = NewConfig()
			err := yaml.Unmarshal([]byte(tc.template), &newConfig)
			if err != nil {
				t.Fatalf("Failed to unmarshal template: %v", err)
			}
			assert.Equal(t, err == nil, tc.expected)
			assert.Equal(t, newConfig.Persistence, tc.config.Persistence)
			assert.Equal(t, newConfig.MessageQueue.DSN, tc.config.MessageQueue.DSN)
			assert.Equal(t, newConfig.DelayMesageQueue.DSN, tc.config.DelayMesageQueue.DSN)
			assert.Equal(t, newConfig.API.QPSLimit, tc.config.API.QPSLimit)
			assert.Equal(t, newConfig.Dispatcher.MaxConcurrency, tc.config.Dispatcher.MaxConcurrency)
			assert.Equal(t, newConfig.Dispatcher.MaxQueueSize, tc.config.Dispatcher.MaxQueueSize)
		})
	}

}
