package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestShortURL(t *testing.T) {
	t.Run("shortUrl testing", func(t *testing.T) {
		res := shortURL()
		require.Equal(t, len(res), 8)
	})
}

func TestHandler(t *testing.T) {
	tests := []struct{
		name string
		url string
		want int
	}{
		{
			"small length of url",
			"/u8432d",
			http.StatusBadRequest,
		},
		{
			"big length of url",
			"/u843241dsSAd",
			http.StatusBadRequest,
		},
		{
			"a lot of / in url",
			"/accept/downfall/rya",
			http.StatusBadRequest,
		},
		{
			"Correct link length but the the route is incorrect ",
			"/acpew/u8",
			http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()
			
			handler(w, req)
			res := w.Result()

			assert.Equal(t, res.StatusCode, test.want)
		})
	}
}