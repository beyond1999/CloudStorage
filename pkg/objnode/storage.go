package objnode

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Storage struct {
	base string
	wal  *WAL
}

func OpenStorage(base string) (*Storage, error) {
	if err := os.MkdirAll(base, 0o755); err != nil {
		return nil, err
	}
	wal, err := OpenWAL(filepath.Join(base, "wal"))
	if err != nil {
		return nil, err
	}
	return &Storage{base: base, wal: wal}, nil
}

func (s *Storage) chunkPath(id []byte, idx uint32) string {
	h := hex.EncodeToString(id)
	return filepath.Join(s.base, "chunks", h[:2], h[2:4], fmt.Sprintf("%s.%d", h, idx))
}

func (s *Storage) PutChunk(id []byte, idx uint32, data []byte) error {
	p := s.chunkPath(id, idx)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	// WAL: 记录将要写入的路径 + 校验
	sum := sha256.Sum256(data)
	rec := fmt.Sprintf("PUT %x %d %x", id, idx, sum)
	if err := s.wal.Append([]byte(rec)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(p+".tmp", data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(p+".tmp", p); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetChunk(id []byte, idx uint32) ([]byte, error) {
	p := s.chunkPath(id, idx)
	return os.ReadFile(p)
}
