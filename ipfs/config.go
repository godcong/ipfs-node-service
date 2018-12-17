package ipfs

import (
	"github.com/json-iterator/go"
	"net/http"
	"net/url"
	"strings"
)

type Host interface {
	Addr() string
	Prefix() string
	Version() string
	Self() string
}

type Config struct {
	addr string
}

func (c *Config) Addr() string {
	return c.addr
}

func (c *Config) SetAddr(addr string) {
	c.addr = addr
}

type api struct {
	Config
	prefix  string
	version string
}

func (a *api) Version() string {
	return a.version
}

func (a *api) SetVersion(version string) {
	a.version = version
}

func (a *api) Prefix() string {
	return a.prefix
}

func (a *api) SetPrefix(prefix string) {
	a.prefix = prefix
}

func newApi(config Config, prefix, version string) *api {
	return &api{Config: config, prefix: prefix, version: version}
}

func NewConfig(addr string) *Config {
	return &Config{addr: addr}
}

func (c Config) Version0() *api {
	return newApi(c, "api", "v0")
}

func (a *api) Key() *Key {
	return &Key{
		api:  a,
		Self: "key",
	}
}

func (a *api) Name() *Name {
	return &Name{
		api:  a,
		Self: "name",
	}
}

func URL(h Host, act string) string {
	return strings.Join([]string{h.Addr(), h.Prefix(), h.Version(), h.Self(), act}, "/")
}

func Get(host string, values url.Values) (map[string]string, error) {
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
