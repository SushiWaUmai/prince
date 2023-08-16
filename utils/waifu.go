package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type animeImage struct {
	Url string `json:"url"`
}

type animeResponse struct {
	Images []animeImage `json:"images"`
}

func GetWaifu(category string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.waifu.im/search/?included_tags=%s", category))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data animeResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	resp, err = http.Get(data.Images[0].Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
