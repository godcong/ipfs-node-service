package ipfs

import (
	"net/url"
	"strconv"
)

// Key ...
type Key struct {
	*api
}

// Self ...
func (*Key) Self() string {
	return "key"
}

// Gen ...
func (k *Key) Gen(name, typ string, size int) (map[string]string, error) {
	v := url.Values{}
	v.Set("arg", name)
	if typ == "" {
		typ = "rsa"
	}
	v.Set("type", typ)
	inSize := "2048"
	if size != 0 {
		inSize = strconv.Itoa(size)
	}
	v.Set("size", inSize)

	host := URL(k, "gen")
	return get(host, v)
}
