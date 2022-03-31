package configscraper

var knownParams = map[string]*Paramaeter{
	"endpoint": {
		Name:     "endpoint",
		Type:     StringType,
		Required: true,
	},
	"tls": {
		Name:     "tls",
		Type:     MapType,
		Required: false,
	},
	"read_buffer_size": {
		Name:     "read_buffer_size",
		Type:     IntType,
		Required: false,
	},
	"write_buffer_size": {
		Name:     "write_buffer_size",
		Type:     IntType,
		Required: false,
	},
	"timeout": {
		Name:     "timeout",
		Type:     DurationType,
		Required: false,
	},
	"headers": {
		Name:     "headers",
		Type:     MapType,
		Required: false,
	},
	"auth": {
		Name:     "auth",
		Type:     MapType,
		Required: false,
	},
	"compression": {
		Name:     "compression",
		Type:     StringType,
		Required: false,
	},
	"max_idle_conns": {
		Name:     "max_idle_conns",
		Type:     IntType,
		Required: false,
	},
	"max_idle_conns_per_host": {
		Name:     "max_idle_conns_per_host",
		Type:     IntType,
		Required: false,
	},
	"max_conns_per_host": {
		Name:     "max_conns_per_host",
		Type:     IntType,
		Required: false,
	},
	"idle_conn_timeout": {
		Name:     "idle_conn_timeout",
		Type:     DurationType,
		Required: false,
	},
	"client_ca_file": {
		Name:     "client_ca_file",
		Type:     StringType,
		Required: false,
	},
	"ca_file": {
		Name:     "ca_file",
		Type:     StringType,
		Required: false,
	},
	"cert_file": {
		Name:     "cert_file",
		Type:     StringType,
		Required: false,
	},
	"key_file": {
		Name:     "key_file",
		Type:     StringType,
		Required: false,
	},
	"min_version": {
		Name:     "min_version",
		Type:     StringType,
		Required: false,
	},
	"max_version": {
		Name:     "max_version",
		Type:     StringType,
		Required: false,
	},
	"reload_interval": {
		Name:     "reload_interval",
		Type:     DurationType,
		Required: false,
	},
}
