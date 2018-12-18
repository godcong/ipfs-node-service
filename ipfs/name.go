package ipfs

import (
	"net/url"
	"strconv"
)

// Name ...
type Name struct {
	*api
}

// Self ...
func (*Name) Self() string {
	return "name"
}

// PublishD ...
func (n *Name) PublishD(ipfsPath string) (map[string]string, error) {
	return n.Publish(ipfsPath, true, "", "", "")
}

// PublishWithKey ...
func (n *Name) PublishWithKey(ipfsPath string, key string) (map[string]string, error) {
	return n.Publish(ipfsPath, true, "", "", key)
}

// Publish ...
func (n *Name) Publish(ipfsPath string, resolve bool, lifetime, ttl, key string) (map[string]string, error) {
	v := url.Values{}
	v.Set("arg", ipfsPath)
	v.Set("resolve", strconv.FormatBool(resolve))
	if lifetime == "" {
		lifetime = "24h"
	}
	v.Set("lifetime", lifetime)
	if ttl != "" {
		v.Set("ttl", ttl)
	}

	if key == "" {
		key = "self"
	}
	v.Set("key", key)
	host := URL(n, "publish")
	return get(host, v)
}
