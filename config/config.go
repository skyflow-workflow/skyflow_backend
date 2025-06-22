package config

type Config struct {
	API         *APIConfig
	Persistence *PersistConfig
}

type PersistConfig struct {
	Database string
}
type APIConfig struct {
	Port int
	Host string
}

type ExecutorConfig struct {
	MaxConcurrency int
	MaxQueueSize   int
}

type ExporterConfig struct {
	MaxConcurrency int
}

type Limiter struct {
	Memory int
	CPU    int
}
