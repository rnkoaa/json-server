package main

import (
	"context"
	"net/http"
)

var (
	httpClient *http.Client
)

type Result struct {
	Error    error
	Response []byte
}

func doRequests(ctx context.Context, urls []int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		for _, userId := range urls {
			res, err := fetchUser(ctx, userId)
			result := Result{Error: err, Response: res}
			select {
			case <-ctx.Done():
				return
			case responses <- result:
			}
		}
	}()

	return responses
}

func main() {
}
