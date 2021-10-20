package isWorkDay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func init() {
	// 设置代理解决 GitHub 被墙问题
	_ = os.Setenv("http_proxy", "http://127.0.0.1:1080")
	_ = os.Setenv("https_proxy", "http://127.0.0.1:1080")
}

func TestWorkDay(t *testing.T) {
	tests := []struct {
		query  string
		result []bool
	}{
		{"?date=2021-10-01&date=2021-10-02&date=2021-10-08", []bool{false, false, true}},
		{"?date=2020-10-01&date=2020-10-02&date=2020-10-08", []bool{false, false, false}},
	}

	for _, test := range tests {
		if actual := getResult(t, test.query); !reflect.DeepEqual(actual, test.result) {
			t.Errorf("Error for query [%s]", test.query)
		}
	}
}

func getResult(t *testing.T, query string) []bool {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", query, nil)
	WorkDay(recorder, request)
	var o Output
	err := json.Unmarshal(recorder.Body.Bytes(), &o)
	if err != nil {
		t.Errorf("error")
	}
	return o.Days
}
