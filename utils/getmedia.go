package utils

import (
	"context"
	"io/ioutil"
	"net/http"
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
		buffer, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(mimeType, "image") {
		imgMsg, err := CreateImgMessage(client, buffer)
		if err != nil {
			return nil, err
		}

		response := &waProto.Message{
			ImageMessage: imgMsg,
		}
		return response, nil
	} else if strings.Contains(mimeType, "audio") {
		uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaAudio)
		if err != nil {
			return nil, err
		}

		audioMsg := &waProto.AudioMessage{
			Mimetype:      &mimeType,
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
	} else {
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

		buffer, err := ioutil.ReadAll(downloadResult)
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
}
