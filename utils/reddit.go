package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type RedditMemeResponse struct {
	PostLink  string   `json:"postLink"`
	SubReddit string   `json:"subreddit"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	NSFW      bool     `json:"nsfw"`
	Spoiler   bool     `json:"spoiler"`
	Author    string   `json:"author"`
	UPs       int      `json:"ups"`
	Preview   []string `json:"preview"`
}

func GetMeme(subreddit string) (*RedditMemeResponse, error) {
	memeEndpoint, err := url.JoinPath("https://meme-api.com/", "gimme", subreddit)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(memeEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	res := &RedditMemeResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetMemeImg(meme *RedditMemeResponse) ([]byte, error) {
	resp, err := http.Get(meme.URL)
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
