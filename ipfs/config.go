package ipfs

import (
	"github.com/ipfs/go-ipfs-cmdkit/files"
	"github.com/json-iterator/go"
	"io"
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
	AddDir(dir string) (map[string]string, error)
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
	self    string
}

// Self ...
func (a *api) Self() string {
	return a.self
}

// SetSelf ...
func (a *api) SetSelf(self string) {
	a.self = self
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

	link := []string{h.Addr(), h.Prefix(), h.Version(), h.Self()}
	if act != "" {
		link = append(link, act)
	}
	return url + strings.Join(link, "/")
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

func post(host string, values url.Values, body io.Reader) (map[string]string, error) {
	req, err := http.NewRequest("POST", host+"?"+values.Encode(), body)
	if err != nil {
		return nil, err
	}

	if fr, ok := body.(*files.MultiFileReader); ok {
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+fr.Boundary())
		req.Header.Set("Content-Disposition", "form-data: name=\"files\"")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	//contentType := resp.Header.Get("Content-Type")
	//parts := strings.Split(contentType, ";")
	//contentType = parts[0]

	m := make(map[string]string)
	dec := jsoniter.NewDecoder(resp.Body)
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
