package days

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDays(t *testing.T) {
	tests := []struct {
		start, end string
		days       int
	}{
		{"2021-01-01", "2021-12-31", 365},
		{"2021-01-22", "2021-05-22", 121},
		{"2020-01-01", "2021-12-23", 723},
		{"2021-01-01", "2020-12-31", -1},
		{"2021-01-01", "2020-03-08", -1},
		{"2021-01-01", "2021-01-01", 1},
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
	Days(recorder, request)
	var o OutputDays
	err := json.Unmarshal(recorder.Body.Bytes(), &o)
	if err != nil {
		t.Errorf("Return error for (start, end) = (%s, %s)", s, e)
	}
	return o.Days
}