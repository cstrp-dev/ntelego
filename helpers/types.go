package helpers

import (
	tlsclient "github.com/bogdanfinn/tls-client"
)

type Data struct {
	Version string `json:"version"`
}

type Response struct {
	ID      string `json:"id"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

type IHelper struct {
	client tlsclient.HttpClient
}
