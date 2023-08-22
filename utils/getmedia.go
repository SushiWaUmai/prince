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
)

var maxVideoLength float64 = 30 * 60

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

		if strings.Contains(mimeType, "text/plain") {
			response := CreateTextMessage(string(buffer))
			return response, nil
		}

		if strings.Contains(mimeType, "image") {
			response, err := CreateImgMessage(client, buffer)

			if err == nil {
				return response, nil
			}
		}

		if strings.Contains(mimeType, "audio") {
			response, err := CreateAudioMessage(client, buffer)

			if err == nil {
				return response, nil
			}
		}

		if strings.Contains(mimeType, "video") {
			response, err := CreateVideoMessage(client, buffer)

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

func getYtDlp(client *whatsmeow.Client, fetchUrl string) (*waProto.Message, error) {
	// yt-dlp
	goutubedl.Path = "yt-dlp"
	result, err := goutubedl.New(context.Background(), fetchUrl, goutubedl.Options{})
	if err != nil {
		return nil, err
	}
	seconds := result.Info.Duration
	if seconds > maxVideoLength {
		return nil, errors.New("Video is too long")
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

	response, err := CreateVideoMessage(client, buffer)
	if err != nil {
		return nil, err
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

	return CreateAudioMessage(client, buffer)
}
