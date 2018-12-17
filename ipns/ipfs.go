package ipns

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

type Key struct {
	*api
	self string
}

func (k *Key) Self() string {
	return k.self
}

func (k *Key) SetSelf(self string) {
	k.self = self
}

func newApi(config Config, prefix, version string) *api {
	return &api{Config: config, prefix: prefix}
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
		self: "key",
	}
}

func (k *Key) Gen(name, typ, size string) (map[string]string, error) {
	v := url.Values{
		"arg":  []string{name},
		"type": []string{typ},
		"size": []string{size},
	}
	resp, err := http.Get(URL(k, "gen") + "?" + v.Encode())
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

func URL(h Host, act string) string {
	return strings.Join([]string{h.Addr(), h.Prefix(), h.Version(), h.Self(), act}, "/")
}
