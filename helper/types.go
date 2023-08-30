package helper

import (
	tlsclient "github.com/bogdanfinn/tls-client"
)

type Data struct {
	Version string `json:"version"`
}

type IHelper struct {
	client tlsclient.HttpClient
}
