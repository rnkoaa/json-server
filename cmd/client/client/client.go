package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	mediaType      = "application/json"
	defaultBaseURL = "http://localhost:8080"
	userAgent      = "goplaceholder/1.0.0"
)

var (
	ErrUnkown = errors.New("unkown rest error")
)

type Client struct {
	httpClient *http.Client
	// Base URL for API requests.
	baseURL *url.URL

	// User agent for client
	UserAgent string

	User    UserService
	Todo    TodoService
	Post    PostService
	Album   AlbumService
	Comment CommentService
	Photo   PhotoService
}

func NewClient(httpClient *http.Client, baseUrl string) *Client {
	c := http.DefaultClient
	if httpClient != nil {
		c = httpClient
	}

	var b = defaultBaseURL
	if baseUrl != "" {
		b = baseUrl
	}

	baseURL, _ := url.Parse(b)

	client := &Client{
		httpClient: c,
		baseURL:    baseURL,
		UserAgent:  userAgent,
	}

	client.User = &UserServiceOp{client}
	client.Todo = &TodoServiceOp{client}
	client.Post = &PostServiceOp{client}
	client.Album = &AlbumServiceOp{client}
	client.Comment = &CommentServiceOp{client}
	client.Photo = &PhotoServiceOp{client}

	return client
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	// RequestID returned from the API, useful to contact support.
	RequestID string `json:"request_id"`
}

// Rate contains the rate limit for the current client.
type Rate struct {
	// The number of request per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	// Reset Timestamp `json:"reset"`
}

// SetBaseURL is a client option for setting the base URL.
func (c *Client) SetBaseURL(bu string) {
	u, err := url.Parse(bu)
	if err != nil {
		// return err
		log.Fatalf("error parsing base url %v", err)
	}

	c.baseURL = u
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.httpClient, req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	if r.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.RequestID, r.Message)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

// Response is a local response object wrapping several properties of http.Response
type Response struct {
	*http.Response
}

func deserialize(b []byte, v interface{}) error {
	if b == nil {
		return errors.New("nil input")
	}
	return json.Unmarshal(b, v)

}
func conflict(res *Response) bool {
	return res.StatusCode == 409
}

func successful(res *Response) bool {
	return res.StatusCode >= 200 && res.StatusCode <= 299
}

func get(ctx context.Context, client *Client, path string, response interface{}) (*Response, error) {
	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	if successful(res) {
		return res, nil
	}

	return res, ErrUnkown
}
