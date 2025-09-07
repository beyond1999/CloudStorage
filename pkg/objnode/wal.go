package objnode

import (
	"os"
)

type WAL struct{ f *os.File }

func OpenWAL(dir string) (*WAL, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(dir+"/wal.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &WAL{f: f}, nil
}

func (w *WAL) Append(rec []byte) error {
	if _, err := w.f.Write(rec); err != nil {
		return err
	}
	if _, err := w.f.Write([]byte("\n")); err != nil {
		return err
	}
	return w.f.Sync()
}

func (w *WAL) Close() error { return w.f.Close() }
