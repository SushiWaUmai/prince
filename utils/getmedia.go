package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/wader/goutubedl"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

func GetMedia(client *whatsmeow.Client, fetchUrl string) (*waProto.Message, error) {
	resp, err := http.Get(fetchUrl)
	var buffer []byte
	var mimeType string
	if err == nil {
		mimeType = resp.Header.Get("Content-Type")
		buffer, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if strings.Contains(mimeType, "image") {
			response, err := getImage(client, buffer)

			if err == nil {
				return response, nil
			}
		}

		if strings.Contains(mimeType, "audio") {
			response, err := getAudio(client, buffer)

			if err == nil {
				return response, nil
			}
		}
	}

	// yt-dlp
	{
		response, err := getYtDlp(client, fetchUrl)
		if err == nil {
			return response, nil
		}
	}

	// spotdl
	{
		response, err := getSpotDl(client, fetchUrl)
		if err == nil {
			return response, nil
		}
	}

	return nil, errors.New("Could not download url")
}

func getImage(client *whatsmeow.Client, buffer []byte) (*waProto.Message, error) {
	imgMsg, err := CreateImgMessage(client, buffer)
	if err != nil {
		return nil, err
	}

	response := &waProto.Message{
		ImageMessage: imgMsg,
	}
	return response, nil
}

func getAudio(client *whatsmeow.Client, buffer []byte) (*waProto.Message, error) {
	uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaAudio)
	if err != nil {
		return nil, err
	}

	audioMsg := &waProto.AudioMessage{
		Mimetype:      proto.String("audio/ogg; codecs=opus"),
		Url:           &uploadResp.URL,
		DirectPath:    &uploadResp.DirectPath,
		MediaKey:      uploadResp.MediaKey,
		FileEncSha256: uploadResp.FileEncSHA256,
		FileSha256:    uploadResp.FileSHA256,
		FileLength:    &uploadResp.FileLength,
	}

	response := &waProto.Message{
		AudioMessage: audioMsg,
	}
	return response, nil
}

func getYtDlp(client *whatsmeow.Client, fetchUrl string) (*waProto.Message, error) {
	// yt-dlp
	goutubedl.Path = "yt-dlp"
	result, err := goutubedl.New(context.Background(), fetchUrl, goutubedl.Options{})
	if err != nil {
		return nil, err
	}
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		return nil, err
	}
	defer downloadResult.Close()

	buffer, err := io.ReadAll(downloadResult)
	if err != nil {
		return nil, err
	}

	uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaVideo)
	if err != nil {
		return nil, err
	}

	videoMsg := &waProto.VideoMessage{
		Mimetype:      proto.String(http.DetectContentType(buffer)),
		Url:           &uploadResp.URL,
		DirectPath:    &uploadResp.DirectPath,
		MediaKey:      uploadResp.MediaKey,
		FileEncSha256: uploadResp.FileEncSHA256,
		FileSha256:    uploadResp.FileSHA256,
		FileLength:    &uploadResp.FileLength,
	}

	response := &waProto.Message{
		VideoMessage: videoMsg,
	}
	return response, nil
}

func getSpotDl(client *whatsmeow.Client, fetchUrl string) (*waProto.Message, error) {
	cmd := exec.CommandContext(
		context.Background(),
		"spotdl",
	)

	cmd.Args = append(cmd.Args, fetchUrl)
	cmd.Args = append(cmd.Args, "--format")
	cmd.Args = append(cmd.Args, "mp3")
	cmd.Args = append(cmd.Args, "--output")
	cmd.Args = append(cmd.Args, "audio.{output-ext}")

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	defer os.Remove("audio.mp3")

	buffer, err := os.ReadFile("audio.mp3")
	if err != nil {
		return nil, err
	}

	audioBuffer, err := Mp3ToOgg(buffer)
	if err != nil {
		return nil, err
	}

	return getAudio(client, audioBuffer)
}
