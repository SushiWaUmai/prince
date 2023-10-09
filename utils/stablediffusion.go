package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/SushiWaUmai/prince/env"
)

type Txt2ImgRequest struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt"`
	Steps          int64  `json:"steps"`
	Width          int64  `json:"width"`
	Height         int64  `json:"height"`
	SamplerName    string `json:"sampler_name"`
}

type Txt2ImgResponse struct {
	Images []string `json:"images"`
	Info   string   `json:"info"`
}

func Txt2Img(prompt string) ([]byte, error) {
	client := &http.Client{}

	payload := Txt2ImgRequest{
		Prompt:         prompt,
		NegativePrompt: "worst quality, bad quality, normal quality, watermarks, image artifacts",
		Steps:          35,
		Width:          512,
		Height:         512,
		SamplerName:    "DPM++ SDE Karras",
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	txt2ImgEndpoint, err := url.JoinPath(env.STABLE_DIFFUSION_ENDPOINT, "sdapi", "v1", "txt2img")
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(txt2ImgEndpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	res := Txt2ImgResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Images) <= 0 {
		return nil, errors.New("Could not generate any images")
	}

	imgData, err := base64.StdEncoding.DecodeString(res.Images[0])
	if err != nil {
		return nil, err
	}

	return imgData, nil
}
