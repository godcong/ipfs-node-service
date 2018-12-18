package ipfs

import (
	"github.com/json-iterator/go"
	"net/http"
	"net/url"
	"strings"
)

// Host ...
type Host interface {
	TLS() bool
	Addr() string
	Prefix() string
	Version() string
	Self() string
}

// API ...
type API interface {
	Key() *Key
	Name() *Name
}

// Config ...
type Config struct {
	tls  bool
	addr string
}

// TLS ...
func (c *Config) TLS() bool {
	return c.tls
}

// SetTLS ...
func (c *Config) SetTLS(tls bool) {
	c.tls = tls
}

// Addr ...
func (c *Config) Addr() string {
	return c.addr
}

// SetAddr ...
func (c *Config) SetAddr(addr string) {
	c.addr = addr
}

type api struct {
	Config
	prefix  string
	version string
}

// Version ...
func (a *api) Version() string {
	return a.version
}

// SetVersion ...
func (a *api) SetVersion(version string) {
	a.version = version
}

// Prefix ...
func (a *api) Prefix() string {
	return a.prefix
}

// SetPrefix ...
func (a *api) SetPrefix(prefix string) {
	a.prefix = prefix
}

func newAPI(config Config, prefix, version string) API {
	return &api{Config: config, prefix: prefix, version: version}
}

// NewConfig ...
func NewConfig(addr string) *Config {
	return &Config{
		tls:  false,
		addr: addr,
	}
}

// VersionAPI ...
func (c Config) VersionAPI(ver string) API {
	return newAPI(c, "api", ver)
}

// Key ...
func (a *api) Key() *Key {
	return &Key{
		api: a,
	}
}

// Name ...
func (a *api) Name() *Name {
	return &Name{
		api: a,
	}
}

// URL ...
func URL(h Host, act string) string {
	url := "http://"
	if h.TLS() {
		url = "https://"
	}
	return url + strings.Join([]string{h.Addr(), h.Prefix(), h.Version(), h.Self(), act}, "/")
}

func get(host string, values url.Values) (map[string]string, error) {
	resp, err := http.Get(host + "?" + values.Encode())
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	dec := jsoniter.NewDecoder(resp.Body)
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
