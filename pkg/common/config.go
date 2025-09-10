package common

type GatewayConfig struct {
	Addr    string
	MetaRPC string
}

type MetaConfig struct {
	Etcd []string
	Addr string
}

type ObjNodeConfig struct{ ID, Addr, DataDir string }
