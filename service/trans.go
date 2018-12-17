package service

import (
	"github.com/godcong/go-ffmpeg/ffmpeg"
	"github.com/godcong/go-ffmpeg/ffprobe"
	"log"
)

// ToM3U8 ...
func ToM3U8(path string, file string) error {
	log.Println("start", path)
	output := path + file
	probe := ffprobe.New(path)

	if probe.Run().IsH264AndAAC() {
		b, err := ffmpeg.CopyToMp4(output, output+"_out.mp4")
		if err != nil {
			log.Println(string(b))
			return err
		}
	} else {
		b, err := ffmpeg.TranToMp4(output, output+"_out.mp4")
		if err != nil {
			log.Println(string(b), err)
			return err
		}
	}

	b, err := ffmpeg.SplitToTS(output+"_out.mp4", "./split/"+file)
	if err != nil {
		log.Println(string(b), err)
		return err
	}
	log.Println("end")

	return nil
}

// ToM3U8WithKey ...
func ToM3U8WithKey(id string) error {

	output := config.Transfer + "/" + id
	source := config.Upload + "/" + id
	probe := ffprobe.New(source)

	err := client.Set(id, StatusTransferring, 0).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	procFunc := ffmpeg.SplitWithKey
	if probe.Run().IsH264AndAAC() {
		procFunc = ffmpeg.QuickSplitWithKey
	}

	b, err := procFunc(source, output, config.KeyInfoFile, "media", config.M3U8)
	if err != nil {
		log.Println(string(b), err)
		return err
	}

	return nil
}
