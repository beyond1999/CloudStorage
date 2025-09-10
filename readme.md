1) gateway（对外入口 / S3 兼容层）

作用：对外提供 HTTP(S) API（常见是 S3 兼容 REST）。做请求解析、鉴权、分片上传、断点续传、多段表单等，然后把元数据相关动作转给 metaserver，把数据读写转给 objnode。

典型职责：

路由与协议：/bucket, /object, Multipart Upload（你有 multiparter.go）

鉴权/签名（AK/SK、STS 可后加）

将对象的元信息（bucket、object key、版本、ETag、大小、分片列表、对象到 node 的映射等）交给 metaserver 管

根据 metaserver 的放置策略/位置映射，把实际数据流量转发到 objnode

统一日志/Tracing（你有 common/logging.go, common/tracing.go）

运行：go run ./cmd/gateway -addr :8080

2) metaserver（元数据与调度）

作用：集中存元数据 & 做放置/路由决策；在规模扩大时也是一致性与事务的核心。

典型职责：

元数据模型（你有 pkg/meta/model.go）：Bucket、Object、Part、Layout（副本/EC 条带）、生命周期规则

对象放置策略：选择哪些 objnode 承载该对象（随机、哈希、CRUSH-like、可用容量、机架感知等）

事务性更新：创建对象记录、完成分片、提交版本、删除标记（delete marker）

健康/心跳：维护各 objnode 的状态和容量

生命周期与策略（你有 pkg/lifecycle/）：过期、转储、归档等

运行：go run ./cmd/metaserver -addr :9090

3) objnode（数据节点 / 存储与纠删码）

作用：真正落盘对象数据；提供数据读写 RPC/HTTP；做校验、纠删码、后台修复等。

典型职责：

数据落地：storage.go 定义对象块/条带在本地磁盘的组织（目录布局、SSTable/Chunk/Extent）

纠删码/副本：erasure.go 做 EC（例如 k+m 编码），或简单多副本

数据面 API：流式 PUT/GET、分片写入/拼接；返回 ETag/CRC 等校验

后台任务：重建、数据清理、压缩、冷热分层

指标：I/O 延迟、吞吐、失败率、容量水位

运行：go run ./cmd/objnode -addr :7070 -data /data/obj1