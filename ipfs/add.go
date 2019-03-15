package ipfs

import (
	"github.com/ipfs/go-ipfs-cmdkit/files"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"path"
	"strconv"
)

// Add ...
func (a *api) AddDir(dir string) (map[string]string, error) {
	stat, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}

	sf, err := files.NewSerialFile(path.Base(dir), dir, false, stat)
	if err != nil {
		return nil, err
	}
	log.Info(sf)
	slf := files.NewSliceFile("", dir, []files.File{sf})
	reader := files.NewMultiFileReader(slf, true)
	a.SetSelf("add")
	host := URL(a, "")
	v := url.Values{}
	v.Set("recursive", strconv.FormatBool(true))
	return post(host, v, reader)
}
