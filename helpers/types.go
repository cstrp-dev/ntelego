package helpers

import (
	tlsclient "github.com/bogdanfinn/tls-client"
)

type Set struct {
	data map[string]bool
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
