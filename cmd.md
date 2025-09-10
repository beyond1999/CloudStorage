

## 1) 生成 gRPC 代码（proto 没跑就会 “找不到包”）

安装生成器（只需一次）：

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# 确保 $GOPATH/bin 在 PATH（Windows 把 %USERPROFILE%\go\bin 加到环境变量）
```

在项目根目录执行（按目录结构）：

```bash
protoc -I pkg/api/proto pkg/api/proto/objectstore.proto --go_out=pkg/api/gen --go_opt=paths=source_relative --go-grpc_out=pkg/api/gen --go-grpc_opt=paths=source_relative
```


运行后应出现：

```
pkg/api/gen/objectstore.pb.go
pkg/api/gen/objectstore_grpc.pb.go
```

最后：

```bash
go mod tidy
```

---

## 2) 常见坑对照表

* ❌ `import "titan/objstore/pkg/api/gen"` 与 `go.mod` 的 `module` 不一致 → **改 import 或改 go\_package，再生成**
* ❌ 生成到错目录（比如没进 `pkg/api/gen`）→ **检查 protoc 命令的输出路径**
* ❌ Windows 环境 `protoc-gen-go` 不在 PATH → **把 `%USERPROFILE%\go\bin` 加到 PATH**
* ❌ `option go_package` 与实际 import 不一致 → **两边统一后重生成**

---

## 3) 快速验证

```bash
go build ./cmd/objnode
go build ./cmd/metaserver
go build ./cmd/gateway
```

都能过就说明路径与生成没问题了。

