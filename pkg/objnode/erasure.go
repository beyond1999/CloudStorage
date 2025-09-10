package objnode

import (
	"bytes"
	"github.com/klauspost/reedsolomon"
)

type EC struct {
	enc  reedsolomon.Encoder
	k, m int
}

func NewEC(k, m int) (*EC, error) {
	enc, err := reedsolomon.New(k, m)
	if err != nil {
		return nil, err
	}
	return &EC{enc: enc, k: k, m: m}, nil
}

// Split 将 data 切为 k+m 片段，每片长度相等（自动 padding）
func (e *EC) Split(data []byte) ([][]byte, error) {
	sz := (len(data) + e.k - 1) / e.k
	shards := make([][]byte, e.k+e.m)
	for i := 0; i < e.k; i++ {
		shards[i] = make([]byte, sz)
	}
	copy(bytes.Join(shards[:e.k], nil), data) // 简化：演示思路
	if err := e.enc.Encode(shards); err != nil {
		return nil, err
	}
	return shards, nil
}
