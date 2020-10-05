package http_test

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ditointernet/go-dito-http/http"

	h "net/http"
)

func TestGet(t *testing.T) {
	tt := []struct {
		name        string
		server      func(t *testing.T) *httptest.Server
		options     http.RequestOptions
		expectedErr error
	}{
		{
			name: "must build query URL with the given QueryParams",
			server: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(h.HandlerFunc(func(w h.ResponseWriter, r *h.Request) {
					if r.URL.RawQuery != "bla=ble" {
						t.Errorf("errorrrr")
					}

					io.WriteString(w, "ping")
				}))
			},
			options: http.RequestOptions{
				QueryParams: map[string]string{
					"bla": "ble",
				},
			},
			expectedErr: nil,
		},
		// {
		// 	name: "must build query URL with the given Headers"
		// },
		// {
		// 	name: "must send Body in Get requests"
		// },
		// {
		// 	name: "must send trace id in Headers"
		// },
		// {
		// 	name: "must fail with error Kind 400"
		// },
		// {
		// 	name: "must fail with error Kind 404"
		// },
		// {
		// 	name: "must fail with error Kind 413 (unknown)"
		// },
		// {
		// 	name: "must fail with error Kind 500"
		// },
		// {
		// 	name: "must fail with request duration exceeds the given timeout"
		// },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := tc.server(t)

			client := http.NewClient(http.ClientInput{})
			_, err := client.Get(context.Background(), ts.URL, tc.options)

			if err != tc.expectedErr {
				t.Fatalf(err.Error())
			}
		})
	}
}
