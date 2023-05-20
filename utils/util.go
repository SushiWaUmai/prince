package utils

import (
	"bytes"
	"io"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func OggToMp3(audioData []byte) ([]byte, error) {
	// ffmpeg -i $inFileName -acodec libmp3lame -y $outFileName
	tmpFile, err := os.CreateTemp("", "audio*.ogg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, bytes.NewReader(audioData))
	if err != nil {
		return nil, err
	}

	tmpFileOut, err := os.CreateTemp("", "audio*.mp3")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFileOut.Name())

	err = ffmpeg.Input(tmpFile.Name()).Output(tmpFileOut.Name(), ffmpeg.KwArgs{
		"acodec": "libmp3lame",
	}).OverWriteOutput().Run()

	if err != nil {
		return nil, err
	}

	return os.ReadFile(tmpFileOut.Name())
}
