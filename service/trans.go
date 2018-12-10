package service

import (
	"github.com/godcong/go-ffmpeg/ffmpeg"
	"github.com/godcong/go-ffmpeg/ffprobe"
	"log"
)

// ToM3U8 ...
func ToM3U8(path string) error {
	log.Println("start")
	stream, err := ffprobe.FilterStream(path)
	if err != nil {
		return err
	}
	acc := ffprobe.CheckAAC(stream)
	h264 := ffprobe.CheckH264(stream)
	if acc && h264 {
		b, err := ffmpeg.CopyToMp4(path, path+"/out.mp4")
		if err != nil {
			log.Println(string(b))
			return err
		}
	}
	b, err := ffmpeg.TranToMp4(path, path+"/out.mp4")
	if err != nil {
		log.Println(string(b))
		return err
	}
	log.Println("end")
	return nil
}
