package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-ipfs-cmdkit/files"
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	//json := jsoniter.ConfigCompatibleWithStandardLibrary

	req, err := http.NewRequest("POST", host+"?"+values.Encode(), body)
	if err != nil {
		return nil, err
	}

	if fr, ok := body.(*files.MultiFileReader); ok {
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+fr.Boundary())
		req.Header.Set("Content-Disposition", "form-data: name=\"files\"")
	}

	resp, err := client().Do(req)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, ";")
	contentType = parts[0]
	m := make(map[string]string)

	switch {
	case resp.StatusCode == http.StatusNotFound:
		m["Message"] = "command not found"
	case contentType == "text/plain":
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response read error: %s\n", err)
		}
		m["Message"] = string(out)
	case contentType == "application/json":
		body, err := ioutil.ReadAll(resp.Body)
		dec := json.NewDecoder(bytes.NewReader(body))
		//var fileInfo []map[string]string
		out := map[string]string{}
		for {
			err = dec.Decode(&out)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			m = out
			//fileInfo = append(fileInfo, out)
		}

	default:
		_, _ = fmt.Fprintf(os.Stderr, "ipfs-shell: warning! unhandled response encoding: %s", contentType)
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ipfs-shell: response read error: %s\n", err)
			return nil, err
		}
		m["Message"] = fmt.Sprintf("unknown ipfs-shell error encoding: %q - %q", contentType, out)
	}

	return m, nil
}

func client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}
}
