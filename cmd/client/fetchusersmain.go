package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// var (
// 	httpClient *http.Client
// )

// type Result struct {
// 	Error    error
// 	Response []byte
// }

func fetchUser(ctx context.Context, id int) ([]byte, error) {
	url := fmt.Sprintf("http://localhost:8080/users/%d", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func doFetchUsersRequests(ctx context.Context, urls []int) <-chan Result {
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

func fetchAllUsers() {
	userIdSize := 10
	httpClient = &http.Client{}
	ctx, _ := context.WithCancel(context.Background())

	var req = make([]int, 0)
	for i := 1; i <= userIdSize; i++ {
		req = append(req, i)
	}

	terminated := doFetchUsersRequests(ctx, req)

	// go func() {
	// 	// Cancel the operation after 1 second.
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Println("Canceling doWork goroutine...")
	// 	cancel()
	// }()

	// 	// var result User
	// 	// json.NewDecoder(resp.Body).Decode(&result)
	// 	fmt.Println(string(b))
	// }

	errCount := 0
	for i := range terminated {
		if i.Error != nil {
			fmt.Println(i.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Println(string(i.Response))
	}
}
