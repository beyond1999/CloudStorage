module CloudStorage

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.32.0 // 用于 S3 兼容签名解析（可选）
	github.com/go-chi/chi/v5 v5.1.0
	github.com/klauspost/reedsolomon v1.12.5
	github.com/prometheus/client_golang v1.20.5
	github.com/rs/zerolog v1.33.0
	go.etcd.io/etcd/client/v3 v3.5.15
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	github.com/cespare/xxhash/v2 v2.3.0
)