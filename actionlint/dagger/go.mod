module dagger/actionlint

go 1.22.7

toolchain go1.23.3

require (
	dagger.io/dagger v0.12.4
	github.com/99designs/gqlgen v0.17.60
	github.com/Khan/genqlient v0.7.0
	github.com/vektah/gqlparser/v2 v2.5.20
	go.opentelemetry.io/otel v1.33.0
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.9.0
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.9.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.33.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.33.0
	go.opentelemetry.io/otel/log v0.9.0
	go.opentelemetry.io/otel/sdk v1.33.0
	go.opentelemetry.io/otel/sdk/log v0.9.0
	go.opentelemetry.io/otel/trace v1.33.0
	go.opentelemetry.io/proto/otlp v1.4.0
	golang.org/x/exp v0.0.0-20241217172543-b2144cdd0a67
	golang.org/x/sync v0.10.0
	google.golang.org/grpc v1.69.0
)

require go.opentelemetry.io/auto/sdk v1.1.0 // indirect

require (
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.24.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.33.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.33.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0
	go.opentelemetry.io/otel/sdk/metric v1.33.0
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241216192217-9240e9c98484 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241216192217-9240e9c98484 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc => go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.0.0-20240518090000-14441aefdf88

replace go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp => go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.3.0

replace go.opentelemetry.io/otel/log => go.opentelemetry.io/otel/log v0.3.0

replace go.opentelemetry.io/otel/sdk/log => go.opentelemetry.io/otel/sdk/log v0.3.0