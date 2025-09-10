package gateway

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"

	gen "CloudStorage/pkg/api/gen"
	"github.com/go-chi/chi/v5"
)

type MetaClient interface {
	AllocatePlacement(ctx context.Context, k, m int) (placement *gen.AllocatePlacementResponse, err error)
	PutObject(ctx context.Context, bucket, key string, size int64, contentType string) (version string, err error)
	CommitObject(ctx context.Context, bucket, key string, chunks []*gen.ChunkRef) (version string, err error)
}

type ObjClient interface {
	PutChunk(ctx context.Context, addr string, chunkID []byte, idx uint32, data []byte) error
}

type Handler struct {
	R    *chi.Mux
	Meta MetaClient
	Obj  ObjClient
	K, M int // EC 参数（例如 4+2）
}

func NewHandler(meta MetaClient, obj ObjClient, k, m int) *Handler {
	h := &Handler{R: chi.NewRouter(), Meta: meta, Obj: obj, K: k, M: m}
	h.R.Put("/v1/objects/{bucket}/*", h.putObject)
	h.R.Get("/v1/objects/{bucket}/*", h.getObject) // 省略
	return h
}

func (h *Handler) putObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	cl := r.Header.Get("Content-Length")
	size, _ := strconv.ParseInt(cl, 10, 64)
	contentType := r.Header.Get("Content-Type")

	if _, err := h.Meta.PutObject(ctx, bucket, key, size, contentType); err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}
	// 读取整个对象（演示；生产建议分片/流式 + gRPC streaming + 背压）
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	sum := sha256.Sum256(body)
	chunkID := sum[:] // 演示：整个对象一个 chunk

	// 申请放置：返回 data/parity 节点
	pl, err := h.Meta.AllocatePlacement(ctx, h.K, h.M)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 简化：把同一份数据发往所有 data 节点的 index=0，parity 暂略
	for i, n := range pl.Placement.DataNodes {
		if err := h.Obj.PutChunk(ctx, n.Address, chunkID, uint32(i), body); err != nil {
			http.Error(w, "store failed:"+err.Error(), 500)
			return
		}
	}
	// 提交元数据
	ver, err := h.Meta.CommitObject(ctx, bucket, key, []*gen.ChunkRef{{Id: chunkID, Index: 0, Stripe: 0}})
	if err != nil {
		http.Error(w, err.Error(), 409)
		return
	}

	w.Header().Set("ETag", "\"sha256-"+hex.EncodeToString(sum[:8])+"\"")
	w.Header().Set("X-Version", ver)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) getObject(w http.ResponseWriter, r *http.Request) {
	// 查元数据 -> 找到 data_nodes -> 从一处/多处拉块 -> 纠删码重构 -> 回包（省略）
	w.WriteHeader(501)
}
