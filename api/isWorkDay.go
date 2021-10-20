package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	// 设置代理解决 GitHub 被墙问题
	_ = os.Setenv("http_proxy", "http://127.0.0.1:1080")
	_ = os.Setenv("https_proxy", "http://127.0.0.1:1080")
}

func WorkDay(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "*")

	query := request.URL.Query()
	dates := query["date"]
	if len(dates) < 1 {
		writer.Write(fail("参数错误"))
		return
	}
	var queryDates []time.Time
	yearsSet := make(map[int]bool)
	for _, d := range dates {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			writer.Write(fail("错误参数：" + d))
			return
		}
		queryDates = append(queryDates, t)
		yearsSet[t.Year()] = true
	}
	var transfer []string
	var holiday []string
	for k := range yearsSet {
		var result *Result
		for i := 0; i < 5; i++ { // 重试 5 次
			var err error
			result, err = getFromUrl(k)
			if err == nil {
				break
			}
		}

		if result == nil {
			writer.Write(fail("数据错误：" + strconv.Itoa(k) + " 年假期数据缺失！"))
			return
		}

		transfer = append(transfer, result.Transfer...)
		holiday = append(holiday, result.Holidays...)
	}

	transfers = buildMap(transfer)
	holidays = buildMap(holiday)

	var result []bool

	for _, t := range queryDates {
		r := isTransfer(t) || (!isHoliday(t) && !isWeekend(t))
		result = append(result, r)
	}

	writer.Write(success(result, "查询成功！"))

}

var transfers map[string]bool
var holidays map[string]bool

type Output struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Days    []bool    `json:"days"`
}

type Result struct {
	Year     int      `json:"year"`
	Holidays []string `json:"holidays"`
	Transfer []string `json:"transfer"`
}

func fail(msg string) []byte {

	marshal, _ := json.Marshal(&Output{
		false, msg, nil,
	})

	return marshal
}

func success(days []bool, msg string) []byte {
	marshal, _ := json.Marshal(&Output{
		true, msg, days,
	})
	return marshal
}

func getFromUrl(year int) (*Result, error) {
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/"+
		"yzsj98/api/main/holidayJsons/%d.json", year))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result Result

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	if err = json.Unmarshal(respByte, &result); err != nil {
		return nil, err
	}
	return &result, nil

}

func buildMap(in []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range in {
		m[s] = true
	}
	return m

}
func isWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func isTransfer(t time.Time) bool {
	format := t.Format("2006-01-02")
	_, existed := transfers[format]
	return existed
}

func isHoliday(t time.Time) bool {
	format := t.Format("2006-01-02")
	_, existed := holidays[format]
	return existed
}