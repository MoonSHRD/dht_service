package models

type DhtConfig struct {
    NPort int64 `json:"node_port"`
    SPort int64 `json:"server_port"`
    SHost string `json:"server_host"`
    PrivateKey string `json:"private_key"`
    BootstrapPeers []string `json:"bootstrap_peers"`
}