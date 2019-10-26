package httphelper

import (
	"net/http"
	"time"
)

var DefaultClient = NewClient()

func NewClient() (client *http.Client) {
	return &http.Client{
		Timeout: time.Minute,
	}
}