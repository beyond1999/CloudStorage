package main

import (
	"CloudStorage/pkg/objnode"
	"flag"
	"github.com/rs/zerolog/log"
)

func main() {
	id := flag.String("id", "obj1", "node id")
	addr := flag.String("addr", ":7001", "listen addr")
	data := flag.String("data", "./data", "data dir")
	flag.Parse()

	st, err := objnode.OpenStorage(*data)
	if err != nil {
		log.Fatal().Err(err).Msg("open storage")
	}
	if err := objnode.Serve(*addr, st); err != nil {
		log.Fatal().Err(err).Msg("serve")
	}
	_ = id // 这里应把节点注册到 etcd /titan/nodes 下，省略
}
