package utils

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/SushiWaUmai/prince/env"
	"github.com/sashabaranov/go-openai"
)

var OpenAIClient = openai.NewClient(env.OPENAI_API_KEY)

func TranscribeAudio(audioData []byte) (string, error) {
	audioData, err := OggToMp3(audioData)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "audio*.mp3")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}

	req, err := OpenAIClient.CreateTranscription(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: tmpFile.Name(),
		},
	)
	if err != nil {
		return "", err
	}

	return req.Text, nil
}
