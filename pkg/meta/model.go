package meta

type Object struct {
	Bucket      string     `json:"bucket"`
	Key         string     `json:"key"`
	Version     string     `json:"version"`
	Size        int64      `json:"size"`
	ContentType string     `json:"content_type"`
	Chunks      []ChunkRef `json:"chunks"`
}

type ChunkRef struct {
	ChunkID []byte `json:"chunk_id"`
	Index   uint32 `json:"index"`
	Stripe  uint32 `json:"stripe"`
	// 省略：节点位置信息（或由 placement 表推得）
}
