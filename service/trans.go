package service

import (
	"github.com/godcong/go-ffmpeg/ffmpeg"
	"github.com/godcong/go-ffmpeg/ffprobe"
	"log"
)

// ToM3U8 ...
func ToM3U8(path string) error {
	log.Println("start")
	probe := ffprobe.New(path)

	if probe.Run().IsH264AndAAC() {
		b, err := ffmpeg.CopyToMp4(path, path+"_out.mp4")
		if err != nil {
			log.Println(string(b))
			return err
		}
	} else {
		b, err := ffmpeg.TranToMp4(path, path+"_out.mp4")
		if err != nil {
			log.Println(string(b), err)
			return err
		}
	}

	b, err := ffmpeg.TransToTS(path+"_out.mp4", path+"_out.ts")
	if err != nil {
		log.Println(string(b), err)
		return err
	}
	log.Println("end")

	return nil
}
