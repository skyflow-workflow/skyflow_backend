module github.com/skyflow-workflow/skyflow_backbend

go 1.22

toolchain go1.22.9

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/go-playground/assert/v2 v2.2.0
	github.com/go-playground/validator/v10 v10.23.0
	github.com/google/go-cmp v0.6.0
	github.com/json-iterator/go v1.1.12
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mmtbak/microlibrary v0.0.0-20250402133427-d63d0f3e4552
	github.com/ohler55/ojg v1.26.1
	github.com/skyflow-workflow/skyflow_backbend/gen/pb v0.0.0-20241109094122-6caef8ab07e4
	google.golang.org/protobuf v1.36.6
	gopkg.in/go-playground/assert.v1 v1.2.1
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.5
	trpc.group/trpc-go/trpc-go v1.0.3
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/ClickHouse/ch-go v0.58.2 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.14.3 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.6.1 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v2.0.0+incompatible // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/gorilla/schema v1.2.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kos-v/dsnparser v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/panjf2000/ants/v2 v2.8.1 // indirect
	github.com/paulmach/orb v0.10.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.43.0 // indirect
	go.opentelemetry.io/otel v1.19.0 // indirect
	go.opentelemetry.io/otel/trace v1.19.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/automaxprocs v1.3.0 // indirect
	go.uber.org/mock v0.5.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.25.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/clickhouse v0.5.1 // indirect
	trpc.group/trpc-go/tnet v1.0.1 // indirect
	trpc.group/trpc/trpc-protocol/pb/go/trpc v1.0.0 // indirect
)

replace github.com/skyflow-workflow/skyflow_backbend/gen/pb => ./gen/pb

replace github.com/kos-v/dsnparser => github.com/mmtbak/dsnparser v0.0.0-20240220012319-a0bde160948e
