package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	URL "net/url"
	"time"

	errors "github.com/ditointernet/go-dito-errors"
	"go.opentelemetry.io/otel/api/trace"
)

type RequestOptions struct {
	Body        []byte
	QueryParams map[string]string
	Headers     map[string]string
	Timeout     time.Duration
}

type Response struct {
	Status int
	Body   []byte
}

// Decode ...
func (r Response) Decode(out interface{}) error {
	if err := json.Unmarshal(r.Body, out); err != nil {
		return errors.New(err)
	}

	return nil
}

type ClientInput struct {
	Tracer trace.Tracer
}

type Client struct {
	in ClientInput
}

func NewClient(in ClientInput) Client {
	return Client{in: in}
}

func (c Client) Get(ctx context.Context, baseURL string, options RequestOptions) (*Response, error) {
	if c.in.Tracer != nil {
		spanCtx, span := c.in.Tracer.Start(ctx, "http.Get")
		ctx = spanCtx
		defer span.End()
	}

	if options.Timeout > 0 {
		tCtx, cancel := context.WithTimeout(ctx, options.Timeout)
		ctx = tCtx
		defer cancel()
	}

	queryValues := URL.Values{}
	for k, v := range options.QueryParams {
		queryValues.Add(k, v)
	}

	url, err := URL.Parse(baseURL)
	if err != nil {
		return nil, errors.New(err)
	}
	url.RawQuery = queryValues.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, errors.New(err)
	}

	for k, v := range options.Headers {
		req.Header.Add(k, v)
	}

	trace.DefaultHTTPPropagator().Inject(ctx, req.Header)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Response{Status: res.StatusCode, Body: body}, errors.New(err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = fmt.Errorf("failed on http request. got status %v", res.StatusCode)
		return &Response{Status: res.StatusCode, Body: body}, errors.New(err, errors.ErrorKind(res.StatusCode))
	}

	return &Response{Status: res.StatusCode, Body: body}, nil
}

func (c Client) Post(ctx context.Context, url string, options RequestOptions) (*Response, error) {
	return nil, nil
}

func (c Client) Delete(ctx context.Context, url string, options RequestOptions) (*Response, error) {
	return nil, nil
}
