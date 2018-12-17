package ipfs

import (
	"net/url"
)

type Key struct {
	*api
}

func (*Key) Self() string {
	return "key"
}

func (k *Key) Gen(name, typ, size string) (map[string]string, error) {
	v := url.Values{}
	v.Set("arg", name)
	v.Set("type", typ)
	v.Set("size", size)

	host := URL(k, "gen")
	return Get(host, v)
}
