package ipfs

import (
	"github.com/json-iterator/go"
	"net/http"
	"net/url"
	"strings"
)

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
