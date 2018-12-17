package ipfs

import (
	"net/url"
	"strconv"
)

type Name struct {
	*api
}

func (*Name) Self() string {
	return "name"
}

func (n *Name) Publish(ipfs string, resolve bool, lifetime, ttl, key string) (map[string]string, error) {
	v := url.Values{}
	v.Set("arg", ipfs)
	v.Set("resolve", strconv.FormatBool(resolve))
	v.Set("lifetime", "24h")
	v.Set("ttl", ttl)
	if key == "" {
		key = "self"
	}
	v.Set("key", key)
	host := URL(n, "publish")
	return Get(host, v)
}
