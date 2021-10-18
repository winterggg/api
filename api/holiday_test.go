package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	// 设置代理解决 GitHub 被墙问题
	_ = os.Setenv("http_proxy", "http://127.0.0.1:1080")
	_ = os.Setenv("https_proxy", "http://127.0.0.1:1080")
}

func TestHandler(t *testing.T) {
	tests := []struct {
		start, end string
		days       int
	}{
		{"2021-01-01", "2021-12-31", 250},
		{"2021-01-22", "2021-05-22", 81},
		{"2020-01-01", "2021-12-23", 495},
		{"2021-01-01", "2020-12-31", -1},
		{"2021-01-01", "2020-03-08", -1},
		{"2021-01-01", "2021-01-01", 0},
		{"2021-01-05", "2021-01-05", 1},
	}

	for _, test := range tests {
		if actual := getResult(test.start, test.end, t); actual != test.days {
			t.Errorf("(start=%s, end=%s): Expected %d, got %d", test.start, test.end, test.days, actual)
		}

	}

}

func getResult(s, e string, t *testing.T) int {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "?start="+s+"&end="+e, nil)
	Handler(recorder, request)
	var o Output
	err := json.Unmarshal(recorder.Body.Bytes(), &o)
	if err != nil {
		t.Errorf("Return error for (start, end) = (%s, %s)", s, e)
	}
	return o.Days
}
