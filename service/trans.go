package service

import (
	"github.com/godcong/ipfs-media-service/ffmpeg"
	"github.com/godcong/ipfs-media-service/ffprobe"
	"github.com/godcong/ipfs-media-service/oss"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// toM3U8WithKey ...
func toM3U8WithKey(id, source, dest string, key string) error {

	output := filepath.Join(dest, id)
	probe := ffprobe.New(source)

	//err := rdsQueue.Set(id, StatusTransferring, 0).Err()
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}

	procFunc := ffmpeg.SplitWithKey
	if probe.Run().IsH264AndAAC() {
		procFunc = ffmpeg.QuickSplitWithKey
	}

	b, err := procFunc(source, output, key, "media", "media.m3u8")
	if err != nil {
		log.Println(string(b), err)
		return err
	}

	return nil
}

// toM3U8 ...
func toM3U8(id string, source, dest string) error {

	output := filepath.Join(dest, id)
	_ = os.MkdirAll(output, os.ModePerm) //ignore err

	//source = source + "/" + id + "/" + fileName
	probe := ffprobe.New(source)

	procFunc := ffmpeg.Split
	if probe.Run().IsH264AndAAC() {
		procFunc = ffmpeg.QuickSplit
	}

	b, err := procFunc(source, output, "media", "media.m3u8")
	if err != nil {
		log.Println(string(b), err)
		return err
	}

	return nil
}

func downloadFromOSS(info *Streamer) error {
	server := oss.Server()

	p := oss.NewProgress()
	p.SetObjectKey(info.ObjectKey)
	p.SetPath(info.FileSource + "/" + info.ID)
	err := server.Download(p)
	if err != nil {
		return err
	}
	return nil
}
