package handler

import (
	"github.com/AlLevykin/cutwell/internal/app/store"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_CreateShortLink(t *testing.T) {
	type args struct {
		lnk    string
		keyLen int
	}
	type want struct {
		code        int
		contentType string
		contentLen  int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"ok",
			args{
				"ya.ru",
				9,
			},
			want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				contentLen:  9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.args.lnk)
			req := httptest.NewRequest(http.MethodGet, "/", body)
			w := httptest.NewRecorder()
			ls := store.NewLinkStore(tt.args.keyLen)
			r := NewRouter(ls)
			r.CreateShortLink(w, req)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			key := string(resBody)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
			if len(key) != tt.want.contentLen {
				t.Errorf("Expected status code %d, got %s len = %d", tt.want.contentLen, key, len(key))
			}
		})
	}
}
