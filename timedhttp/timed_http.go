package timedhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// Get replaces the http package introducing a timeout
func Get(theURL string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", theURL, nil)
	if err != nil {
		return
	}
	timeout := viper.GetInt("HTTP_TIMEOUT")
	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(timeout)*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	client := http.DefaultClient
	return client.Do(req)
}
