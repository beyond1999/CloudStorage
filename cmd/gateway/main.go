package main

import (
	"CloudStorage/pkg/gateway"
	"flag"
	"github.com/rs/zerolog/log"
	"net/http"
)

// 这里把 MetaClient/ObjClient 的具体实现注入。
func main() {
	addr := flag.String("addr", ":8080", "listen addr")
	flag.Parse()

	gw := gateway.New()
	// TODO: gw.h = gateway.NewHandler(metaClient, objClient, 4, 2)

	log.Info().Str("addr", *addr).Msg("Gateway listening")
	_ = http.ListenAndServe(*addr, gw.Handler())
}
